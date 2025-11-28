package constants

import "github.com/rs/zerolog"

const (
	ConfigFileName = ".env"

	// MySQL URL with the following format: HOST:PORT.
	MySQLURL = "MYSQL_URL"

	// MySQL user.
	MySQLUser = "MYSQL_USER"

	// MySQL password.
	MySQLPassword = "MYSQL_PASSWORD"

	// MySQL database name.
	MySQLDatabase = "MYSQL_DATABASE"

	// RabbitMQ address.
	RabbitMQAddress = "RABBITMQ_ADDRESS"

	//nolint:gosec // False positive.
	// Auth token used when logged in to Twitter.
	TwitterAuthToken = "TWITTER_AUTH_TOKEN"

	//nolint:gosec // False positive.
	// CSRF token used when logged in to Twitter.
	TwitterCSRFToken = "TWITTER_CSRF_TOKEN"

	// Number of tweets retrieved per call.
	TwitterTweetCount = "TWEET_COUNT"

	// User Agent to call Twitter API.
	TwitterUserAgent = "USER_AGENT"

	// Metric port.
	MetricPort = "METRIC_PORT"

	// Zerolog values from [trace, debug, info, warn, error, fatal, panic].
	LogLevel = "LOG_LEVEL"

	// Boolean; used to register commands at development guild level or globally.
	Production = "PRODUCTION"

	defaultMySQLURL          = "localhost:3306"
	defaultMySQLUser         = ""
	defaultMySQLPassword     = ""
	defaultMySQLDatabase     = "kaellybot"
	defaultRabbitMQAddress   = "amqp://localhost:5672"
	defaultTwitterAuthToken  = ""
	defaultTwitterCSRFToken  = ""
	defaultTwitterTweetCount = 20
	defaultTwitterUserAgent  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0"
	defaultMetricPort        = 2112
	defaultLogLevel          = zerolog.InfoLevel
	defaultProduction        = false
)

func GetDefaultConfigValues() map[string]any {
	return map[string]any{
		MySQLURL:          defaultMySQLURL,
		MySQLUser:         defaultMySQLUser,
		MySQLPassword:     defaultMySQLPassword,
		MySQLDatabase:     defaultMySQLDatabase,
		RabbitMQAddress:   defaultRabbitMQAddress,
		TwitterAuthToken:  defaultTwitterAuthToken,
		TwitterCSRFToken:  defaultTwitterCSRFToken,
		TwitterTweetCount: defaultTwitterTweetCount,
		TwitterUserAgent:  defaultTwitterUserAgent,
		MetricPort:        defaultMetricPort,
		LogLevel:          defaultLogLevel.String(),
		Production:        defaultProduction,
	}
}
