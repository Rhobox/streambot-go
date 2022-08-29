package workers

import (
	"streambot/db"
	"streambot/db/models"

	"gorm.io/gorm"
)

func CleanStaleMessages() {
	log.Debug("Marking old stream messages for deletion")
	live_streams := []models.TwitchStream{}

	db.Conn.Transaction(func(tx *gorm.DB) error {
		tx.Find(&live_streams)

		usernames := make([]string, len(live_streams))
		for idx, stream := range live_streams {
			usernames[idx] = stream.Username
		}

		// Clear out people that are no longer streaming
		if err := tx.Unscoped().Where("username NOT IN ?", usernames).Delete(&models.Stream{}).Error; err != nil {
			return err
		}

		// Clear out people that are no longer streaming the same game
		for _, stream := range live_streams {
			if err := tx.Unscoped().Where("username = ? AND game_id != ?", stream.Username, stream.GameID).Delete(&models.Stream{}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
