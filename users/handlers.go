package users

import (
	"errors"
	"github.com/antunesgabriel/babytl_backend/src/infrastructure/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/antunesgabriel/babytl_backend/database"
)

func HandlerStore(c *gin.Context) {
	db := database.GetDatabase()

	var user models.User

	err := c.ShouldBindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "INVALID_PARAMS",
		})

		return
	}

	var exist models.User

	result := db.First(&exist, "email = ?", user.Email)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusExpectationFailed, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	if result.RowsAffected > 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "USER_EXISTS",
		})

		return
	}

	err = db.Create(&user).Error

	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	c.JSON(http.StatusCreated, user)
}

func HandlerShow(c *gin.Context) {
	db := database.GetDatabase()
	userId := c.GetUint("authId")

	var user models.User

	if db.Where("ID = ?", userId).First(&user).Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "USER_NOT_FOUND",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func HandlerUpdate(c *gin.Context) {
	authId := c.GetUint("authId")

	db := database.GetDatabase()

	var user models.User
	var updateUserDTO UpdateUserDTO

	if c.ShouldBindJSON(&updateUserDTO) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "INVALID_PARAMS",
		})

		return
	}

	if db.First(&user, authId).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	user.FirebaseId = updateUserDTO.FirebaseId
	user.FirstName = updateUserDTO.FirstName
	user.LastName = updateUserDTO.LastName
	user.Premium = updateUserDTO.Premium
	user.Phone = updateUserDTO.Phone

	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "SUCCESS",
	})
}
