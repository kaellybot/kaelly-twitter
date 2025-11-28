package dtos

import "time"

// Tweet type.
type Tweet struct {
	ConversationID    string
	GIFs              []GIF
	Hashtags          []string
	HTML              string
	ID                string
	InReplyToStatus   *Tweet
	InReplyToStatusID string
	IsQuoted          bool
	IsPin             bool
	IsReply           bool
	IsRetweet         bool
	IsSelfThread      bool
	Likes             int
	Name              string
	Mentions          []Mention
	PermanentURL      string
	Photos            []Photo
	Place             *Place
	QuotedStatus      *Tweet
	QuotedStatusID    string
	Replies           int
	Retweets          int
	RetweetedStatus   *Tweet
	RetweetedStatusID string
	Text              string
	Thread            []*Tweet
	TimeParsed        time.Time
	Timestamp         int64
	URLs              []string
	UserID            string
	Username          string
	Videos            []Video
	Views             int
	SensitiveContent  bool
}

// Mention type.
type Mention struct {
	ID       string
	Username string
	Name     string
}

// URL represents a URL with display, expanded, and index data.
type URL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	URL         string `json:"url"`
	Indices     []int  `json:"indices"`
}

// Photo type.
type Photo struct {
	ID  string
	URL string
}

// Video type.
type Video struct {
	ID      string
	Preview string
	URL     string
	HLSURL  string
}

// GIF type.
type GIF struct {
	ID      string
	Preview string
	URL     string
}

type Place struct {
	ID          string `json:"id"`
	PlaceType   string `json:"place_type"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	CountryCode string `json:"country_code"`
	Country     string `json:"country"`
	BoundingBox struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"bounding_box"`
}

type TimelineV2 struct {
	Data struct {
		User struct {
			Result struct {
				TimelineV2 struct {
					Timeline struct {
						Instructions []struct {
							ModuleItems []Item  `json:"moduleItems"`
							Entries     []Entry `json:"entries"`
							Entry       Entry   `json:"entry"`
							Type        string  `json:"type"`
						} `json:"instructions"`
					} `json:"timeline"`
				} `json:"timeline_v2"`

				Timeline struct {
					Timeline struct {
						Instructions []struct {
							Entries []Entry `json:"entries"`
							Entry   Entry   `json:"entry"`
							Type    string  `json:"type"`
						} `json:"instructions"`
					} `json:"timeline"`
				} `json:"timeline"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
}

type Item struct {
	EntryID string `json:"entryId"`
	Item    struct {
		ItemContent struct {
			ItemType         string `json:"itemType"`
			TweetDisplayType string `json:"tweetDisplayType"`
			TweetResults     struct {
				Result Result `json:"result"`
			} `json:"tweet_results"`
			CursorType string `json:"cursorType"`
			Value      string `json:"value"`
		} `json:"itemContent"`
	} `json:"item"`
}

type Entry struct {
	Content struct {
		CursorType  string `json:"cursorType"`
		Value       string `json:"value"`
		Items       []Item `json:"items"`
		ItemContent struct {
			ItemType         string `json:"itemType"`
			TweetDisplayType string `json:"tweetDisplayType"`
			TweetResults     struct {
				Result Result `json:"result"`
			} `json:"tweet_results"`
			UserDisplayType string `json:"userDisplayType"`
			UserResults     struct {
				Result UserResult `json:"result"`
			} `json:"user_results"`
			CursorType string `json:"cursorType"`
			Value      string `json:"value"`
		} `json:"itemContent"`
	} `json:"content"`
}

type Result struct {
	Typename string `json:"__typename"`
	RawTweet
	Tweet RawTweet `json:"tweet"`
}

type RawTweet struct {
	Core struct {
		UserResults struct {
			Result LegacyUserResult `json:"result"`
		} `json:"user_results"`
	} `json:"core"`
	Views struct {
		Count string `json:"count"`
	} `json:"views"`
	NoteTweet struct {
		NoteTweetResults struct {
			Result struct {
				Text string `json:"text"`
			} `json:"result"`
		} `json:"note_tweet_results"`
	} `json:"note_tweet"`
	QuotedStatusResult struct {
		Result *Result `json:"result"`
	} `json:"quoted_status_result"`
	Legacy LegacyRawTweet `json:"legacy"`
	Card   struct {
		RestID string `json:"rest_id"`
		Legacy struct {
			BindingValues []struct {
				Key   string `json:"key"`
				Value struct {
					StringValue string `json:"string_value"`
				} `json:"value"`
			} `json:"binding_values"`
		} `json:"legacy"`
	} `json:"card"`
}

type LegacyUserResult struct {
	ID             string `json:"rest_id"`
	IsBlueVerified bool   `json:"is_blue_verified"`
	Core           struct {
		CreatedAt  string `json:"created_at"`
		Name       string `json:"name"`
		ScreenName string `json:"screen_name"`
	} `json:"core"`
	Legacy LegacyUser `json:"legacy"`
}

type LegacyUser struct {
	Description string `json:"description"`
	Entities    struct {
		URL struct {
			Urls []struct {
				ExpandedURL string `json:"expanded_url"`
			} `json:"urls"`
		} `json:"url"`
	} `json:"entities"`
	FavouritesCount         int      `json:"favourites_count"`
	FollowersCount          int      `json:"followers_count"`
	FriendsCount            int      `json:"friends_count"`
	ListedCount             int      `json:"listed_count"`
	Location                string   `json:"location"`
	PinnedTweetIDsStr       []string `json:"pinned_tweet_ids_str"`
	ProfileBannerURL        string   `json:"profile_banner_url"`
	ProfileImageURLHTTPS    string   `json:"profile_image_url_https"`
	Protected               bool     `json:"protected"`
	StatusesCount           int      `json:"statuses_count"`
	Verified                bool     `json:"verified"`
	FollowedBy              bool     `json:"followed_by"`
	Following               bool     `json:"following"`
	CanDm                   bool     `json:"can_dm"`
	CanMediaTag             bool     `json:"can_media_tag"`
	DefaultProfile          bool     `json:"default_profile"`
	DefaultProfileImage     bool     `json:"default_profile_image"`
	FastFollowersCount      int      `json:"fast_followers_count"`
	HasCustomTimelines      bool     `json:"has_custom_timelines"`
	IsTranslator            bool     `json:"is_translator"`
	MediaCount              int      `json:"media_count"`
	NeedsPhoneVerification  bool     `json:"needs_phone_verification"`
	NormalFollowersCount    int      `json:"normal_followers_count"`
	PossiblySensitive       bool     `json:"possibly_sensitive"`
	ProfileInterstitialType string   `json:"profile_interstitial_type"`
	TranslatorType          string   `json:"translator_type"`
	WantRetweets            bool     `json:"want_retweets"`
	WithheldInCountries     []string `json:"withheld_in_countries"`
}

type LegacyUserV2 struct {
	FollowedBy          bool   `json:"followed_by"`
	Following           bool   `json:"following"`
	CanDm               bool   `json:"can_dm"`
	CanMediaTag         bool   `json:"can_media_tag"`
	CreatedAt           string `json:"created_at"`
	DefaultProfile      bool   `json:"default_profile"`
	DefaultProfileImage bool   `json:"default_profile_image"`
	Description         string `json:"description"`
	Entities            struct {
		Description struct {
			Urls []URL `json:"urls"`
		} `json:"description"`
		URL struct {
			Urls []URL `json:"urls"`
		} `json:"url"`
	} `json:"entities"`
	FastFollowersCount      int      `json:"fast_followers_count"`
	FavouritesCount         int      `json:"favourites_count"`
	FollowersCount          int      `json:"followers_count"`
	FriendsCount            int      `json:"friends_count"`
	HasCustomTimelines      bool     `json:"has_custom_timelines"`
	IsTranslator            bool     `json:"is_translator"`
	ListedCount             int      `json:"listed_count"`
	Location                string   `json:"location"`
	MediaCount              int      `json:"media_count"`
	Name                    string   `json:"name"`
	NormalFollowersCount    int      `json:"normal_followers_count"`
	PinnedTweetIDsStr       []string `json:"pinned_tweet_ids_str"`
	PossiblySensitive       bool     `json:"possibly_sensitive"`
	ProfileBannerURL        string   `json:"profile_banner_url"`
	ProfileImageURLHTTPS    string   `json:"profile_image_url_https"`
	ProfileInterstitialType string   `json:"profile_interstitial_type"`
	ScreenName              string   `json:"screen_name"`
	StatusesCount           int      `json:"statuses_count"`
	TranslatorType          string   `json:"translator_type"`
	URL                     string   `json:"url"`
	Verified                bool     `json:"verified"`
	WantRetweets            bool     `json:"want_retweets"`
	WithheldInCountries     []any    `json:"withheld_in_countries"`
}

type LegacyRawTweet struct {
	ConversationIDStr string `json:"conversation_id_str"`
	CreatedAt         string `json:"created_at"`
	FavoriteCount     int    `json:"favorite_count"`
	FullText          string `json:"full_text"`
	Entities          struct {
		Hashtags []struct {
			Text string `json:"text"`
		} `json:"hashtags"`
		Media []struct {
			MediaURLHttps string `json:"media_url_https"`
			Type          string `json:"type"`
			URL           string `json:"url"`
		} `json:"media"`
		URLs         []URL `json:"urls"`
		UserMentions []struct {
			IDStr      string `json:"id_str"`
			Name       string `json:"name"`
			ScreenName string `json:"screen_name"`
		} `json:"user_mentions"`
	} `json:"entities"`
	ExtendedEntities struct {
		Media []ExtendedMedia `json:"media"`
	} `json:"extended_entities"`
	IDStr                 string `json:"id_str"`
	InReplyToStatusIDStr  string `json:"in_reply_to_status_id_str"`
	Place                 Place  `json:"place"`
	ReplyCount            int    `json:"reply_count"`
	RetweetCount          int    `json:"retweet_count"`
	RetweetedStatusIDStr  string `json:"retweeted_status_id_str"`
	RetweetedStatusResult struct {
		Result *Result `json:"result"`
	} `json:"retweeted_status_result"`
	QuotedStatusIDStr string `json:"quoted_status_id_str"`
	SelfThread        struct {
		IDStr string `json:"id_str"`
	} `json:"self_thread"`
	Time      time.Time `json:"time"`
	UserIDStr string    `json:"user_id_str"`
	Views     struct {
		State string `json:"state"`
		Count string `json:"count"`
	} `json:"ext_views"`
}

type ExtendedMedia struct {
	IDStr                    string `json:"id_str"`
	MediaURLHttps            string `json:"media_url_https"`
	ExtSensitiveMediaWarning struct {
		AdultContent    bool `json:"adult_content"`
		GraphicViolence bool `json:"graphic_violence"`
		Other           bool `json:"other"`
	} `json:"ext_sensitive_media_warning"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	VideoInfo struct {
		Variants []struct {
			Type    string `json:"content_type"`
			Bitrate int    `json:"bitrate"`
			URL     string `json:"url"`
		} `json:"variants"`
	} `json:"video_info"`
}

type UserResult struct {
	Typename                   string                `json:"__typename"`
	ID                         string                `json:"id"`
	RestID                     string                `json:"rest_id"`
	AffiliatesHighlightedLabel struct{}              `json:"affiliates_highlighted_label"`
	HasGraduatedAccess         bool                  `json:"has_graduated_access"`
	IsBlueVerified             bool                  `json:"is_blue_verified"`
	ProfileImageShape          string                `json:"profile_image_shape"`
	Legacy                     LegacyUserV2          `json:"legacy"`
	LegacyExtendedProfile      LegacyExtendedProfile `json:"legacy_extended_profile"`
	IsProfileTranslatable      bool                  `json:"is_profile_translatable"`
	VerificationInfo           VerificationInfo      `json:"verification_info"`
	HighlightsInfo             HighlightsInfo        `json:"highlights_info"`
	UserSeedTweetCount         int                   `json:"user_seed_tweet_count"`
	PremiumGiftingEligible     bool                  `json:"premium_gifting_eligible"`
	CreatorSubscriptionsCount  int                   `json:"creator_subscriptions_count"`
}

type LegacyExtendedProfile struct {
	Birthdate struct {
		Day            int    `json:"day"`
		Month          int    `json:"month"`
		Year           int    `json:"year"`
		Visibility     string `json:"visibility"`
		YearVisibility string `json:"year_visibility"`
	} `json:"birthdate"`
}

type VerificationInfo struct {
	IsIdentityVerified bool `json:"is_identity_verified"`
	Reason             struct {
		Description struct {
			Text     string `json:"text"`
			Entities []struct {
				FromIndex int `json:"from_index"`
				ToIndex   int `json:"to_index"`
				Ref       struct {
					URL     string `json:"url"`
					URLType string `json:"url_type"`
				} `json:"ref"`
			} `json:"entities"`
		} `json:"description"`
		VerifiedSinceMsec string `json:"verified_since_msec"`
	} `json:"reason"`
}

type HighlightsInfo struct {
	CanHighlightTweets bool   `json:"can_highlight_tweets"`
	HighlightedTweets  string `json:"highlighted_tweets"`
}

type UnifiedCard struct {
	Type          string   `json:"type"`
	Components    []string `json:"components"`
	MediaEntities map[string]struct {
		ID            int64  `json:"id"`
		IDStr         string `json:"id_str"`
		MediaURLHTTPS string `json:"media_url_https"`
		Type          string `json:"type"`
		OriginalInfo  struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"original_info"`
		SourceUserID int64 `json:"source_user_id"`
		VideoInfo    struct {
			AspectRatio    []int `json:"aspect_ratio"`
			DurationMillis int   `json:"duration_millis"`
			Variants       []struct {
				Bitrate     int    `json:"bitrate,omitempty"`
				ContentType string `json:"content_type"`
				URL         string `json:"url"`
			} `json:"variants"`
		} `json:"video_info"`
	} `json:"media_entities"`
}
