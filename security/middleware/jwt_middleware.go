package middleware

import (
	"erp-2c/model"

	"github.com/go-chi/chi/v5"
)

func UserIdentity(user model.User) {

}

func GetUserFromContext(ctx *chi.Context) (*model.User, error) {
	return nil, nil
}

func addUserToContext(user model.User) {

}
