package use_cases

import (
	"erp-2c/lib/sl"
	"erp-2c/model"
	"erp-2c/security"
	"erp-2c/service"
	"erp-2c/store"
	"errors"
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

func (a *AuthService) SignUp(signUp model.SignUp) (*model.UserDomain, error) {
	const op = "service.usescases.auth.SignUp"

	passHash, err := generatePasswordHash(signUp.Password)
	if err != nil {
		slog.Error("failed generate password hash", slog.String("login", signUp.Login), sl.ErrWithOP(err, op))
		return nil, fmt.Errorf("user login %s, %w", signUp.Login, err)
	}
	signUp.Password = passHash

	saved, err := a.userService.Save(signUp)
	if err != nil {
		slog.Error("failed to get jwt token", slog.String("login", signUp.Login), sl.ErrWithOP(err, op))

		return nil, fmt.Errorf("failed to save user with login %s, %w", signUp.Login, err)
	}
	return saved, nil
}

func (a *AuthService) SignIn(signIn model.SignIn) (string, error) {
	const op = "service.usescases.auth.SignIn"

	userDomain, err := a.userService.GetByLogin(signIn.Login)
	if err != nil {
		return "", err
	}

	res := checkPassword(userDomain.Password, signIn.Password)
	if !res {
		return "", errors.New(fmt.Sprintf("user password invalid, user login %s", signIn.Login))
	}
	token, err := security.GenerateToken(userDomain.Id, userDomain.UserRole)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generatePasswordHash(password string) (string, error) {
	const op = "service.usesacses.auth.generatePasswordHash"

	b, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("failed generate password hash %w %s", err, op)
	}
	return string(b), err
}

func checkPassword(hashedPassword string, passwordToCheck string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(passwordToCheck),
	)

	return err == nil
}
