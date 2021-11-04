package albums

import (
	"net/http"
	"strconv"

	"github.com/antunesgabriel/babytl_backend/database"
	"github.com/antunesgabriel/babytl_backend/entities"
	"github.com/gin-gonic/gin"
)

func HandlerStore(c *gin.Context) {
	db := database.GetDatabase()
	var album entities.Album
	var user entities.User

	if c.ShouldBindJSON(&album) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "INVALID_FIELDS",
		})

		return
	}

	authId := c.GetUint("authId")

	if db.Where("ID = ?", authId).First(&user).Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "USER_NOT_FOUND",
		})

		return
	}

	album.UserID = user.ID
	album.User = user

	if db.Create(&album).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Falha ao criar album",
		})

		return
	}

	if db.Save(&album).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Falha ao salvar album",
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "CREATED",
	})
}

func HandlerIndex(c *gin.Context) {
	db := database.GetDatabase()
	authId := c.GetUint("authId")

	var albums []entities.Album

	if db.Where("user_id = ?", authId).Find(&albums).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Falha ao buscar albums",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"albums": albums,
	})
}

func HandlerUpdate(c *gin.Context) {
	db := database.GetDatabase()
	albumId, err := strconv.Atoi(c.Param("albumId"))

	var updateAlbumDTO UpdateAlbumDTO

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "INVALID_PARAMS",
		})

		return
	}

	if c.ShouldBindJSON(&updateAlbumDTO) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "INVALID_PARAMS",
		})

		return
	}

	var album entities.Album

	if db.Where("ID = ?", albumId).First(&album).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	album.Title = updateAlbumDTO.Title
	album.Gender = updateAlbumDTO.Gender

	if db.Save(&album).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "UPDATED",
		"album": album,
	})
}

func HandlerDestroy(c *gin.Context) {
	db := database.GetDatabase()
	albumId, err := strconv.Atoi(c.Param("albumId"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "INVALID_PARAMS",
		})

		return
	}

	var album entities.Album

	if db.Where("ID = ?", albumId).Delete(&album).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "DELETED",
		"album": album,
	})
}