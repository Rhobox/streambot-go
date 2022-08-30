package models

import (
	"gorm.io/gorm"
)

type TwitchStream struct {
	gorm.Model
	Username     string `gorm:"index"`
	DisplayName  string
	GameID       string `gorm:"index"`
	Description  string
	ThumbnailURL string
	Speedrun     bool
}
