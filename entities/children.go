package entities

type Children struct {
	Base
	Name     string
	ThumbUrl string
	UserID   string
	User     `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
