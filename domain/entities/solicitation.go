package entities

import "time"

type Solicitation struct {
	Base
	AttendedAt time.Timer `json:"attendedAt" gorm:"type:datetime"`
	FilesUrl   string     `json:"filesUrl" gorm:"type:varchar(255)"`
	UserID     string
	User       User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
