package entities

import "gorm.io/gorm"

type TimeLine struct {
	gorm.Model
	SnapUrl string `json:"snapUrl" gorm:"type:varchar(255);not null"`
	AlbumID uint
	Album   Album `json:"album" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
