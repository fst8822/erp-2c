package use_cases

import (
	"erp-2c/lib/sl"
	"erp-2c/model"
	"erp-2c/store"
	"log/slog"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	store    *store.Store
	validate *validator.Validate
}

func NewUserService(store *store.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (u *UserService) Save(userToSave model.SignUp) (*model.UserDomain, error) {
	const op = "service.usecase.user.SAVE"

	userDB := model.UserDB{
		FirstName: userToSave.FirstName,
		Email:     userToSave.Email,
		Login:     userToSave.Login,
		Password:  userToSave.Password,
		UserRole:  userToSave.UserRole,
	}

	saved, err := u.store.UserRepo.Save(userDB)
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
	const op = "service.usecase.user.GetById"

	found, err := u.store.UserRepo.GetById(userId)
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
	const op = "service.usecase.user.GetById"

	found, err := u.store.UserRepo.GetByLogin(userLogin)
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
