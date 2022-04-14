package snaps

import (
	"errors"
	"fmt"
	models2 "github.com/antunesgabriel/babytl_backend/src/infrastructure/models"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/antunesgabriel/babytl_backend/database"
	"github.com/antunesgabriel/babytl_backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const FOLDER = "snaps"

func HandlerIndex(c *gin.Context) {
	db := database.GetDatabase()

	albumId, err1 := strconv.Atoi(c.Param("albumId"))

	if err1 != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "INVALID_PARAMS",
		})

		return
	}

	var snaps []models2.Snap

	if db.Where("album_id = ?", albumId).Find(&snaps).Error != nil {
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

	timeNow := time.Now()
	year, month, day := timeNow.Date()
	startDay := time.Date(year, month, day, 0, 0, 0, 0, timeNow.Location())
	endDay := time.Date(year, month, day, 23, 59, 59, 0, timeNow.Location())

	var snapExist models2.Snap

	result := db.First(&snapExist, "album_id = ? AND created_at BETWEEN ? AND ?", albumId, startDay, endDay)
	errFind := result.Error

	if errFind != nil && !errors.Is(errFind, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":    "INTERNAL",
			"_details": errFind.Error(),
		})

		return
	}

	if result.RowsAffected > 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Já foi feito um registro para o dia de hoje.",
			"snap":    snapExist,
			"error":   "NOT_PERMITED_SNAP",
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

	var album models2.Album
	var snap models2.Snap

	if db.First(&album, albumId).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	snap.FileName = newFileName
	snap.AlbumID = album.ID
	snap.Album = album

	snapUrl, errUpload := workerUpload(dir, snap.ID, snap.AlbumID)

	if errUpload != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	snap.SnapUrl = snapUrl

	if db.Create(&snap).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":    "INTERNAL",
			"_details": "On insert",
		})

		return
	}

	fmt.Println("=== Upload Feito")

	c.JSON(http.StatusCreated, gin.H{
		"message": "CREATED",
		"snap":    snap,
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

	var snap models2.Snap

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

func workerUpload(dir string, snapId uint, albumId uint) (string, error) {
	fmt.Println("INIT ROUTINE")
	defer fmt.Println("FINISH ROUTINE")

	s3Handler, err := utils.NewS3Handler()

	if err != nil {
		utils.RegisterLog("[UPLOAD SNAP - LINE: 183]", err.Error())

		return "", err
	}

	folder := path.Join(FOLDER, fmt.Sprint("ALBUM_ID_", albumId))

	fileUrl, errUpload := s3Handler.UploadFile(dir, folder)

	if errUpload != nil {
		utils.RegisterLog("[UPLOAD SNAP - LINE: 192]", errUpload.Error())

		return "", errUpload
	}

	errRemove := os.Remove(dir)

	if errRemove != nil {
		utils.RegisterLog("[UPLOAD SNAP - LINE: 200]", err.Error())
	}

	return fileUrl, nil
}

func HandlerShow(c *gin.Context) {
	albumId, err := strconv.Atoi(c.Param("albumId"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "INVALID_PARAMS",
		})

		return
	}

	db := database.GetDatabase()

	var snap models2.Snap

	timeNow := time.Now()
	year, month, day := timeNow.Date()

	startDay := time.Date(year, month, day, 0, 0, 0, 0, timeNow.Location())
	endDay := time.Date(year, month, day, 23, 58, 58, 0, timeNow.Location())

	errFirst := db.First(&snap, "album_id = ? AND created_at BETWEEN ? AND ?", albumId, startDay, endDay).Error

	if errFirst != nil {

		if !errors.Is(errFirst, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "INTERNAL",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"snap": nil,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"snap": snap,
	})
}
