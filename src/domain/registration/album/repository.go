package album

type Repository interface {
	FindByUserId(userId uint) ([]*Album, error)

	Create(*Album) error

	Update(*Album) error

	Destroy(albumId uint) error

	FindById(id uint) (*Album, error)
}
