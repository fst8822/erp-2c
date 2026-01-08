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
	return a.userService.Save(userToSave)
}

func (a *AuthService) SignIn(login string, password string) (string, error) {
	return "", nil
}

func generatePasswordHash(password string) (string, error) {
	const op = "service.use.generatePasswordHash"
	slog.Info("Begin get hash password", slog.String("op", op))

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
	const op = "service.use.checkPassword"
	slog.Info("Begin check hash password", slog.String("op", op))

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	return err == nil
}
