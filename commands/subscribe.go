package commands

import (
	"fmt"
	"streambot/twitch"

	"streambot/db"
	"streambot/db/models"
)

var CommandSubscribe string
var CommandSpeedrun string

func init() {
	CommandSubscribe = "subscribe"
	CommandSpeedrun = "speedrun"

	if Config.Debug {
		CommandSubscribe = "dsubscribe"
		CommandSpeedrun = "dspeedrun"
	}
}

func Subscribe(c *Command) {
	c.Reply("Sure! One moment while I look up that game.")

	gameID, err := twitch.GameID(c.RawArguments)
	if err != nil || gameID == "" {
		c.Reply(fmt.Sprintf("Unexpected error: %v", err))
		return
	}

	reservation := models.Reservation{
		GuildID:      c.Event.GuildID,
		ChannelID:    c.Event.ChannelID,
		GameID:       gameID,
		SpeedrunOnly: c.Name == CommandSpeedrun,
	}
	db.Conn.Create(&reservation)

	c.Reply("Subscribed!")
}
