package solicitations

import (
	"net/http"
	"time"

	"github.com/antunesgabriel/babytl_backend/database"
	"github.com/antunesgabriel/babytl_backend/entities"
	"github.com/gin-gonic/gin"
)

const (
	PREMIUM_WAIT_HOURS = time.Hour * 24 * 30
	NORMAL_WAIT_HOURS = time.Hour * 24 * 30 * 3
)

func HandlerStore (c *gin.Context) {
	db := database.GetDatabase()
	authId := c.GetUint("authId")

	var user entities.User
	var solicitations []entities.Solicitation

	if db.Where("ID = ?", authId).First(&user).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	if db.Model(&user).Association("Solicitations").Find(&solicitations) != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	if len(solicitations) != 0 {
		var lastSolicitationDate time.Time

		for idx, item := range solicitations {
			if idx == 0 {
				lastSolicitationDate = item.CreatedAt

				continue
			}

			if lastSolicitationDate.Before(item.CreatedAt) {
				lastSolicitationDate = item.CreatedAt
			}
		}

		if user.Premium {
			nextDateToSolicitation := lastSolicitationDate.Add(PREMIUM_WAIT_HOURS)

			message, isValid := validateToNewSolicictation(nextDateToSolicitation)

			if !isValid {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": message,
				})

				return
			}
			
		} else {
			nextDateToSolicitation := lastSolicitationDate.Add(NORMAL_WAIT_HOURS)

			message, isValid := validateToNewSolicictation(nextDateToSolicitation)

			if !isValid {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": message,
				})

				return
			}
		}
	}

	var newSolicitation  entities.Solicitation

	if db.Model(&user).Association("Solicitations").Append(&newSolicitation) != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
			"_datails": "On append new solicitation",
		})

		return
	}

	//TODO: Implement routine to download snaps and zip and update solicitation

	c.JSON(http.StatusCreated, gin.H{
		"message": "SUCCESS",
	})
}

func validateToNewSolicictation (nextDateToSolicitation time.Time) (string, bool) {
	now := time.Now()

	if !now.After(nextDateToSolicitation) {
		diff := nextDateToSolicitation.Sub(now)

		message := buildMessageDiff(diff)

		return message, false
	}

	return "", true
}

func buildMessageDiff(diff time.Duration) (message string) {
	date := time.Now().Add(diff)

	message = "Sua próxima data de solicitação será no dia: " + date.Format("02/01/2006") + " ás: " + date.Format("15:04") + "hrs"

	return 
}