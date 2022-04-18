package user

import (
	"errors"
	"time"
)

type UseCase struct {
	Repository Repository
}

func (uc *UseCase) CreateUser(firstName, lastName, email, password string) error {
	user, err := NewUser(firstName, lastName, email, password)

	if err != nil {
		return err
	}

	exist, err := uc.Repository.FindByEmail(user.Email())

	if err != nil {
		return err
	}

	if exist.Email() != "" {
		return errors.New("email has been in using")
	}

	err = uc.Repository.Create(user)

	return err
}

func (uc *UseCase) AddMoreUserInfo(id uint, phone string, birthDate *time.Time) error {
	user, err := uc.Repository.FindById(id)

	if err != nil {
		return err
	}

	user.AddWhatsApp(phone)
	user.DefineBirthDate(birthDate)

	err = uc.Repository.Update(user)

	return err
}

func (uc *UseCase) ChangePremiumStatus(id uint, isPremium bool) error {
	user, err := uc.Repository.FindById(id)

	if err != nil {
		return err
	}

	if isPremium {
		user.ActivePremium()
	} else {
		user.DisablePremium()
	}

	err = uc.Repository.Update(user)

	return err
}
