package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a gRPC interceptor that validates the authorization bearer token
func UnaryServerInterceptor(v *Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata context")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		tokenStr, err := ExtractBearerToken(authHeaders[0])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header: %v", err)
		}

		claims, err := v.Validate(tokenStr)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// Attach claims to context for business logic to consume
		newCtx := WithClaims(ctx, claims)

		// Proceed to the handler
		return handler(newCtx, req)
	}
}
