package service

import "github.com/go-chi/chi/v5"

type AuthService interface {
	SignUp()
	SignIn()
	GetAuthenticatedUserFromContext(c *chi.Context)
}

type UserService interface {
	Save()
	GetById()
	GetByName()
}
type ProductService interface {
	Save()
	GetById()
	GetByName()
	GetAll()
	UpdateById()
	DeleteById()
}
