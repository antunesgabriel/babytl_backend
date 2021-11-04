package entities

import "gorm.io/gorm"

type Album struct {
	gorm.Model
	Title    string `json:"title" gorm:"type:varchar(255);not null"`
	ThumbUrl string `json:"thumbUrl" gorm:"type:varchar(255)"`
	UserID   uint
	User     User `json:"user" gorm:"foreignKey:UserID"`
}
