package usecases

import (
	"github.com/antunesgabriel/babytl_backend/src/application/dtos"
	u "github.com/antunesgabriel/babytl_backend/src/domain/user"
)

type ChangeUserDetails struct {
	UserRepository u.Repository
}

func (uc *ChangeUserDetails) Execute(dto *dtos.ChangeUserDetailsDTO) error {
	user, err := uc.UserRepository.FindById(dto.UserId)

	if err != nil {
		return err
	}

	user.AddWhatsApp(dto.WhatsApp)
	user.DefineBirthDate(dto.BirthDate)

	err = uc.UserRepository.Update(user)

	return err
}
