package entities

import "gorm.io/gorm"

type Album struct {
	gorm.Model
	Title    string `json:"title" gorm:"type:varchar(255);not null"`
	ThumbUrl string `json:"thumbUrl" gorm:"type:varchar(255)"`
	Gender   string `json:"gender" gorm:"type:varchar(1)"`
	UserID   uint
	User     User `json:"-" gorm:"foreignKey:UserID"`
}
