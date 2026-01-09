package model

type User struct {
	ID        int64  `json:"ID"`
	FirstName string `json:"first_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Login     string `json:"login" validate:"required"`
	Password  string `json:"password" validate:"required"`
	UserRole  string `json:"user_role,omitempty"`
}

type UserDB struct {
	ID        int64  `db:"ID"`
	FirstName string `db:"first_name"`
	Email     string `db:"email"`
	Login     string `db:"login"`
	Password  string `db:"password"`
	UserRole  string `db:"user_role"`
}
