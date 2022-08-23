package discord

import (
	"github.com/aricodes-oss/std"
	"github.com/bwmarrin/discordgo"
	"streambot/config"
)

var Session *discordgo.Session
var logger = std.Logger
var appConfig = config.Config

func init() {
	if appConfig.DiscordBotToken == "" {
		panic("DISCORD_BOT_TOKEN environment variable not set")
	}

	s, err := discordgo.New("Bot " + appConfig.DiscordBotToken)
	if err != nil {
		panic(err)
	}
	Session = s
}
