package snaps

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

const FOLDER = "snaps"

func HandlerIndex(c *gin.Context) {
	db := database.GetDatabase()

	albumId, err1 := strconv.Atoi(c.Query("albumId"))
	month, err2 := time.Parse("2006-01-02", c.Query("month"))

	if err1 != nil || err2 != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "INVALID_PARAMS",
		})

		return
	}

	var snaps []entities.Snap

	now := time.Now()

	if db.Where("created_at BETWEEN ? AND ? AND album_id = ?", month, now, albumId).Find(&snaps).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"snaps": snaps,
	})
}

func HandlerStore(c *gin.Context) {
	db := database.GetDatabase()
	authId := c.GetUint("authId")
	albumId, err := strconv.Atoi(c.PostForm("albumId"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "INVALID_PARAMS",
		})

		return
	}

	snapFile, fileErr := c.FormFile("snap")

	if fileErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "INVALID_PARAMS",
		})

		return
	}

	oldFilename := filepath.Base(snapFile.Filename)
	ext := filepath.Ext(oldFilename)

	timeUnix := time.Now().Unix()

	newFileName := fmt.Sprint("user_id_", authId, "_album_id_", albumId, "_", timeUnix, ext)

	dir := filepath.Join("tmp", newFileName)

	if errUpload := c.SaveUploadedFile(snapFile, dir); errUpload != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":    "INTERNAL",
			"_details": errUpload.Error(),
		})

		return
	}

	var snap entities.Snap
	var album entities.Album

	if db.First(&album, albumId).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	snap.FileName = newFileName
	snap.AlbumID = album.ID
	snap.Album = album

	if db.Create(&snap).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":    "INTERNAL",
			"_details": "On insert",
		})

		return
	}

	go workerUpload(dir, snap.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "CREATED",
	})
}

func HandlerDestroy(c *gin.Context) {
	snapId, err := strconv.Atoi(c.Param("snapId"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "INVALID_PARAMS",
		})

		return
	}

	db := database.GetDatabase()

	var snap entities.Snap

	if db.Where("ID = ?", snapId).Delete(&snap).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":    "INTERNAL",
			"_details": "On delete",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "SUCCESS",
	})
}

func workerUpload(dir string, snapId uint) {
	fmt.Println("INIT ROUTINE")
	defer fmt.Println("FINISH ROUTINE")

	db := database.GetDatabase()
	var snap entities.Snap

	if db.First(&snap, snapId).Error != nil {
		log.Fatalln("error on acess snap")

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

	snap.SnapUrl = fileUrl

	db.Save(&snap)

	errRemove := os.Remove(dir)

	if errRemove != nil {
		fmt.Println("error on remove file to tmp", errRemove)
	}
}
