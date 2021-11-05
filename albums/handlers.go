package albums

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/antunesgabriel/babytl_backend/database"
	"github.com/antunesgabriel/babytl_backend/entities"
	"github.com/antunesgabriel/babytl_backend/utils"
	"github.com/gin-gonic/gin"
)

const FOLDER = "thumbs"

func HandlerStore(c *gin.Context) {
	db := database.GetDatabase()
	var album entities.Album
	var user entities.User

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

	if db.Create(&album).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Falha ao criar album",
		})

		return
	}

	go workerUpload(dir, album.ID)

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

	var album entities.Album

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

func workerUpload(dir string, albumId uint) {
	fmt.Println("INIT ROUTINE")
	defer fmt.Println("FINISH ROUTINE")

	db := database.GetDatabase()
	var album entities.Album

	if db.First(&album, albumId).Error != nil {
		log.Fatalln("error on acess album")

		return
	}

	s3Handler, err := utils.NewS3Handler()

	if err != nil {
		log.Fatalf("[error] on connect s3: %s", err.Error())

		return
	}

	fileUrl, errUpload := s3Handler.UploadFile(dir, FOLDER)

	if errUpload != nil {
		log.Fatalf("error on upload s3: %v", errUpload)

		return
	}

	album.ThumbUrl = fileUrl

	db.Save(&album)

	errRemove := os.Remove(dir)

	if errRemove != nil {
		fmt.Println("error on remove file to tmp", errRemove)
	}
}
