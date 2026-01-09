package use_cases

import (
	"erp-2c/lib/sl"
	"erp-2c/lib/types"
	"erp-2c/model"
	"erp-2c/store"
	"fmt"
	"log/slog"
)

type UserService struct {
	store *store.Store
}

func NewUserService(store *store.Store) *UserService {
	return &UserService{store: store}
}

func (u *UserService) Save(userToSave model.User) (*model.User, error) {
	const op = "service.usecase.user.SAVE"
	slog.With("op", op)

	_, err := u.store.UserRepo.GetByLogin(userToSave.Login)
	if err != nil {
		return nil, fmt.Errorf("%w, %s", types.ErrAlreadyExist, op)
	}

	userDB := model.UserDB{
		FirstName: userToSave.FirstName,
		Email:     userToSave.Email,
		Login:     userToSave.Login,
		Password:  userToSave.Password,
		UserRole:  userToSave.UserRole,
	}

	saved, err := u.store.UserRepo.Save(userDB)
	if err != nil {
		slog.Error("Failed to save user", slog.String("OP", op), sl.Err(err))
		return nil, err
	}

	user := &model.User{
		ID:        saved.ID,
		FirstName: saved.FirstName,
		Email:     saved.Email,
		Login:     saved.Login,
		UserRole:  saved.UserRole,
	}
	return user, nil
}

func (u *UserService) GetById(userId int) (*model.User, error) {
	const op = "service.usecase.user.GetById"
	slog.With("op", op)

	return nil, nil
}

func (u *UserService) GetByName(userName string) (*model.User, error) {
	const op = "service.usecase.user.GetByName"
	slog.With("op", op)

	return nil, nil
}
