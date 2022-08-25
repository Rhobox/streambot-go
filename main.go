package main

import (
	"github.com/aricodes-oss/std"
	"github.com/bwmarrin/discordgo"

	"streambot/commands"
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Only process commands from external users and not any other messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	command, err := commands.Parse(s, m)
	if err != nil {
		logger.Debugf("Error parsing command: %v", err)
		return
	}

	switch command.Name {
	case commands.CommandSubscribe, commands.CommandSpeedrun:
		commands.Subscribe(command)
	case commands.CommandUnsubscribe:
		commands.Unsubscribe(command)
	case "ping":
		command.Reply("Pong!")
	case "sleep":
		time.Sleep(time.Second * 5)
		command.Reply("Slept!")
	}
}
