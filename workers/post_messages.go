package workers

import (
	"streambot/channels"
	"streambot/db"
	"streambot/db/models"
	"sync"
	"time"
)

func post_messages(rid uint) {
	reservation := models.Reservation{}
	live_streams := []models.TwitchStream{}

	// Check to make sure we actually have this reservation
	result := db.Conn.Preload("Streams").Where("id = ?", rid).First(&reservation)
	if result.RowsAffected == 0 {
		log.Debug("Bailed at finding res")
		return
	}

	known_users := make([]string, len(reservation.Streams))
	for idx, val := range reservation.Streams {
		known_users[idx] = val.Username
	}

	// Get the list of streams preloaded by another worker into the database
	query := db.Conn.Where("game_id = ?", reservation.GameID)
	if len(known_users) > 0 {
		query = query.Where("username NOT IN ?", known_users)
	}
	query.Find(&live_streams)

	if len(live_streams) == 0 {
		log.Debug("Bailed at finding streams")
		return
	}

	// Iterate through individually so that the database stays in sync with the channel
	for _, stream := range live_streams {
		msg, err := discordClient.ChannelMessageSend(reservation.ChannelID, stream.Username)
		if err != nil {
			log.Warn(err)
			continue
		}

		message_record := models.Stream{
			Username:      stream.Username,
			ReservationID: reservation.ID,
			MessageID:     msg.ID,
		}
		db.Conn.Create(&message_record)
	}
}

func PostMessagesWorker(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	subgroup := sync.WaitGroup{}

	for {
		select {
		case <-channels.Running:
			return
		case <-ticker.C:
			clientLock.Lock()
			reservations := []models.Reservation{}
			db.Conn.Find(&reservations)

			for _, reservation := range reservations {
				log.Debugf("Posting new messages for reservation ID %v", reservation.ID)
				go func(rid uint) {
					subgroup.Add(1)
					defer subgroup.Done()

					post_messages(rid)
				}(reservation.ID)
			}

			subgroup.Wait()
			clientLock.Unlock()
		}
	}
}
