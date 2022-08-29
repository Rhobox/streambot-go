package commands

import (
	"fmt"
	"streambot/twitch"
	"streambot/workers"

	"streambot/db"
	"streambot/db/models"
)

var CommandUnsubscribe string

func init() {
	CommandUnsubscribe = "unsubscribe"

	if Config.Debug {
		CommandUnsubscribe = "dunsubscribe"
	}
}

func Unsubscribe(c *Command) {
	reservation := models.Reservation{}

	gameID, err := twitch.GameID(c.RawArguments)
	if err != nil || gameID == "" {
		c.Reply(fmt.Sprintf("Unexpected error: %v", err))
		return
	}

	db.Conn.Where(&models.Reservation{
		GuildID: c.Event.GuildID,
		GameID:  gameID,
	}).First(&reservation)

	if reservation.ID == 0 {
		c.Reply("No matching reservation found.")
		return
	}

	// Clear out the channel before it falls out of scope of the reservation
	db.Conn.Unscoped().Where(&models.Stream{ReservationID: reservation.ID}).Delete(&models.Stream{})
	workers.CleanChannelsWorker()

	db.Conn.Delete(&reservation)

	c.Reply("Unsubscribed!")
}
