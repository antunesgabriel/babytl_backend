package user

type Repository interface {
	Create(*User) error

	Update(*User) error

	FindByEmail(email string) (*User, error)

	FindById(id uint) (*User, error)
}
