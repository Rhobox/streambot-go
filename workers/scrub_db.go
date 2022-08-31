package workers

import (
	"streambot/db"
	"streambot/db/models"
	"sync"

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

func ScrubDBWorker() {
	subgroup := sync.WaitGroup{}

	reservations := []models.Reservation{}
	db.Conn.Find(&reservations)

	// Running these operations synchronously on purpose
	for _, reservation := range reservations {
		log.Debugf("Scrubbing messages table for %v", reservation.ID)
		subgroup.Add(1)

		go func(rid uint) {
			defer subgroup.Done()

			scrub_db(rid)
		}(reservation.ID)
	}

	subgroup.Wait()
}
