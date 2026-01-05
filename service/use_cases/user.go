package use_cases

import "erp-2c/store"

type UserService struct {
	store *store.Store
}

func NewUserService(store *store.Store) *UserService {
	return &UserService{store: store}
}

func (u UserService) Save() {
	//TODO implement me
	panic("implement me")
}

func (u UserService) GetById() {
	//TODO implement me
	panic("implement me")
}

func (u UserService) GetByName() {
	//TODO implement me
	panic("implement me")
}
