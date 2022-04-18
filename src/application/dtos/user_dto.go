package dtos

import "time"

type CreateUserDTO struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type ChangeUserDetailsDTO struct {
	UserId    uint       `json:"userId"`
	BirthDate *time.Time `json:"birthDate"`
	WhatsApp  string     `json:"whatsApp"`
}

type ChangeUserSignatureDTO struct {
	UserId    uint `json:"userId"`
	IsPremium bool `json:"premium"`
}
