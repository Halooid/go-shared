package logging

import (
	"context"
	"log"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a gRPC interceptor that logs info for all calls and stacktraces for errors
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			st, ok := status.FromError(err)
			// Only log stacktraces for unexpected errors (Internal or Unknown)
			if !ok || st.Code() == codes.Unknown || st.Code() == codes.Internal || st.Code() == codes.DataLoss || st.Code() == codes.DeadlineExceeded {
				log.Printf("gRPC Error: %v\nMethod: %s\nDuration: %v\nStacktrace:\n%s", err, info.FullMethod, duration, debug.Stack())
			} else {
				log.Printf("gRPC Error: %v\nMethod: %s\nDuration: %v", err, info.FullMethod, duration)
			}
		} else {
			log.Printf("gRPC Success: %s\nDuration: %v", info.FullMethod, duration)
		}

		return resp, err
	}
}
