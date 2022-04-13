package entities

import "gorm.io/gorm"

type Snap struct {
	gorm.Model
	FileName string `json:"fileName" gorm:"type:varchar(255);not null"`
	SnapUrl  string `json:"snapUrl" gorm:"type:varchar(255)"`
	AlbumID  uint   `json:"albumId"`
	Album    Album  `json:"album" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
