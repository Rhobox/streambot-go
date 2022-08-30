package workers

import (
	"streambot/constants"
	"streambot/db"
	"streambot/db/models"
	"streambot/twitch"

	"sync"

	"gorm.io/gorm"
)

func fetch_streams(gameID string) {
	streams, err := twitch.Streams(gameID)
	if err != nil {
		log.Warnf("Failed to fetch streams for %v", gameID)
		return
	}

	records := make([]models.TwitchStream, len(streams))

	for idx, val := range streams {
		records[idx] = models.TwitchStream{
			Username:     val.UserLogin,
			DisplayName:  val.UserName,
			GameID:       gameID,
			Description:  val.Title,
			ThumbnailURL: val.ThumbnailURL,
			Speedrun:     false,
		}

		for _, tagID := range val.TagIDs {
			if tagID == constants.SPEEDRUN_TAG_ID {
				records[idx].Speedrun = true
				break
			}
		}
	}

	// Batch full delete, then batch insert in the same transaction
	db.Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("game_id = ?", gameID).Delete(&models.TwitchStream{}).Error; err != nil {
			log.Warnf("Failed to batch delete twitch streams for %v: %v", gameID, err)
			return err
		}

		if len(records) == 0 {
			return nil
		}

		if err := tx.Create(&records).Error; err != nil {
			log.Warnf("Failed to batch insert twitch streams for %v: %v", gameID, err)
			return err
		}

		return nil
	})
}

func StreamsWorker() {
	subgroup := sync.WaitGroup{}

	reservations := []models.Reservation{}
	db.Conn.Distinct("game_id", "name").Select("*").Find(&reservations)

	for _, reservation := range reservations {
		if reservation.Name == "" {
			name, err := twitch.GameName(reservation.GameID)
			if err == nil {
				reservation.Name = name
				log.Infof("Filled in missing name for %v (%v)", reservation.GameID, reservation.Name)
				db.Conn.Save(&reservation)
			}
		}

		log.Debugf("Fetching game ID %v (%v)", reservation.GameID, reservation.Name)
		subgroup.Add(1)
		go func(gid string) {
			defer subgroup.Done()

			fetch_streams(gid)
		}(reservation.GameID)
	}

	// Wait for all calls to come back so we only run these queries
	// at *max* as often as the ticker allows
	subgroup.Wait()

}
