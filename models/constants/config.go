package constants

import "github.com/rs/zerolog"

const (
	ConfigFileName = ".env"

	// Username used to login to Twitter.
	TwitterUsername = "TWITTER_USERNAME"

	// Password used to login to Twitter.
	TwitterPassword = "TWITTER_PASSWORD"

	// Auth token used when logged in to Twitter.
	TwitterAuthToken = "TWITTER_AUTH_TOKEN"

	// CSRF token used when logged in to Twitter.
	TwitterCSRFToken = "TWITTER_CSRF_TOKEN"

	// Number of tweets retrieved per call.
	TwitterTweetCount = "TWEET_COUNT"

	// Discord webhook URL to publish errors that matter
	DiscordWebhookUrl = "DISCORD_WEBHOOK_URL"

	// Metric port.
	MetricPort = "METRIC_PORT"

	// Zerolog values from [trace, debug, info, warn, error, fatal, panic].
	LogLevel = "LOG_LEVEL"

	// Boolean; used to register commands at development guild level or globally.
	Production = "PRODUCTION"

	defaultTwitterUsername   = ""
	defaultTwitterPassword   = ""
	defaultTwitterAuthToken  = ""
	defaultTwitterCSRFToken  = ""
	defaultTwitterTweetCount = 20
	defaultDiscordWebhookUrl = ""
	defaultMetricPort        = 2112
	defaultLogLevel          = zerolog.InfoLevel
	defaultProduction        = false
)

func GetDefaultConfigValues() map[string]any {
	return map[string]any{
		TwitterUsername:   defaultTwitterUsername,
		TwitterPassword:   defaultTwitterPassword,
		TwitterAuthToken:  defaultTwitterAuthToken,
		TwitterCSRFToken:  defaultTwitterCSRFToken,
		TwitterTweetCount: defaultTwitterTweetCount,
		DiscordWebhookUrl: defaultDiscordWebhookUrl,
		MetricPort:        defaultMetricPort,
		LogLevel:          defaultLogLevel.String(),
		Production:        defaultProduction,
	}
}
