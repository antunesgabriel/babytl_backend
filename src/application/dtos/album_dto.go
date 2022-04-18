package dtos

type FindUserAlbumsDTO struct {
	UserId uint `json:"userId"`
}

type AlbumDTO struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	ThumbUrl string `json:"thumbUrl"`
	Gender   string `json:"gender"`
	UserId   uint   `json:"userId"`
}

type CreateAlbumDTO struct {
	UserId   uint   `json:"userId"`
	Title    string `json:"title"`
	Gender   string `json:"gender"`
	ThumbUrl string `json:"thumbUrl"`
}

type DeleteUserAlbum struct {
	UserId  uint `json:"userId"`
	AlbumId uint `json:"albumId"`
}
