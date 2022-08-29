package workers

import (
	"fmt"
	"streambot/db"
	"streambot/db/models"
	"strings"
	"sync"

	"github.com/aricodes-oss/std"
	"github.com/bwmarrin/discordgo"
	"github.com/clinet/discordgo-embed"
)

func make_embed(stream *models.TwitchStream) *discordgo.MessageEmbed {
	thumbnailUrl := strings.Replace(stream.ThumbnailURL, "{width}", "72", 1)
	thumbnailUrl = strings.Replace(thumbnailUrl, "{height}", "72", 1)

	return embed.NewEmbed().
		SetTitle(stream.DisplayName).
		SetDescription(stream.Description).
		SetURL(fmt.Sprintf("https://twitch.tv/%s", stream.Username)).
		SetColor(0x0099FF).
		SetThumbnail(thumbnailUrl).MessageEmbed
}

func post_messages(rid uint) {
	reservation := models.Reservation{}
	live_streams := []models.TwitchStream{}

	// Check to make sure we actually have this reservation
	result := db.Conn.Preload("Streams").Where("id = ?", rid).First(&reservation)
	if result.RowsAffected == 0 {
		return
	}

	known_users := *std.Map(reservation.Streams, func(val models.Stream, idx int) string { return val.Username })

	// Get the list of streams preloaded by another worker into the database
	query := db.Conn.Where("game_id = ?", reservation.GameID)
	if len(known_users) > 0 {
		query = query.Where("username NOT IN ?", known_users)
	}

	// Filter on the speedrun tag if necessary
	if reservation.SpeedrunOnly {
		query = query.Where(&models.TwitchStream{Speedrun: true})
	}
	query.Find(&live_streams)

	if len(live_streams) == 0 {
		return
	}

	// Iterate through individually so that the database stays in sync with the channel
	for _, stream := range live_streams {
		embed := make_embed(&stream)
		msg, err := discordClient.ChannelMessageSendEmbed(reservation.ChannelID, embed)
		if err != nil {
			log.Warn(err)
			continue
		}

		message_record := models.Stream{
			Username:      stream.Username,
			ReservationID: reservation.ID,
			MessageID:     msg.ID,
			GameID:        reservation.GameID,
		}
		db.Conn.Create(&message_record)
	}
}

func PostMessagesWorker() {
	subgroup := sync.WaitGroup{}
	reservations := []models.Reservation{}
	db.Conn.Find(&reservations)

	for _, reservation := range reservations {
		log.Debugf("Posting new messages for reservation ID %v", reservation.ID)
		subgroup.Add(1)
		go func(rid uint) {
			defer subgroup.Done()

			post_messages(rid)
		}(reservation.ID)
	}

	subgroup.Wait()
}
