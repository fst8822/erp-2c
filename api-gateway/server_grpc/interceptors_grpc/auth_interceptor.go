package interceptors_grpc

import (
	"api-gateway/model"
	"api-gateway/security"
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthServerInterceptorUnary(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (any, error) {
	const op = "api-gateway.server_grpc.token_interceptor.AuthServerInterceptorUnary"

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is missing")
	}

	ctxNew, err := prepareContext(ctx, md)
	if err != nil {
		return nil, err
	}

	return handler(ctxNew, req)
}

type wrapperCTX struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrapperCTX) Context() context.Context {
	return w.ctx
}
func AuthServerInterceptorStream(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	const op = "api-gateway.server_grpc.token_interceptor.AuthServerInterceptorStream"

	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is missing")
	}

	ctxNew, err := prepareContext(ss.Context(), md)
	if err != nil {
		return err
	}

	return handler(srv, &wrapperCTX{
		ServerStream: ss,
		ctx:          ctxNew,
	})
}

func prepareContext(ctx context.Context, md metadata.MD) (context.Context, error) {
	const op = "api-gateway.server_grpc.token_interceptor.prepareContext"

	authHeader := md.Get(authorization)
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization header missing")
	}

	token := strings.TrimPrefix(authHeader[0], bearer)
	claimsCustom, err := security.ParseToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	claims, ok := claimsCustom.Claims.(*security.CustomClaims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	userID, ok := security.GetUserIdFromClaims(claims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	userRole, ok := security.GetRoleFromClaims(claims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	ctxNew := model.SetUserIDContext(ctx, userID)
	ctxNew = model.SetUserRoleContext(ctxNew, userRole)
	return ctxNew, nil
}
