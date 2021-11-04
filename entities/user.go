package entities

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Base
	FirstName string     `json:"firstName" gorm:"type:varchar(255);not null"`
	LastName  string     `json:"lastName" gorm:"type:varchar(255);not null"`
	Email     string     `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string     `json:"-" gorm:"type:varchar(255);not null"`
	Premium   bool       `json:"premium" gorm:"bool"`
	Phone     string     `json:"phone" gorm:"type:varchar(255)"`
	BirthDate time.Timer `json:"birthDate" gorm:"type:datetime"`
}

func (user *User) Prepare() error {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(passwordBytes)

	return nil
}
