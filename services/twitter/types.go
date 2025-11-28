package twitter

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-twitter/repositories/twitteraccounts"
)

const (
	routingkey = "news.twitter"
	baseURL    = "https://x.com/i/api/graphql/lZRf8IC-GTuGxDwcsHW8aw/UserTweets"

	//nolint:lll // Clearer as it is.
	bearerToken = "AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"
)

type Service interface {
	DispatchNewTweets() error
}

type Impl struct {
	userAgent           string
	authToken           string
	csrfToken           string
	tweetCount          int
	broker              amqp.MessageBroker
	twitterAccountsRepo twitteraccounts.Repository
}
