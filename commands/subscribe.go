package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"streambot/twitch"

	"streambot/db"
	"streambot/db/models"
	"streambot/workers"
)

var (
	CommandSubscribe string
	CommandSpeedrun  string
)

func init() {
	CommandSubscribe = "subscribe"
	CommandSpeedrun = "speedrun"

	if Config.Debug {
		CommandSubscribe = "dsubscribe"
		CommandSpeedrun = "dspeedrun"
	}
}

func userCanManageMessages(c *Command) bool {
	perms, err := c.Session.UserChannelPermissions(c.Session.State.User.ID, c.Event.ChannelID)
	if err != nil {
		return false
	}

	return perms&discordgo.PermissionManageMessages != 0
}

func Subscribe(c *Command) {
	c.Reply("Sure! One moment while I look up that game.")

	gameID, err := twitch.GameID(c.RawArguments)
	if err != nil || gameID == "" {
		c.Reply(fmt.Sprintf("Unexpected error: %v", err))
		return
	}

	if !userCanManageMessages(c) {
		c.Reply("User does not have permission to manage messages. This permission is required to function.")
		return
	}

	existing := models.Reservation{}
	db.Conn.Where(&models.Reservation{GameID: gameID, GuildID: c.Event.GuildID, ChannelID: c.Event.ChannelID}).First(&existing)

	if existing.ID != 0 {
		c.Reply(fmt.Sprintf("This channel is already subscribed to %v streams", c.RawArguments))
	}

	reservation := models.Reservation{
		GuildID:      c.Event.GuildID,
		ChannelID:    c.Event.ChannelID,
		GameID:       gameID,
		Name:         c.RawArguments,
		SpeedrunOnly: c.Name == CommandSpeedrun,
	}
	db.Conn.Create(&reservation)

	c.Reply("Subscribed! Gathering streams now, this may take a minute.")
	workers.StreamsWorker()
}
