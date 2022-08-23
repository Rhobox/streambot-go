package main

import (
	"github.com/aricodes-oss/std"
	"github.com/bwmarrin/discordgo"

	"streambot/config"
	"streambot/discord"
	"time"
)

var logger = std.Logger
var session = discord.Session

func main() {
	session.AddHandler(messageCreate)
	session.Identify.Intents = discordgo.IntentsAll

	err := session.Open()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	logger.Info("Discord connected successfully!")

	if config.Config.Debug {
		logger.Warn("Running with debug enabled, be careful!")
		logger = logger.WithDebug()
	}

	std.WaitForKill()
}

func messageIsCommand(message string) bool {
	return len(message) >= 2 && string(message[0]) == "!"
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Only process commands from external users and not any other messages
	if m.Author.ID == s.State.User.ID || !messageIsCommand(m.Content) {
		return
	}

	logger.Debug(m.Content)

	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "!sleep" {
		time.Sleep(time.Second * 5)
		s.ChannelMessageSend(m.ChannelID, "Slept!")
	}
}
