package use_cases

import (
	"erp-2c/lib/sl"
	"erp-2c/lib/types"
	"erp-2c/model"
	"erp-2c/security"
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

func (a *AuthService) SignUp(signUp model.SignUp) (*model.UserDomain, error) {
	const op = "service.use_cases.auth.SignUp"

	passHash, err := generatePasswordHash(signUp.Password)
	if err != nil {
		slog.Error("failed generate password hash", slog.String("login", signUp.Login), sl.Err(err))
		return nil, types.NewAppErr("Inspected error", types.ErrPasswordHash)
	}
	signUp.Password = passHash

	saved, err := a.userService.Save(signUp)
	if err != nil {
		return nil, err
	}
	return saved, nil
}

func (a *AuthService) SignIn(signIn model.SignIn) (string, error) {
	const op = "service.use_cases.auth.SignIn"

	userDomain, err := a.userService.GetByLogin(signIn.Login)
	if err != nil {
		return "", err
	}

	if !checkPassword(userDomain.Password, signIn.Password) {
		return "", types.NewAppErr(fmt.Sprintf("user with login %s password invalid", signIn.Login),
			types.ErrUnauthorized)
	}

	token, err := security.GenerateToken(userDomain.Id, userDomain.UserRole)
	if err != nil {
		slog.Error("failed generate Token", slog.String("login", signIn.Login), sl.Err(err))
		return "", types.NewAppErr("Inspected error", types.ErrGeneratedToken)
	}

	return token, nil
}

func generatePasswordHash(password string) (string, error) {
	const op = "service.use_cases.auth.generatePasswordHash"

	b, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
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
