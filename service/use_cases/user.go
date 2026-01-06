package use_cases

import (
	"erp-2c/model"
	"erp-2c/store"
)

type UserService struct {
	store *store.Store
}

func NewUserService(store *store.Store) *UserService {
	return &UserService{store: store}
}

func (u *UserService) Save(userToSave model.User) (model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) GetById(userId int) (model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) GetByName(userName string) (model.User, error) {
	//TODO implement me
	panic("implement me")
}
