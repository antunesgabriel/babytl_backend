package user

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	id         uint
	firstName  string
	lastName   string
	email      string
	firebaseId string
	password   string
	premium    bool
	whatsApp   string
	birthDate  *time.Time
}

func NewUser(firstName, lastName, email, password string) (*User, error) {
	user := new(User)

	user.firstName = strings.TrimSpace(firstName)
	user.lastName = strings.TrimSpace(lastName)
	user.email = strings.TrimSpace(email)
	user.password = strings.TrimSpace(password)

	err := user.validation()

	return user, err
}

func (u *User) ID() uint {
	return u.id
}

func (u *User) FirstName() string {
	return u.firstName
}

func (u *User) LastName() string {
	return u.lastName
}

func (u *User) Email() string {
	return u.email
}

func (u *User) FirebaseId() string {
	return u.firebaseId
}

func (u *User) Password() string {
	return u.password
}

func (u *User) IsPremium() bool {
	return u.premium
}

func (u *User) WhatsApp() string {
	return u.whatsApp
}

func (u *User) BirthDate() *time.Time {
	return u.birthDate
}

func (u *User) validation() error {
	if u.email == "" {
		return errors.New("email is required")
	}

	if u.firstName == "" {
		return errors.New("email is required")
	}

	if u.lastName == "" {
		return errors.New("email is required")
	}

	if u.password == "" {
		return errors.New("email is required")
	}

	return nil
}

func (u *User) AddWhatsApp(whatsApp string) {
	u.whatsApp = whatsApp
}

func (u *User) DefineBirthDate(birthDate *time.Time) {
	u.birthDate = birthDate
}

func (u *User) ActivePremium() {
	u.premium = true
}

func (u *User) DisablePremium() {
	u.premium = false
}
