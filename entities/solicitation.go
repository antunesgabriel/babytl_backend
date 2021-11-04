package entities

import (
	"time"

	"gorm.io/gorm"
)

type Solicitation struct {
	gorm.Model
	AttendedAt time.Time `json:"attendedAt"`
	ZipUrl     string    `json:"zipUrl" gorm:"type:varchar(255)"`
	UserID     uint
	User       User `json:"user" gorm:"foreignKey:UserID"`
}
