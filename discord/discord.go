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
	s, err := discordgo.New("Bot " + appConfig.DiscordBotToken)
	if err != nil {
		panic(err)
	}
	Session = s
}
