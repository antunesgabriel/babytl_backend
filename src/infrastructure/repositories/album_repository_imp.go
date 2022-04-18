package repositories

import "gorm.io/gorm"

type AlbumRepositoryImp struct {
	DB *gorm.DB
}
