package security

import (
	"context"
	"erp-2c/lib/response"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const (
	authorization            = "Authorization"
	bearer                   = "Bearer"
	userIdKey     contextKey = "id"
	userRole      contextKey = "role"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get(authorization)
		if auth == "" {
			response.Unauthorized(fmt.Sprintf("Missing  authorization header")).SendResponse(w, r)
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != bearer {
			response.Unauthorized("Invalid  authorization header format").SendResponse(w, r)
			return
		}

		token, err := ParseToken(parts[1])
		if err != nil || !token.Valid {
			response.Unauthorized("Invalid token").SendResponse(w, r)
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			response.Unauthorized("Invalid claims").SendResponse(w, r)
			return
		}

		id, ok := GetUserIdFromClaims(claims)
		if !ok {
			response.Unauthorized("Invalid token: user with Id not found").SendResponse(w, r)
			return
		}
		role, ok := GetRoleFromClaims(claims)
		if !ok {
			response.Unauthorized("Invalid token: user role not found").SendResponse(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, id)
		ctx = context.WithValue(ctx, userRole, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
