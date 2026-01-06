package jwt

import "erp-2c/model"

type JWT struct {
	token string
}

func GenerateToken(user model.User) (JWT, error) {
	return JWT{}, nil
}

func ParseToken(jwt string) (int, error) {
	return 0, nil
}

func GetRoleFromJwt() (string, error) {
	return "", nil
}
