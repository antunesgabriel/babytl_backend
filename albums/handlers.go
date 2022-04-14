package albums

import (
	"fmt"
	models2 "github.com/antunesgabriel/babytl_backend/src/infrastructure/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/antunesgabriel/babytl_backend/database"
	"github.com/antunesgabriel/babytl_backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const FOLDER = "thumbs"

func HandlerStore(c *gin.Context) {
	db := database.GetDatabase()
	var album models2.Album
	var user models2.User

	title := c.PostForm("title")
	gender := c.PostForm("gender")
	authId := c.GetUint("authId")

	thumbFile, err := c.FormFile("thumb")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "INVALID_PARAMS",
		})

		return
	}

	oldFilename := filepath.Base(thumbFile.Filename)
	ext := filepath.Ext(oldFilename)

	timeUnix := time.Now().Unix()

	newFileName := fmt.Sprint("user_id_", authId, "_album_", title, "_", timeUnix, ext)

	if err := os.MkdirAll("tmp", os.ModePerm); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":    "INTERNAL",
			"_details": err.Error(),
		})
	}

	dir := filepath.Join("tmp", newFileName)

	if errUpload := c.SaveUploadedFile(thumbFile, dir); errUpload != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":    "INTERNAL",
			"_details": errUpload.Error(),
		})

		return
	}

	if db.Where("ID = ?", authId).First(&user).Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "USER_NOT_FOUND",
		})

		return
	}

	album.Title = title
	album.Gender = gender
	album.UserID = user.ID
	album.User = user

	thumbUrl, errUpload := workerUpload(dir, album.ID, db)

	if errUpload != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	album.ThumbUrl = thumbUrl

	if db.Create(&album).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Falha ao criar album",
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

	var albums []models2.Album

	if db.Where("user_id = ?", authId).Find(&albums).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	for idx, album := range albums {
		albums[idx].SnapsCount = db.Model(&album).Association("Snaps").Count()
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

	var album models2.Album

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
		"album":   album,
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

	var album models2.Album

	if db.Where("ID = ?", albumId).Delete(&album).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "DELETED",
		"album":   album,
	})
}

func HandlerShow(c *gin.Context) {
	albumId, err := strconv.Atoi(c.Param("albumId"))

	db := database.GetDatabase()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "INVALID_PARAMS",
		})

		return
	}

	var album models2.Album

	findError := db.Preload("Snaps").First(&album, albumId).Error

	if findError != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "NOT_FOUND",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"album": album,
	})
}

func workerUpload(dir string, albumId uint, db *gorm.DB) (string, error) {
	s3Handler, err := utils.NewS3Handler()

	if err != nil {
		utils.RegisterLog("[UPLOAD ALBUM THUMB TO S3] Line 201:", err.Error())

		return "", err
	}

	fileUrl, errUpload := s3Handler.UploadFile(dir, FOLDER)

	if errUpload != nil {
		utils.RegisterLog("[UPLOAD ALBUM THUMB TO S3] Line 209:", errUpload.Error())

		return "", errUpload
	}

	errRemove := os.Remove(dir)

	if errRemove != nil {
		utils.RegisterLog("[UPLOAD ALBUM THUMB TO S3] Line 217:", errRemove.Error())

		return "", errRemove
	}

	return fileUrl, nil
}
