package model

type SignUp struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Login     string `json:"login"`
	Password  string `json:"password,omitempty"`
	UserRole  string `json:"user_role"`
}

type SignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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
