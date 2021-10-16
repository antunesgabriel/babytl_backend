package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        string         `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Timer     `json:"createdAt" gorm:"type:datetime"`
	UpdatedAt time.Timer     `json:"updatedAt" gorm:"type:datetime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"type:datetime" sql:"index"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) {
	tx.Statement.SetColumn("CreatedAt", time.Now())
	tx.Statement.SetColumn("ID", uuid.NewV4().String())
}
