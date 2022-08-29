package workers

import (
	"streambot/channels"
	"streambot/db"
	"streambot/db/models"
	"sync"
	"time"

	"gorm.io/gorm"
)

func scrub_db(rid uint) {
	allReservations := []models.Reservation{}
	db.Conn.Preload("Streams").Find(&allReservations)

	db.Conn.Transaction(func(tx *gorm.DB) error {
		for _, res := range allReservations {
			for _, stream := range res.Streams {
				_, err := discordClient.ChannelMessage(res.ChannelID, stream.MessageID)

				if err != nil {
					tx.Unscoped().Delete(stream)
				}
			}
		}

		return nil
	})
}

func ScrubDBWorker(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-channels.Running:
			return
		case <-ticker.C:
			clientLock.Lock()
			reservations := []models.Reservation{}
			db.Conn.Find(&reservations)

			// Running these operations synchronously on purpose
			for _, reservation := range reservations {
				log.Debugf("Scrubbing messages table for %v", reservation.ID)
				scrub_db(reservation.ID)
			}
			clientLock.Unlock()
		}
	}
}
