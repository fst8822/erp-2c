package store

type ProductRepository interface {
	Save()
	GetById()
	GetAll()
	UpdateById()
	DeleteById()
}

type UserRepository interface {
	Save()
	GetById()
	GetByName()
}
