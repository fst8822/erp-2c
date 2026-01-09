package use_cases

import (
	"erp-2c/model"
	"erp-2c/service"
	"erp-2c/store"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	store       *store.Store
	userService service.UserService
}

func NewAuthService(store *store.Store, userService service.UserService) *AuthService {
	return &AuthService{
		store:       store,
		userService: userService,
	}
}

func (a *AuthService) SignUp(userToSave model.User) (*model.User, error) {
	//todo error needs to be handle, or use slog?
	const op = "service.usescases.auth.SignUp"
	slog.With("op", op)

	passHash, err := generatePasswordHash(userToSave.Password)
	//todo where to handle the error or use logging: 1. where i call method or in method which we call
	if err != nil {
		return nil, err
	}
	userToSave.Password = passHash
	return a.userService.Save(userToSave)
}

func (a *AuthService) SignIn(login string, password string) (string, error) {
	const op = "service.usescases.auth.SignIn"
	slog.With("op", op)
	return "", nil
}

func generatePasswordHash(password string) (string, error) {
	const op = "service.usesacses.auth.generatePasswordHash"
	slog.With("op", op)

	b, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("failed generate password hash %w %s", err, op)
	}
	return string(b), err
}
func checkPassword(password string, hash string) bool {
	const op = "service.use.auth.checkPassword"
	slog.With("op", op)

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	return err == nil
}
