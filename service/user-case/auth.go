package user_case

import "erp-2c/store"

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
