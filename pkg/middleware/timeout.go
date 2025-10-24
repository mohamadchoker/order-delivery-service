package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/company/order-delivery-service/internal/constants"
)

// TimeoutUnaryInterceptor adds a timeout to each request
func TimeoutUnaryInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	if timeout == 0 {
		timeout = constants.DefaultContextTimeout
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// Channel to handle response
		type result struct {
			resp interface{}
			err  error
		}
		resultChan := make(chan result, 1)

		// Execute handler in goroutine
		go func() {
			resp, err := handler(ctx, req)
			resultChan <- result{resp: resp, err: err}
		}()

		// Wait for result or timeout
		select {
		case res := <-resultChan:
			return res.resp, res.err
		case <-ctx.Done():
			return nil, status.Error(codes.DeadlineExceeded, "request timeout exceeded")
		}
	}
}
