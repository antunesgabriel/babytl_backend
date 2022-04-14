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
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Premium:   false,
	}

	err := ur.DB.Create(&userModel).Error

	return err
}

func (ur *UserRepositoryImpl) Update(user user.User) error {
	var userModel models.User

	if err := ur.DB.First(&userModel, "ID = ?", user.Id).Error; err != nil {
		return err
	}

	userModel.ID = user.Id
	userModel.FirstName = user.FirstName
	userModel.LastName = user.LastName
	userModel.Email = user.Email
	userModel.Premium = user.Premium
	userModel.BirthDate = user.BirthDate
	userModel.Password = user.Password

	err := ur.DB.Save(&userModel).Error

	return err
}

func (ur *UserRepositoryImpl) FindByEmail(email string) (*user.User, error) {
	var userModel models.User

	if err := ur.DB.First(&userModel, "email = ?", email).Error; err != nil {
		return new(user.User), err
	}

	user := user.User{
		Id:        userModel.ID,
		FirstName: userModel.FirstName,
		LastName:  userModel.LastName,
		Email:     userModel.Email,
		Password:  userModel.Password,
		Premium:   userModel.Premium,
		BirthDate: userModel.BirthDate,
		Phone:     userModel.Phone,
	}

	return &user, nil
}

func (ur *UserRepositoryImpl) FindById(id uint) (*user.User, error) {
	var userModel models.User

	if err := ur.DB.First(&userModel, "ID = ?", id).Error; err != nil {
		return new(user.User), err
	}

	user := user.User{
		Id:        userModel.ID,
		FirstName: userModel.FirstName,
		LastName:  userModel.LastName,
		Email:     userModel.Email,
		Password:  userModel.Password,
		Premium:   userModel.Premium,
		BirthDate: userModel.BirthDate,
		Phone:     userModel.Phone,
	}

	return &user, nil
}
