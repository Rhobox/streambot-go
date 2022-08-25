package models

import (
	"gorm.io/gorm"
)

type Stream struct {
	gorm.Model
	ReservationID uint
	Username      string
	MessageID     int `gorm:"index"`
}
