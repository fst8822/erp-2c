package use_cases

import (
	"api-gateway/lib/sl"
	"api-gateway/model"
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type userRepositoryInt interface {
	Save(userToSave model.UserDB) (*model.UserDB, error)
	GetById(userId int64) (*model.UserDB, error)
	GetByLogin(login string) (*model.UserDB, error)
	BeginTxx(ctx context.Context) (*sqlx.Tx, error)
}

type UserService struct {
	userRepo userRepositoryInt
}

func NewUserService(userRepo userRepositoryInt) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) Save(userToSave model.SignUp) (*model.UserDomain, error) {
	const op = "service.use_cases.user.SAVE"

	userDB := model.UserDB{
		FirstName: userToSave.FirstName,
		Email:     userToSave.Email,
		Login:     userToSave.Login,
		Password:  userToSave.Password,
		UserRole:  userToSave.UserRole,
	}

	saved, err := u.userRepo.Save(userDB)
	if err != nil {
		slog.Error("Failed to save user", sl.ErrWithOP(err, op))
		return nil, err
	}

	user := model.UserDomain{
		Id:        saved.Id,
		FirstName: saved.FirstName,
		Email:     saved.Email,
		Login:     saved.Login,
		UserRole:  saved.UserRole,
	}

	slog.Info("User created", slog.Int64("id", user.Id))
	return &user, nil
}

func (u *UserService) GetById(userId int64) (*model.UserDomain, error) {
	const op = "service.use_cases.user.GetById"

	found, err := u.userRepo.GetById(userId)
	if err != nil {
		slog.Error("failed to find user", sl.ErrWithOP(err, op))
		return nil, err
	}

	user := &model.UserDomain{
		Id:        found.Id,
		FirstName: found.FirstName,
		Email:     found.Email,
		Login:     found.Login,
		UserRole:  found.UserRole,
	}

	return user, nil
}

func (u *UserService) GetByLogin(userLogin string) (*model.UserDomain, error) {
	const op = "service.use_cases.user.GetWithItemsById"

	found, err := u.userRepo.GetByLogin(userLogin)
	if err != nil {
		slog.Error("failed to find user", sl.ErrWithOP(err, op))
		return nil, err
	}

	user := model.UserDomain{
		Id:        found.Id,
		FirstName: found.FirstName,
		Email:     found.Email,
		Login:     found.Login,
		Password:  found.Password,
		UserRole:  found.UserRole,
	}

	return &user, nil
}
