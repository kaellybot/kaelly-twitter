package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kaellybot/kaelly-twitter/models/dtos"
	"github.com/rs/zerolog/log"
)

var (
	reHashtag    = regexp.MustCompile(`\B(\#\S+\b)`)
	reTwitterURL = regexp.MustCompile(`https:(\/\/t\.co\/([A-Za-z0-9]|[A-Za-z]){10})`)
	reUsername   = regexp.MustCompile(`\B(\@\S{1,15}\b)`)
)

func (service *Impl) fetchTweets(userID string, maxTweetsNbr int) ([]*dtos.Tweet, error) {
	variables := map[string]any{
		"userId":                                 userID,
		"count":                                  maxTweetsNbr,
		"includePromotedContent":                 false,
		"withQuickPromoteEligibilityTweetFields": false,
		"withVoice":                              false,
	}

	features := map[string]bool{
		"rweb_video_screen_enabled":                                               false,
		"profile_label_improvements_pcf_label_in_post_enabled":                    false,
		"responsive_web_profile_redirect_enabled":                                 false,
		"rweb_tipjar_consumption_enabled":                                         true,
		"verified_phone_label_enabled":                                            false,
		"creator_subscriptions_tweet_preview_api_enabled":                         true,
		"responsive_web_graphql_timeline_navigation_enabled":                      true,
		"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
		"premium_content_api_read_enabled":                                        false,
		"communities_web_enable_tweet_community_results_fetch":                    true,
		"c9s_tweet_anatomy_moderator_badge_enabled":                               true,
		"responsive_web_grok_analyze_button_fetch_trends_enabled":                 false,
		"responsive_web_grok_analyze_post_followups_enabled":                      true,
		"responsive_web_jetfuel_frame":                                            true,
		"responsive_web_grok_share_attachment_enabled":                            true,
		"articles_preview_enabled":                                                true,
		"responsive_web_edit_tweet_api_enabled":                                   true,
		"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
		"view_counts_everywhere_api_enabled":                                      true,
		"longform_notetweets_consumption_enabled":                                 true,
		"responsive_web_twitter_article_tweet_consumption_enabled":                true,
		"tweet_awards_web_tipping_enabled":                                        false,
		"responsive_web_grok_show_grok_translated_post":                           false,
		"responsive_web_grok_analysis_button_from_backend":                        true,
		"creator_subscriptions_quote_tweet_preview_enabled":                       false,
		"freedom_of_speech_not_reach_fetch_enabled":                               true,
		"standardized_nudges_misinfo":                                             true,
		"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
		"longform_notetweets_rich_text_read_enabled":                              true,
		"longform_notetweets_inline_media_enabled":                                true,
		"responsive_web_grok_image_annotation_enabled":                            true,
		"responsive_web_grok_imagine_annotation_enabled":                          true,
		"responsive_web_grok_community_note_auto_translation_is_enabled":          false,
		"responsive_web_enhance_cards_enabled":                                    false,
	}

	u, errURL := url.Parse(baseURL)
	if errURL != nil {
		return nil, errURL
	}

	q := u.Query()
	q.Set("variables", mapToJSONString(variables))
	q.Set("features", mapToJSONString(features))
	u.RawQuery = q.Encode()

	req, errReq := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	if errReq != nil {
		return nil, errReq
	}

	req.Header.Set("X-CSRF-Token", service.csrfToken)
	req.Header.Set("authorization", fmt.Sprintf("Bearer %v", bearerToken))
	req.Header.Set("Cookie", fmt.Sprintf("auth_token=%v; ct0=%v", service.authToken, service.csrfToken))
	req.Header.Set("User-Agent", service.userAgent)

	client := &http.Client{}
	resp, errDo := client.Do(req)
	if errDo != nil {
		return nil, errDo
	}
	defer resp.Body.Close()

	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return nil, errRead
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status %s: %s", resp.Status, body)
	}

	var timeline dtos.TimelineV2
	if errTimeline := json.Unmarshal(body, &timeline); errTimeline != nil {
		return nil, errTimeline
	}

	return extractTweets(&timeline), nil
}

//nolint:gocognit,goconst // Clearer as it is.
func extractTweets(timeline *dtos.TimelineV2) []*dtos.Tweet {
	var tweets []*dtos.Tweet
	for _, instruction := range timeline.Data.User.Result.Timeline.Timeline.Instructions {
		for _, entry := range instruction.Entries {
			if entry.Content.CursorType == "Bottom" {
				continue
			}
			if entry.Content.ItemContent.TweetResults.Result.Typename == "Tweet" ||
				entry.Content.ItemContent.TweetResults.Result.Typename == "TweetWithVisibilityResults" {
				if tweet := parseResult(&entry.Content.ItemContent.TweetResults.Result); tweet != nil {
					tweets = append(tweets, tweet)
				}
			}
			if len(entry.Content.Items) > 0 {
				for _, item := range entry.Content.Items {
					if tweet := parseResult(&item.Item.ItemContent.TweetResults.Result); tweet != nil {
						tweets = append(tweets, tweet)
					}
				}
			}
		}
	}

	return tweets
}

//nolint:gocognit,nestif // Clearer as it is.
func parseResult(result *dtos.Result) *dtos.Tweet {
	if result.NoteTweet.NoteTweetResults.Result.Text != "" {
		result.Legacy.FullText = result.NoteTweet.NoteTweetResults.Result.Text
	}
	legacy := &result.Legacy
	user := &result.Core.UserResults.Result
	if result.Typename == "TweetWithVisibilityResults" {
		legacy = &result.Tweet.Legacy
		user = &result.Tweet.Core.UserResults.Result
	}
	tw := parseLegacyTweet(user, legacy)
	if tw == nil {
		return nil
	}
	if tw.Views == 0 && result.Views.Count != "" {
		tw.Views, _ = strconv.Atoi(result.Views.Count)
	}
	if result.QuotedStatusResult.Result != nil {
		tw.QuotedStatus = parseResult(result.QuotedStatusResult.Result)
	}

	// Get videos from cards
	for _, v := range result.Tweet.Card.Legacy.BindingValues {
		if v.Key == "unified_card" {
			var card dtos.UnifiedCard
			err := json.Unmarshal([]byte(v.Value.StringValue), &card)
			if err != nil {
				continue
			}

			for _, media := range card.MediaEntities {
				if media.Type == "video" {
					var vid dtos.Video

					vid.ID = media.IDStr
					vid.Preview = media.MediaURLHTTPS

					var bitrate int
					for _, variant := range media.VideoInfo.Variants {
						switch variant.ContentType {
						case "video/mp4":
							if variant.Bitrate > bitrate {
								bitrate = variant.Bitrate
								vid.URL = variant.URL
							}
						case "application/x-mpegURL":
							vid.HLSURL = variant.URL
						}
					}

					if vid.URL != "" {
						tw.Videos = append(tw.Videos, vid)
					}
				}
			}
		}
	}

	return tw
}

//nolint:gocognit,gocyclo,cyclop,funlen // Clearer as it is.
func parseLegacyTweet(user *dtos.LegacyUserResult, tweet *dtos.LegacyRawTweet) *dtos.Tweet {
	tweetID := tweet.IDStr
	if tweetID == "" {
		return nil
	}
	text := expandURLs(tweet.FullText, tweet.Entities.URLs, tweet.ExtendedEntities.Media)
	username := user.Core.ScreenName
	name := user.Core.Name
	tw := &dtos.Tweet{
		ConversationID: tweet.ConversationIDStr,
		ID:             tweetID,
		Likes:          tweet.FavoriteCount,
		Name:           name,
		PermanentURL:   fmt.Sprintf("https://x.com/%s/status/%s", username, tweetID),
		Replies:        tweet.ReplyCount,
		Retweets:       tweet.RetweetCount,
		Text:           text,
		UserID:         user.ID,
		Username:       username,
	}

	tm, err := time.Parse(time.RubyDate, tweet.CreatedAt)
	if err == nil {
		tw.TimeParsed = tm
		tw.Timestamp = tm.Unix()
	}

	if tweet.Place.ID != "" {
		tw.Place = &tweet.Place
	}

	if tweet.QuotedStatusIDStr != "" {
		tw.IsQuoted = true
		tw.QuotedStatusID = tweet.QuotedStatusIDStr
	}
	if tweet.InReplyToStatusIDStr != "" {
		tw.IsReply = true
		tw.InReplyToStatusID = tweet.InReplyToStatusIDStr
	}
	if tweet.RetweetedStatusIDStr != "" || tweet.RetweetedStatusResult.Result != nil {
		tw.IsRetweet = true
		tw.RetweetedStatusID = tweet.RetweetedStatusIDStr
		if tweet.RetweetedStatusResult.Result != nil {
			legacy := &tweet.RetweetedStatusResult.Result.Legacy
			userRT := &tweet.RetweetedStatusResult.Result.Core.UserResults.Result
			if tweet.RetweetedStatusResult.Result.Typename == "TweetWithVisibilityResults" {
				legacy = &tweet.RetweetedStatusResult.Result.Tweet.Legacy
				userRT = &tweet.RetweetedStatusResult.Result.Tweet.Core.UserResults.Result
			}
			tw.RetweetedStatus = parseLegacyTweet(userRT, legacy)
			tw.RetweetedStatusID = tw.RetweetedStatus.ID
		}
	}

	if tweet.Views.Count != "" {
		views, viewsErr := strconv.Atoi(tweet.Views.Count)
		if viewsErr != nil {
			views = 0
		}
		tw.Views = views
	}

	for _, pinned := range user.Legacy.PinnedTweetIDsStr {
		if tweet.IDStr == pinned {
			tw.IsPin = true
			break
		}
	}

	for _, hash := range tweet.Entities.Hashtags {
		tw.Hashtags = append(tw.Hashtags, hash.Text)
	}

	for _, mention := range tweet.Entities.UserMentions {
		tw.Mentions = append(tw.Mentions, dtos.Mention{
			ID:       mention.IDStr,
			Username: mention.ScreenName,
			Name:     mention.Name,
		})
	}

	for _, media := range tweet.ExtendedEntities.Media {
		switch media.Type {
		case "photo":
			photo := dtos.Photo{
				ID:  media.IDStr,
				URL: media.MediaURLHttps,
			}

			tw.Photos = append(tw.Photos, photo)
		case "video":
			video := dtos.Video{
				ID:      media.IDStr,
				Preview: media.MediaURLHttps,
			}

			maxBitrate := 0
			for _, variant := range media.VideoInfo.Variants {
				if variant.Type == "application/x-mpegURL" {
					video.HLSURL = variant.URL
				}
				if variant.Bitrate > maxBitrate {
					video.URL = strings.TrimSuffix(variant.URL, "?tag=10")
					maxBitrate = variant.Bitrate
				}
			}

			tw.Videos = append(tw.Videos, video)
		case "animated_gif":
			gif := dtos.GIF{
				ID:      media.IDStr,
				Preview: media.MediaURLHttps,
			}

			// Twitter's API doesn't provide bitrate for GIFs, (it's always set to zero).
			// Therefore we check for `>=` instead of `>` in the loop below.
			// Also, GIFs have just a single variant today. Just in case that changes in the future,
			// and there will be multiple variants, we'll pick the one with the highest bitrate,
			// if other one will have a non-zero bitrate.
			maxBitrate := 0
			for _, variant := range media.VideoInfo.Variants {
				if variant.Bitrate >= maxBitrate {
					gif.URL = variant.URL
					maxBitrate = variant.Bitrate
				}
			}

			tw.GIFs = append(tw.GIFs, gif)
		}

		if !tw.SensitiveContent {
			sensitive := media.ExtSensitiveMediaWarning
			tw.SensitiveContent = sensitive.AdultContent || sensitive.GraphicViolence || sensitive.Other
		}
	}

	for _, url := range tweet.Entities.URLs {
		tw.URLs = append(tw.URLs, url.ExpandedURL)
	}

	tw.HTML = tweet.FullText
	tw.HTML = reHashtag.ReplaceAllStringFunc(tw.HTML, func(hashtag string) string {
		return fmt.Sprintf(`<a href="https://twitter.com/hashtag/%s">%s</a>`,
			strings.TrimPrefix(hashtag, "#"),
			hashtag,
		)
	})
	tw.HTML = reUsername.ReplaceAllStringFunc(tw.HTML, func(username string) string {
		return fmt.Sprintf(`<a href="https://twitter.com/%s">%s</a>`,
			strings.TrimPrefix(username, "@"),
			username,
		)
	})
	var foundedMedia []string
	tw.HTML = reTwitterURL.ReplaceAllStringFunc(tw.HTML, func(tco string) string {
		for _, entity := range tweet.Entities.URLs {
			if tco == entity.URL {
				return fmt.Sprintf(`<a href="%s">%s</a>`, entity.ExpandedURL, entity.ExpandedURL)
			}
		}
		for _, entity := range tweet.ExtendedEntities.Media {
			if tco == entity.URL {
				foundedMedia = append(foundedMedia, entity.MediaURLHttps)
				return fmt.Sprintf(`<br><a href="%s"><img src="%s"/></a>`, tco, entity.MediaURLHttps)
			}
		}
		return tco
	})
	for _, photo := range tw.Photos {
		url := photo.URL
		if stringInSlice(url, foundedMedia) {
			continue
		}
		tw.HTML += fmt.Sprintf(`<br><img src="%s"/>`, url)
	}
	for _, video := range tw.Videos {
		url := video.Preview
		if stringInSlice(url, foundedMedia) {
			continue
		}
		tw.HTML += fmt.Sprintf(`<br><img src="%s"/>`, url)
	}
	for _, gif := range tw.GIFs {
		url := gif.Preview
		if stringInSlice(url, foundedMedia) {
			continue
		}
		tw.HTML += fmt.Sprintf(`<br><img src="%s"/>`, url)
	}
	tw.HTML = strings.ReplaceAll(tw.HTML, "\n", "<br>")
	return tw
}

func mapToJSONString(data any) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Warn().Err(err).Msg("Returning empty string")
		return ""
	}
	return string(jsonBytes)
}

func expandURLs(text string, urls []dtos.URL, extendedMediaEntities []dtos.ExtendedMedia) string {
	expandedText := text
	for _, url := range urls {
		expandedText = strings.ReplaceAll(expandedText, url.URL, url.ExpandedURL)
	}
	for _, entity := range extendedMediaEntities {
		expandedText = strings.ReplaceAll(expandedText, entity.URL, entity.MediaURLHttps)
	}

	return expandedText
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
