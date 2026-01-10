package pg

import (
	"database/sql"
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
	query := `
        INSERT INTO users (first_name, email, login, password, user_role) 
        VALUES (:first_name, :email, :login, :password, :user_role) 
        RETURNING id`

	rows, err := u.db.NamedQuery(query, userToSave)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == constraintViolation {
			return nil, fmt.Errorf("user with this email or login already exist: %w", err)
		}
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	defer rows.Close()

	rows.Next()
	if err := rows.Scan(&userToSave.Id); err != nil {
		return nil, fmt.Errorf("failed to scan user %w", err)
	}
	return &userToSave, nil
}

func (u *UserRepository) GetById(userId int) (*model.UserDB, error) {
	userDB := &model.UserDB{}

	query := `SELECT * FROM users where id = $1`
	if err := u.db.Get(userDB, query, userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", userId)
		}
		return nil, fmt.Errorf("failed to get user %w", err)
	}
	return userDB, nil
}
