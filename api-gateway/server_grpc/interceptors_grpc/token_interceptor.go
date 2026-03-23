package interceptors_grpc

import (
	"api-gateway/lib/sl"
	"api-gateway/model"
	"api-gateway/security"
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorization = "authorization"
	bearer        = "Bearer "
)

func TokenInterceptorUnary() grpc.UnaryClientInterceptor {
	const op = "api-gateway.server_grpc.token_interceptor.TokenInterceptorUnary"

	return func(
		ctx context.Context,
		method string,
		req,
		reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctxNew, err := prepareOutgoingContext(ctx)
		if err != nil {
			slog.Error("Error prepare Auth Context", sl.ErrWithOP(err, op))
			return err
		}
		return invoker(ctxNew, method, req, reply, cc, opts...)
	}
}

func TokenInterceptorStream() grpc.StreamClientInterceptor {
	const op = "api-gateway.server_grpc.token_interceptor.TokenInterceptorStream"
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption) (grpc.ClientStream, error) {

		ctxNew, err := prepareOutgoingContext(ctx)
		if err != nil {
			slog.Error("Error prepare Auth Context", sl.ErrWithOP(err, op))
			return nil, err
		}
		return streamer(ctxNew, desc, cc, method, opts...)
	}
}

func prepareOutgoingContext(ctx context.Context) (context.Context, error) {
	const op = "api-gateway.server_grpc.token_interceptor.prepareAuthContext"

	userID, ok := model.UserIDFromContext(ctx)
	if !ok {
		slog.Error("Missing user id from context", slog.String("op", op))
		return nil, status.Error(codes.Internal, "Missing user id")
	}
	userRole, ok := model.UserRoleFromContext(ctx)
	if !ok {
		slog.Error("Missing user role from context", slog.String("op", op))
		return nil, status.Error(codes.Internal, "Missing user role")
	}
	tokenJWT, err := security.GenerateToken(userID, userRole)
	if err != nil {
		slog.Error("Error generate jwt token", sl.ErrWithOP(err, op))
		return nil, status.Error(codes.Internal, "Error generate jwt token")
	}
	md := metadata.Pairs(authorization, bearer+tokenJWT)
	return metadata.NewOutgoingContext(ctx, md), nil
}
