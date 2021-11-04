package api

import (
	"net/http"

	"github.com/antunesgabriel/babytl_backend/albums"
	"github.com/antunesgabriel/babytl_backend/api/middlewares"
	"github.com/antunesgabriel/babytl_backend/auth"
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

		authGroup := firstVersion.Group("auth")
		{
			authGroup.POST("local", auth.HandlerLoginWithEmail)
		}

		userGroup := firstVersion.Group("users")
		{
			userGroup.POST("", users.HandlerStore)
			userGroup.GET("me", middlewares.AuthMiddleware(), users.HandlerShow)
		}

		albumGroup := firstVersion.Group("albums", middlewares.AuthMiddleware())
		{
			albumGroup.GET("", albums.HandlerIndex)
			albumGroup.POST("", albums.HandlerStore)
			albumGroup.DELETE(":albumId", albums.HandlerDestroy)
			albumGroup.PUT(":albumId", albums.HandlerUpdate)
		}
	}

	return router
}
