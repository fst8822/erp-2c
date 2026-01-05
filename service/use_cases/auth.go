package use_cases

import (
	"erp-2c/store"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	store *store.Store
}

func NewAuthService(store *store.Store) *AuthService {
	return &AuthService{store: store}
}

func (a AuthService) SignUp() {
	//TODO implement me
	panic("implement me")
}

func (a AuthService) SignIn() {
	//TODO implement me
	panic("implement me")
}

func generatePasswordHash(password string) (string, error) {
	const op = "service.use"
	b, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("failed generate password hash %w %s", err, op)
	}
	return string(b), err
}
