package api

import (
	"net/http"

	"github.com/antunesgabriel/babytl_backend/users"
	"github.com/gin-gonic/gin"
)

func ConfigureRoutes(router *gin.Engine) *gin.Engine {
	firstVersion := router.Group("api/v1")
	{
		health := firstVersion.Group("health")
		{
			health.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"answerIs": 42,
				})
			})
		}

		usersGroup := firstVersion.Group("users")
		{
			usersGroup.POST("", users.HandlerStore)
		}

	}

	return router
}
