package models

import (
	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model
	GuildID      string
	ChannelID    string
	GameID       string
	SpeedrunOnly bool
	Streams      []Stream
}
