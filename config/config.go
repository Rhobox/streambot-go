package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
)

type AppConfig struct {
	TwitchClientID     string
	TwitchClientSecret string
	DiscordBotToken    string
	Debug              bool
}

var Config *AppConfig

func init() {
	Config = &AppConfig{}
	Config.TwitchClientID = os.Getenv("TWITCH_CLIENT_ID")
	Config.TwitchClientSecret = os.Getenv("TWITCH_CLIENT_SECRET")
	Config.DiscordBotToken = os.Getenv("DISCORD_BOT_TOKEN")
	Config.Debug = os.Getenv("DEBUG") != ""
}
