package auth

import (
	"github.com/antunesgabriel/babytl_backend/src/infrastructure/models"
	"net/http"

	"github.com/antunesgabriel/babytl_backend/configs"
	"github.com/antunesgabriel/babytl_backend/database"
	"github.com/gin-gonic/gin"
)

func HandlerLoginWithEmail(c *gin.Context) {
	var authDTO AuthWithEmailDTO

	if c.ShouldBindJSON(&authDTO) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "INVALID_CREDENTIALS",
		})

		return
	}

	var user models.User
	db := database.GetDatabase()

	if err := db.Where("email = ?", authDTO.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "USER_NOT_EXIST",
		})

		return
	}

	if user.CheckPass(authDTO.Password) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "INVALID_PASSWORD",
		})

		return
	}

	jwtConfig := configs.NewJWT()

	token, err := jwtConfig.GenerateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"type":  "Bearer",
	})

}
