package twitch

import (
	"time"

	"github.com/nicklaw5/helix"
	"streambot/config"
)

var conf = config.Config
var Helix *helix.Client

func getAuthToken() {
	resp, err := Helix.RequestAppAccessToken([]string{})
	if err != nil {
		panic(err)
	}

	Helix.SetAppAccessToken(resp.Data.AccessToken)
}

func init() {
	if conf.TwitchClientID == "" || conf.TwitchClientSecret == "" {
		panic("Missing either TWITCH_CLIENT_ID or TWITCH_CLIENT_SECRET")
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:     conf.TwitchClientID,
		ClientSecret: conf.TwitchClientSecret,
	})
	if err != nil {
		panic(err)
	}

	Helix = client
	getAuthToken()

	// This library has no endpoint for refreshing the token
	// so we just spin up a background task to get a new one
	// every now and then
	go (func() {
		ticker := time.Tick(5 * time.Second)

		for range ticker {
			getAuthToken()
		}
	})()
}
