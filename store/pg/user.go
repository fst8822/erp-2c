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

func (u *UserRepository) Save(userToSave model.UserDB) (*model.UserDB, error) {

	u.db.Query("INSERT INTO ")
	userToSave.ID = 1
	return &userToSave, nil
}

func (u *UserRepository) GetById(userId int) (*model.UserDB, error) {
	return &model.UserDB{}, nil
}

func (u *UserRepository) GetByLogin(login string) (*model.UserDB, error) {
	return &model.UserDB{}, nil
}
