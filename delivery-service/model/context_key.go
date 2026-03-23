package model

import "context"

type contextKey string

const (
	userIdKey   contextKey = "id"
	userRoleKey contextKey = "role"
)

func SetUserIDContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIdKey, userID)
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIdKey).(int64)
	return userID, ok
}

func SetUserRoleContext(ctx context.Context, userRole string) context.Context {
	return context.WithValue(ctx, userRoleKey, userRole)
}

func UserRoleFromContext(ctx context.Context) (string, bool) {
	userRole, ok := ctx.Value(userRoleKey).(string)
	return userRole, ok
}
