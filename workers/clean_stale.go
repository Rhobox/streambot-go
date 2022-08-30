package workers

import (
	"streambot/db"
	"streambot/db/models"

	"gorm.io/gorm"
)

func CleanStaleMessages() {
	live_streams := []models.TwitchStream{}

	db.Conn.Transaction(func(tx *gorm.DB) error {
		tx.Find(&live_streams)
		if len(live_streams) == 0 {
			return nil
		}

		log.Debug("Marking old stream messages for deletion")

		usernames := make([]string, len(live_streams))
		for idx, stream := range live_streams {
			usernames[idx] = stream.Username
		}

		// Clear out people that are no longer streaming
		if err := tx.Unscoped().Where("username NOT IN ?", usernames).Delete(&models.Stream{}).Error; err != nil {
			log.Warnf("Failed to clear out stale usernames: %v", err)
			return err
		}

		// Clear out people that are no longer streaming the same game
		for _, stream := range live_streams {
			if err := tx.Unscoped().Where("username = ?", stream.Username).Where("game_id != ?", stream.GameID).Delete(&models.Stream{}).Error; err != nil {
				log.Warnf("Failed to clear out streams with different games: %v", err)
				return err
			}
		}

		return nil
	})
}
