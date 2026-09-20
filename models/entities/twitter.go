package entities

import (
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
)

// TwitterAccount is identified by (id, game): the same account may be followed in
// several games, each with its own locale and its own last_update cursor.
type TwitterAccount struct {
	ID         string    `gorm:"primaryKey"`
	Game       amqp.Game `gorm:"primaryKey"`
	Name       string
	Locale     amqp.Language
	LastUpdate time.Time `gorm:"not null; default:current_timestamp"`
}
