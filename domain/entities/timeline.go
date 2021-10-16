package entities

type TimeLine struct {
	Base
	SnapUrl    string `json:"snapUrl" gorm:"type:varchar(255);not null"`
	ChildrenID string
	Children   Children `json:"children" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
