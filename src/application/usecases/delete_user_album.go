package usecases

import (
	"errors"
	"github.com/antunesgabriel/babytl_backend/src/application/dtos"
	alb "github.com/antunesgabriel/babytl_backend/src/domain/album"
)

type DeleteUserAlbum struct {
	AlbumRepository alb.Repository
}

func (uc *DeleteUserAlbum) Execute(dto *dtos.DeleteUserAlbum) error {
	album, err := uc.AlbumRepository.FindByUserIdAndAlbumId(dto.UserId, dto.AlbumId)

	if album == nil {
		return errors.New("album is not exist for this user")
	}

	err = uc.AlbumRepository.Destroy(dto.AlbumId)

	return err
}
