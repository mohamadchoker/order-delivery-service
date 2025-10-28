package middleware

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingUnaryInterceptor creates a gRPC unary interceptor that logs requests with request ID and status code
func LoggingUnaryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		requestID := GetRequestID(ctx)

		// Call handler
		resp, err := handler(ctx, req)

		// Extract gRPC status code
		grpcStatus := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				grpcStatus = st.Code()
			} else {
				grpcStatus = codes.Unknown
			}
		}

		// Log request
		duration := time.Since(start)
		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("grpc_code", grpcStatus.String()),
			zap.Duration("duration", duration),
			zap.String("request_id", requestID),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Error("gRPC request failed", fields...)
		} else {
			logger.Info("gRPC request completed", fields...)
		}

		return resp, err
	}
}
