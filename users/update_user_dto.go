package users

import "time"

type UpdateUserDTO struct {
	FirstName  string     `json:"firstName" gorm:"type:varchar(255);not null"`
	LastName   string     `json:"lastName" gorm:"type:varchar(255);not null"`
	FirebaseId string     `json:"firebaseId" gorm:"type:varchar(255)"`
	Premium    bool       `json:"premium" gorm:"bool"`
	Phone      string     `json:"phone" gorm:"type:varchar(255)"`
	BirthDate  *time.Time `json:"birthDate"`
}
