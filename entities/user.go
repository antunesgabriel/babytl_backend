package entities

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string    `json:"firstName" gorm:"type:varchar(255);not null"`
	LastName  string    `json:"lastName" gorm:"type:varchar(255);not null"`
	Email     string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string    `json:"password" gorm:"type:varchar(255);not null"`
	Premium   bool      `json:"premium" gorm:"bool"`
	Phone     string    `json:"phone" gorm:"type:varchar(255)"`
	BirthDate time.Time `json:"birthDate"`
	Albums    []Album   `json:"albums"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("[ERROR]:", err)
		return err
	}

	u.Password = string(passwordBytes)

	return
}

func (u *User) CheckPass(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
