package pg

import (
	"erp-2c/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u UserRepository) Save(userToSave model.User) (model.UserDB, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) GetById(userId int) (model.UserDB, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) GetByName(userName string) (model.UserDB, error) {
	//TODO implement me
	panic("implement me")
}
