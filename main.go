package main

import (
	"github.com/aricodes-oss/std"
	"github.com/bwmarrin/discordgo"

	"streambot/channels"
	"streambot/commands"
	"streambot/config"
	"streambot/discord"
	"streambot/workers"

	"sync"
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

	wg := &sync.WaitGroup{}
	workers.LaunchAll(wg)

	std.WaitForKill()
	logger.Info("Kill signal received, shutting down workers...")

	close(channels.Running)
	wg.Wait()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Only process commands from external users and not any other messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	command, err := commands.Parse(s, m)
	if err != nil {
		if err != commands.ErrNotACommand {
			logger.Debugf("Error parsing command: %v", err)
		}

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
