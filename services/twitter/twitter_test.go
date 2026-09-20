package twitter

import (
	"testing"
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-twitter/models/dtos"
	"github.com/kaellybot/kaelly-twitter/models/entities"
)

type stubRepository struct {
	accounts []entities.TwitterAccount
	saved    *[]entities.TwitterAccount
}

func (stub stubRepository) GetTwitterAccounts() ([]entities.TwitterAccount, error) {
	return stub.accounts, nil
}

func (stub stubRepository) Save(account entities.TwitterAccount) error {
	*stub.saved = append(*stub.saved, account)
	return nil
}

func newTestService(accounts []entities.TwitterAccount, saved *[]entities.TwitterAccount,
	published *[]*amqp.RabbitMQMessage) *Impl {
	broker := amqp.Mock{
		EmitFunc: func(msg *amqp.RabbitMQMessage, _ amqp.Exchange, _, _ string) error {
			*published = append(*published, msg)
			return nil
		},
	}
	return &Impl{
		broker:              &broker,
		twitterAccountsRepo: stubRepository{accounts: accounts, saved: saved},
	}
}

// News carries the account's game and locale, so the notifier posts it with the
// right bot.
func TestPublishTweetCarriesAccountGame(t *testing.T) {
	t.Parallel()

	for _, game := range []amqp.Game{amqp.Game_DOFUS_GAME, amqp.Game_DOFUS_TOUCH} {
		t.Run(game.String(), func(t *testing.T) {
			t.Parallel()

			saved := make([]entities.TwitterAccount, 0)
			published := make([]*amqp.RabbitMQMessage, 0)
			service := newTestService(nil, &saved, &published)

			err := service.publishTweet(entities.TwitterAccount{
				ID: "account", Name: "account", Game: game, Locale: amqp.Language_FR,
			}, &dtos.Tweet{ID: "tweet1", Timestamp: 100, Text: "text"})
			if err != nil {
				t.Fatalf("publishTweet() = %v, want nil", err)
			}

			if len(published) != 1 {
				t.Fatalf("published %d news, want 1", len(published))
			}
			if published[0].GetGame() != game {
				t.Errorf("Game = %v, want %v", published[0].GetGame(), game)
			}
			if !published[0].IsGameSet() {
				t.Error("news game must be set, otherwise the broker refuses to publish it")
			}
			if published[0].GetLanguage() != amqp.Language_FR {
				t.Errorf("Language = %v, want %v", published[0].GetLanguage(), amqp.Language_FR)
			}
			if published[0].GetType() != amqp.RabbitMQMessage_NEWS_TWITTER {
				t.Errorf("Type = %v, want %v", published[0].GetType(),
					amqp.RabbitMQMessage_NEWS_TWITTER)
			}
		})
	}
}

// A row with no game can never be published, so it must not cost a Twitter call.
func TestDispatchNewTweetsSkipsAccountsWithoutGame(t *testing.T) {
	t.Parallel()

	saved := make([]entities.TwitterAccount, 0)
	published := make([]*amqp.RabbitMQMessage, 0)
	service := newTestService([]entities.TwitterAccount{
		{ID: "account", Game: amqp.Game_ANY_GAME, LastUpdate: time.Unix(0, 0).UTC()},
	}, &saved, &published)

	// fetchTweets would need a live HTTP client; reaching it at all is the failure.
	if err := service.DispatchNewTweets(); err != nil {
		t.Fatalf("DispatchNewTweets() = %v, want nil", err)
	}

	if len(published) != 0 {
		t.Errorf("published %d news, want 0", len(published))
	}
	if len(saved) != 0 {
		t.Errorf("saved %d account(s), want 0", len(saved))
	}
}
