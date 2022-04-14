package solicitations

import (
	"fmt"
	models2 "github.com/antunesgabriel/babytl_backend/src/infrastructure/models"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/antunesgabriel/babytl_backend/database"
	"github.com/antunesgabriel/babytl_backend/utils"
	"github.com/gin-gonic/gin"
)

const (
	PREMIUM_WAIT_HOURS   = time.Hour * 24 * 30
	NORMAL_WAIT_HOURS    = time.Hour * 24 * 30 * 3
	FOLDER               = "snaps"
	FOLDER_SOLICITATIONS = "solicitations"
)

func HandlerStore(c *gin.Context) {
	db := database.GetDatabase()
	authId := c.GetUint("authId")

	var solicitationDTO CreateSolicitationDTO

	if c.ShouldBindJSON(&solicitationDTO) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "INVALID_PARAMS",
		})

		return
	}

	var user models2.User
	var solicitations []models2.Solicitation
	var album models2.Album

	if db.Where("ID = ?", authId).First(&user).Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "INTERNAL",
		})

		return
	}

	if db.First(&album, solicitationDTO.AlbumID).Error != nil {
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

		if user.Premium {
			forThisAlbum := false
			var lastSolicitationDate time.Time

			for idx, item := range solicitations {
				if item.AlbumID != solicitationDTO.AlbumID {
					continue
				}

				forThisAlbum = true

				if idx == 0 {
					lastSolicitationDate = item.CreatedAt

					continue
				}

				if lastSolicitationDate.Before(item.CreatedAt) {
					lastSolicitationDate = item.CreatedAt
				}
			}

			if forThisAlbum {
				nextDateToSolicitation := lastSolicitationDate.Add(PREMIUM_WAIT_HOURS)

				message, isValid := validateToNewSolicictation(nextDateToSolicitation, user.Premium)

				if !isValid {
					c.JSON(http.StatusUnauthorized, gin.H{
						"message": message,
					})

					return
				}
			}

		} else {
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

			nextDateToSolicitation := lastSolicitationDate.Add(NORMAL_WAIT_HOURS)

			message, isValid := validateToNewSolicictation(nextDateToSolicitation, user.Premium)

			if !isValid {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": message,
				})

				return
			}
		}
	}

	var newSolicitation models2.Solicitation

	newSolicitation.AlbumID = album.ID
	newSolicitation.Album = album

	if db.Model(&user).Association("Solicitations").Append(&newSolicitation) != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":    "INTERNAL",
			"_datails": "On append new solicitation",
		})

		return
	}

	go workerSolicitation(solicitationDTO.AlbumID, newSolicitation.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "SUCCESS",
	})
}

func validateToNewSolicictation(nextDateToSolicitation time.Time, isPremium bool) (string, bool) {
	now := time.Now()

	if !now.After(nextDateToSolicitation) {
		diff := nextDateToSolicitation.Sub(now)

		message := buildMessageDiff(diff, isPremium)

		return message, false
	}

	return "", true
}

func buildMessageDiff(diff time.Duration, isPremium bool) (message string) {
	date := time.Now().Add(diff)

	if isPremium {
		message = "Próxima solicitação: " + date.Format("02/01/2006") + " ás: " + date.Format("15:04") + "hrs"

		return
	}

	message = "Próxima solicitação: " + date.Format("02/01/2006") + " ás: " + date.Format("15:04") + "hrs"

	return
}

func workerSolicitation(albumId, solicitationId uint) {
	fmt.Println("MAKE NEW SOLICITATION")

	defer fmt.Println("CLOSE SOLICITATION")

	db := database.GetDatabase()

	var snaps []models2.Snap

	if err := db.Where("album_id = ?", albumId).Find(&snaps).Error; err != nil {
		indentify := fmt.Sprintf("WORKER_SOLICITATION - ON SELECT SNAPS ALBUM_ID: %s", fmt.Sprint(albumId))

		utils.RegisterLog(indentify, err.Error())

		return
	}

	folderName := path.Join(FOLDER, fmt.Sprint("ALBUM_ID_", albumId))
	oldPutFolder := path.Join("tmp", folderName)

	if err := os.MkdirAll(oldPutFolder, os.ModePerm); err != nil {
		indentify := fmt.Sprintf("WORKER_SOLICITATION - ON CREATE FOLDER TO ALBUM_ID: %s", fmt.Sprint(albumId))

		utils.RegisterLog(indentify, err.Error())

		return
	}

	s3Hanlder, err := utils.NewS3Handler()

	if err != nil {
		indentify := fmt.Sprintf("WORKER_SOLICITATION - ON OPEN SESSION TO DOWNLOAD ALBUM_ID: %s", fmt.Sprint(albumId))

		utils.RegisterLog(indentify, err.Error())

		return
	}

	files := make([]string, len(snaps))

	for idx, snap := range snaps {
		files[idx] = snap.FileName
	}

	if err := s3Hanlder.DownloadFiles(files, oldPutFolder, folderName); err != nil {
		indentify := fmt.Sprintf("WORKER_SOLICITATION - ON DOWNLOAD ALBUM_ID: %s", fmt.Sprint(albumId))

		utils.RegisterLog(indentify, err.Error())

		return
	}

	zipName := oldPutFolder + ".zip"

	if err := utils.ZipSource(oldPutFolder, zipName); err != nil {
		indentify := fmt.Sprintf("WORKER_SOLICITATION - ON ZIP ALBUM_ID: %s", fmt.Sprint(albumId))

		utils.RegisterLog(indentify, err.Error())

		return
	}

	var solicitation models2.Solicitation

	if err := db.First(&solicitation, solicitationId).Error; err != nil {
		indentify := fmt.Sprintf("WORKER_SOLICITATION - ON CLOSE ZIP ALBUM_ID: %s", fmt.Sprint(albumId))

		utils.RegisterLog(indentify, err.Error())

		return
	}

	fileUrl, err := s3Hanlder.UploadFile(zipName, FOLDER_SOLICITATIONS)

	if err != nil {
		indentify := fmt.Sprintf("WORKER_SOLICITATION - ON UPLOAD ZIP ALBUM_ID: %s", fmt.Sprint(albumId))

		utils.RegisterLog(indentify, err.Error())

		return
	}

	solicitation.ZipUrl = fileUrl

	db.Save(solicitation)

	if err := os.RemoveAll(oldPutFolder); err != nil {
		indentify := fmt.Sprintf("WORKER_SOLICITATION - ON REMOVE OUTPUT ALBUM_ID: %s", fmt.Sprint(albumId))

		utils.RegisterLog(indentify, err.Error())

		return
	}

	if err := os.Remove(zipName); err != nil {
		indentify := fmt.Sprintf("WORKER_SOLICITATION - ON REMOVE ZIP ALBUM_ID: %s", fmt.Sprint(albumId))

		utils.RegisterLog(indentify, err.Error())

		return
	}
}
