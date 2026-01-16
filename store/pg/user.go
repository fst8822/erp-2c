package pg

import (
	"database/sql"
	"erp-2c/lib/types"
	"erp-2c/model"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sqlx.DB
}

const (
	constraintViolation = "23505"
)

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Save(userToSave model.UserDB) (*model.UserDB, error) {
	query := `INSERT INTO users (first_name, email, login, password, user_role) 
        	  VALUES (:first_name, :email, :login, :password, :user_role) 
              RETURNING id`

	rows, err := u.db.NamedQuery(query, userToSave)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == constraintViolation {
			return nil, types.NewAppErr(fmt.Sprintf("user with email %s or login %s already exist",
				userToSave.Email, userToSave.Login), types.ErrAlreadyExist)
		}
		return nil, types.NewAppErr("inspected SQL error", fmt.Errorf(err.Error(), types.ErrInspectedSQL))
	}
	defer rows.Close()

	rows.Next()
	if err := rows.Scan(&userToSave.Id); err != nil {
		return nil, types.NewAppErr("inspected SQL error", fmt.Errorf(err.Error(), types.ErrInspectedSQL))
	}
	return &userToSave, nil
}

func (u *UserRepository) GetById(userId int64) (*model.UserDB, error) {
	userDB := &model.UserDB{}

	query := `SELECT * FROM users where id = $1`
	if err := u.db.Get(userDB, query, userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NewAppErr(fmt.Sprintf("user with id %d not found", userId), types.ErrNotFound)
		}
		return nil, types.NewAppErr("inspected SQL error", fmt.Errorf(err.Error(), types.ErrInspectedSQL))
	}
	return userDB, nil
}

func (u *UserRepository) GetByLogin(login string) (*model.UserDB, error) {
	userDB := &model.UserDB{}

	query := `SELECT * FROM users where login = $1`
	if err := u.db.Get(userDB, query, login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.NewAppErr(fmt.Sprintf("user with login %s not found", login), types.ErrNotFound)
		}
		return nil, types.NewAppErr("inspected SQL error", fmt.Errorf(err.Error(), types.ErrInspectedSQL))
	}
	return userDB, nil
}
