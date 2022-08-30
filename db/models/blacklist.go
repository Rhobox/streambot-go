package models

import (
	"gorm.io/gorm"
)

type BlacklistedUser struct {
	gorm.Model
	Username  string `gorm:"index"`
	ChannelID string `gorm:"index"`
}
