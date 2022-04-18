package usecases

import (
	"github.com/antunesgabriel/babytl_backend/src/application/dtos"
	u "github.com/antunesgabriel/babytl_backend/src/domain/user"
)

type ChangeUserSignature struct {
	UserRepository u.Repository
}

func (uc *ChangeUserSignature) Execute(dto *dtos.ChangeUserSignatureDTO) error {
	user, err := uc.UserRepository.FindById(dto.UserId)

	if err != nil {
		return err
	}

	if dto.IsPremium {
		user.ActivePremium()
	} else {
		user.DisablePremium()
	}

	err = uc.UserRepository.Update(user)

	return err
}
