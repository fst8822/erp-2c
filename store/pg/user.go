package pg

import (
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
	queryFields         = "first_name, email, login, password, user_role"
	queryValues         = ":first_name, :email, :login, :password, :user_role"
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
		if errors.As(err, &pqErr) {
			return nil, fmt.Errorf("user with this email or login already exist: %w", err)
		}
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows iteration error: %w", err)
		}
		return nil, fmt.Errorf("no id returned after insert")
	}

	if err := rows.Scan(&userToSave.ID); err != nil {
		return nil, fmt.Errorf("failed to scan user ID: %w", err)
	}
	return &userToSave, nil
}

func (u *UserRepository) GetById(userId int) (*model.UserDB, error) {
	return &model.UserDB{}, nil
}

func (u *UserRepository) GetByLogin(login string) (*model.UserDB, error) {
	return &model.UserDB{}, nil
}
