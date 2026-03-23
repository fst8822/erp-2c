package model

type SignUp struct {
	FirstName string `json:"first_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Login     string `json:"login" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required"`
	UserRole  string `json:"user_role" validate:"required"`
}

type SignIn struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserDomain struct {
	Id        int64
	FirstName string
	Email     string
	Login     string
	Password  string
	UserRole  string
}
type UserDB struct {
	Id        int64  `db:"id"`
	FirstName string `db:"first_name"`
	Email     string `db:"email"`
	Login     string `db:"login"`
	Password  string `db:"password"`
	UserRole  string `db:"user_role"`
}
