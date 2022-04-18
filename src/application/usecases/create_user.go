package usecases

import (
	"errors"
	"github.com/antunesgabriel/babytl_backend/src/application/dtos"
	u "github.com/antunesgabriel/babytl_backend/src/domain/user"
)

type CreateAUser struct {
	UserRepository u.Repository
}

func (uc *CreateAUser) Execute(dto *dtos.CreateUserDTO) error {
	user, err := u.NewUser(
		dto.FirstName,
		dto.LastName,
		dto.Email,
		dto.Email,
		0)

	if err != nil {
		return err
	}

	exist, err := uc.UserRepository.FindByEmail(user.Email())

	if err != nil {
		return err
	}

	if exist.Email() != "" {
		return errors.New("email has been in using")
	}

	err = uc.UserRepository.Create(user)

	return err
}
