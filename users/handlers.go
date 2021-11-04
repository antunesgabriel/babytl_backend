package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/antunesgabriel/babytl_backend/database"
	"github.com/antunesgabriel/babytl_backend/entities"
)

func HandlerStore(c *gin.Context) {
	db := database.GetDatabase()

	var user entities.User

	err := c.ShouldBindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot bind JSON.",
		})

		return
	}

	err = db.Create(&user).Error

	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, user)
}
