package user

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	Id         uint
	FirstName  string
	LastName   string
	Email      string
	FirebaseId string
	Password   string
	Premium    bool
	Phone      string
	BirthDate  *time.Time
}

func NewUser(firstName, lastName, email, password string) (*User, error) {
	user := new(User)

	user.FirstName = strings.TrimSpace(firstName)
	user.LastName = strings.TrimSpace(lastName)
	user.Email = strings.TrimSpace(email)
	user.Password = strings.TrimSpace(password)

	err := user.validation()

	return user, err
}

func (u *User) validation() error {
	if u.Email == "" {
		return errors.New("email is required")
	}

	if u.FirstName == "" {
		return errors.New("email is required")
	}

	if u.LastName == "" {
		return errors.New("email is required")
	}

	if u.Password == "" {
		return errors.New("email is required")
	}

	return nil
}

func (u *User) AddContact(phone string) {
	u.Phone = phone
}

func (u *User) DefineBirthDate(birthDate *time.Time) {
	u.BirthDate = birthDate
}

func (u *User) ActivePremium() {
	u.Premium = true
}

func (u *User) DisablePremium() {
	u.Premium = true
}
