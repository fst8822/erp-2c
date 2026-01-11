package security

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// todo move env
var secretKey = []byte("2cc801953a01607bb319c9cd5d6f131d29be53cdc69a8793acda750372f21672")

type CustomClaims struct {
	jwt.RegisteredClaims
	*CustomerInfo
}
type CustomerInfo struct {
	Id   int64  `json:"id"`
	Role string `json:"role"`
}

func GenerateToken(userId int64, userRole string) (string, error) {
	slog.Info("Call generate token jwt for",
		slog.Int64("UserId", userId), slog.String("UserRole", userRole))

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		CustomerInfo: &CustomerInfo{
			Id:   userId,
			Role: userRole,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func ParseToken(tokenStr string) (*jwt.Token, error) {
	slog.Info("Begin parse token")

	return jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})
}

func GetRoleFromClaims(claims *CustomClaims) (string, bool) {
	slog.Info("Call GetRoleFromClaims ")

	if claims == nil || claims.CustomerInfo == nil {
		return "", false
	}
	return claims.Role, true
}

func GetUserIdFromClaims(claims *CustomClaims) (int64, bool) {
	slog.Info("Call GetUserIdFromClaims ")

	if claims == nil || claims.CustomerInfo == nil {
		return 0, false
	}
	return claims.Id, true
}
