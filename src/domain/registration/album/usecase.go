package album

type UseCase struct {
	Repository Repository
}

func (uc *UseCase) ListUserAlbums(userId uint) ([]*Album, error) {
	albums, err := uc.Repository.FindByUserId(userId)

	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (uc *UseCase) CreateAlbum(userId uint, title, thumbUrl, gender string) error {
	album, err := NewAlbum(userId, title, thumbUrl, gender)

	if err != nil {
		return err
	}

	err = uc.Repository.Create(album)

	return err
}

func (uc *UseCase) UpdateAlbum(albumId uint, newTitle, newGender string) error {
	album, err := uc.Repository.FindById(albumId)

	if err != nil {
		return err
	}

	if newTitle != "" && newTitle != album.Title() {
		err = album.ChangeTitle(newTitle)
	}

	if newGender != "" && newGender != album.Gender() {
		err = album.ChangeGender(newTitle)
	}

	if err != nil {
		err = uc.Repository.Update(album)
	}

	return err
}

func (uc *UseCase) DestroyAlbum(albumId uint) error {
	err := uc.Repository.Destroy(albumId)

	return err
}
