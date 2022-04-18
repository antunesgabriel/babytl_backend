package repositories

import (
	"github.com/antunesgabriel/babytl_backend/src/domain/user"
	"github.com/antunesgabriel/babytl_backend/src/infrastructure/models"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func (ur *UserRepositoryImpl) Create(user user.User) error {
	userModel := models.User{
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		Email:     user.Email(),
		Password:  user.Password(),
		Premium:   false,
	}

	err := ur.DB.Create(&userModel).Error

	return err
}

func (ur *UserRepositoryImpl) Update(user user.User) error {
	var userModel models.User

	if err := ur.DB.First(&userModel, "ID = ?", user.ID()).Error; err != nil {
		return err
	}

	userModel.ID = user.ID()
	userModel.FirstName = user.FirstName()
	userModel.LastName = user.LastName()
	userModel.Email = user.Email()
	userModel.Premium = user.IsPremium()
	userModel.BirthDate = user.BirthDate()
	userModel.Password = user.Password()

	err := ur.DB.Save(&userModel).Error

	return err
}

func (ur *UserRepositoryImpl) FindByEmail(email string) (*user.User, error) {
	var userModel models.User

	if err := ur.DB.First(&userModel, "email = ?", email).Error; err != nil {
		return new(user.User), err
	}

	user := ur.parseModelToEntity(userModel)

	return user, nil
}

func (ur *UserRepositoryImpl) FindById(id uint) (*user.User, error) {
	var userModel models.User

	if err := ur.DB.First(&userModel, "ID = ?", id).Error; err != nil {
		return new(user.User), err
	}

	user := ur.parseModelToEntity(userModel)

	return user, nil
}

func (ur *UserRepositoryImpl) parseModelToEntity(userModel models.User) *user.User {
	user := user.NewUser(
		userModel.FirstName,
		userModel.LastName,
		userModel.Email,
		userModel.Password,
	)

	user.AddWhatsApp(userModel.WhatsApp)

	if userModel.Premium {
		user.ActivePremium()
	} else {
		user.DisablePremium()
	}

	user.DefineBirthDate(userModel.BirthDate)

	return user
}
