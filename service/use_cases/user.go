package use_cases

import (
	"erp-2c/lib/sl"
	"erp-2c/model"
	"erp-2c/store"
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

	userDB := model.UserDB{
		FirstName: userToSave.FirstName,
		Email:     userToSave.Email,
		Login:     userToSave.Login,
		Password:  userToSave.Password,
		UserRole:  userToSave.UserRole,
	}
	//todo before saved user, do i need check, exist user invoke method getById and handle error,
	//todo or handle err from u.store.UserRepo.Save (example )
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
	slog.Info("User created", slog.Int64("id", userToSave.ID))
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
