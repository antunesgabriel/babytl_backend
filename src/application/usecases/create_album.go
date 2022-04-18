package usecases

import (
	"github.com/antunesgabriel/babytl_backend/src/application/dtos"
	alb "github.com/antunesgabriel/babytl_backend/src/domain/album"
)

type CreateAlbum struct {
	AlbumRepository alb.Repository
}

func (uc *CreateAlbum) Execute(dto *dtos.CreateAlbumDTO) error {

	album, err := alb.NewAlbum(dto.UserId, dto.Title, dto.ThumbUrl, dto.Gender, 0)

	if err != nil {
		return err
	}

	err = uc.AlbumRepository.Create(album)

	return err
}
