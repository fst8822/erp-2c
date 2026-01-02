package pg

import "github.com/jmoiron/sqlx"

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u UserRepository) Save() {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) GetById() {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) GetByName() {
	//TODO implement me
	panic("implement me")
}
