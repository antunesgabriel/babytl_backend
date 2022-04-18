package repositories

import (
	"github.com/antunesgabriel/babytl_backend/src/infrastructure/models"
	"gorm.io/gorm"
)

type AlbumRepositoryImp struct {
	DB *gorm.DB
}

func (rep *AlbumRepositoryImp) FindByUserId() {

}

func (rep *AlbumRepositoryImp) parseModelToEntity(albumModel *models.Album) {

}
