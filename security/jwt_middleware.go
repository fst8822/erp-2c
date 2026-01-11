package security

import (
	"context"
	"erp-2c/lib/response"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
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
			resp := response.Unauthorized(fmt.Sprintf("Missing  authorization header %s", auth))
			render.Status(r, resp.Code)
			render.JSON(w, r, resp)
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != bearer {
			resp := response.Unauthorized("Invalid  authorization header format")
			render.Status(r, resp.Code)
			render.JSON(w, r, resp)
			return
		}

		token, err := ParseToken(parts[1])
		if err != nil || !token.Valid {
			resp := response.Unauthorized("Invalid token")
			render.Status(r, resp.Code)
			render.JSON(w, r, resp)
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			resp := response.Unauthorized("Invalid claims")
			render.Status(r, resp.Code)
			render.JSON(w, r, resp)
			return
		}

		id, ok := GetUserIdFromClaims(claims)
		if !ok {
			resp := response.Unauthorized("Invalid token: user Id not found")
			render.Status(r, resp.Code)
			render.JSON(w, r, resp)
			return
		}
		role, ok := GetRoleFromClaims(claims)
		if !ok {
			resp := response.Unauthorized("Invalid token: user role not found")
			render.Status(r, resp.Code)
			render.JSON(w, r, resp)
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, id)
		ctx = context.WithValue(ctx, userRole, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
