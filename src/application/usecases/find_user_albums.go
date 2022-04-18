package usecases

import (
	"github.com/antunesgabriel/babytl_backend/src/application/dtos"
	alb "github.com/antunesgabriel/babytl_backend/src/domain/album"
)

type FindUserAlbums struct {
	AlbumRepository alb.Repository
}

func (uc *FindUserAlbums) Execute(dto *dtos.FindUserAlbumsDTO) (*[]dtos.AlbumDTO, error) {
	albums, err := uc.AlbumRepository.FindByUserId(dto.UserId)

	var albumsDTO []dtos.AlbumDTO
	{
	}

	for _, album := range albums {
		albumsDTO = append(albumsDTO, dtos.AlbumDTO{
			ID:       album.ID(),
			Title:    album.Title(),
			ThumbUrl: album.ThumbUrl(),
			Gender:   album.Gender(),
			UserId:   album.UserId(),
		})
	}

	return &albumsDTO, err
}
