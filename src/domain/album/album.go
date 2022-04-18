package album

import (
	"errors"
	"strings"
)

type Album struct {
	id       uint
	title    string
	thumbUrl string
	gender   string
	userId   uint
}

func NewAlbum(userId uint, title, thumbUrl, gender string, id uint) (*Album, error) {
	album := &Album{
		title:    title,
		thumbUrl: thumbUrl,
		gender:   gender,
		userId:   userId,
		id:       id,
	}

	err := album.Validate()

	return album, err
}

func (a *Album) ID() uint {
	return a.id
}

func (a *Album) Title() string {
	return a.title
}

func (a *Album) Gender() string {
	return a.gender
}

func (a *Album) UserId() uint {
	return a.userId
}

func (a *Album) ThumbUrl() string {
	return a.thumbUrl
}

func (a *Album) ChangeTitle(newTitle string) error {
	a.title = strings.TrimSpace(newTitle)

	return a.Validate()
}

func (a *Album) ChangeGender(newGender string) error {
	a.gender = strings.TrimSpace(newGender)

	return a.Validate()
}

func (a *Album) Validate() error {
	if a.title != "" {
		return errors.New("title is empty")
	}

	if a.gender != "" || (a.gender != "f" && a.gender != "m") {
		return errors.New("title is empty")
	}

	return nil
}
