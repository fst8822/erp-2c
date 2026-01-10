package model

type User struct {
	Id        int64  `json:"id" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Login     string `json:"login" validate:"required"`
	Password  string `json:"password" validate:"required"`
	UserRole  string `json:"user_role,omitempty"`
}

type UserDB struct {
	Id        int64  `db:"id"`
	FirstName string `db:"first_name"`
	Email     string `db:"email"`
	Login     string `db:"login"`
	Password  string `db:"password"`
	UserRole  string `db:"user_role"`
}
