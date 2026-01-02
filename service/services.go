package service

type AuthService interface {
	SignUp()
	SignIn()
}

type UserService interface {
	Save()
	GetById()
	GetByName()
}
type ProductService interface {
	Save()
	GetById()
	GetAll()
	UpdateById()
	DeleteById()
}
