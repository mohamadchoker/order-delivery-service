package middleware

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/mohamadchoker/order-delivery-service/internal/constants"
)

type requestIDKey struct{}

// RequestIDUnaryInterceptor adds a request ID to the context
func RequestIDUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		requestID := extractRequestID(ctx)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add request ID to context
		ctx = context.WithValue(ctx, requestIDKey{}, requestID)

		// Add request ID to outgoing metadata
		if err := grpc.SetHeader(ctx, metadata.Pairs(constants.RequestIDHeader, requestID)); err != nil {
			// Log but don't fail the request
			_ = err
		}

		return handler(ctx, req)
	}
}

// extractRequestID extracts request ID from incoming metadata
func extractRequestID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(constants.RequestIDHeader)
	if len(values) > 0 {
		return values[0]
	}

	return ""
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey{}).(string); ok {
		return requestID
	}
	return ""
}
