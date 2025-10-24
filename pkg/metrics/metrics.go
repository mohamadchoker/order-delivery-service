package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/company/order-delivery-service/internal/constants"
)

var (
	// RequestsTotal counts total number of gRPC requests
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: constants.MetricsNamespace,
			Subsystem: constants.MetricsSubsystem,
			Name:      "grpc_requests_total",
			Help:      "Total number of gRPC requests",
		},
		[]string{"method", "code"},
	)

	// RequestDuration tracks request duration in seconds
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: constants.MetricsNamespace,
			Subsystem: constants.MetricsSubsystem,
			Name:      "grpc_request_duration_seconds",
			Help:      "Duration of gRPC requests in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// ActiveRequests tracks number of active requests
	ActiveRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: constants.MetricsNamespace,
			Subsystem: constants.MetricsSubsystem,
			Name:      "grpc_requests_active",
			Help:      "Number of active gRPC requests",
		},
		[]string{"method"},
	)

	// DeliveryAssignmentsTotal counts delivery assignments by status
	DeliveryAssignmentsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: constants.MetricsNamespace,
			Subsystem: constants.MetricsSubsystem,
			Name:      "delivery_assignments_total",
			Help:      "Total number of delivery assignments",
		},
		[]string{"status", "operation"},
	)

	// DatabaseQueriesTotal counts database queries
	DatabaseQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: constants.MetricsNamespace,
			Subsystem: constants.MetricsSubsystem,
			Name:      "database_queries_total",
			Help:      "Total number of database queries",
		},
		[]string{"operation", "status"},
	)

	// DatabaseQueryDuration tracks database query duration
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: constants.MetricsNamespace,
			Subsystem: constants.MetricsSubsystem,
			Name:      "database_query_duration_seconds",
			Help:      "Duration of database queries in seconds",
			Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"operation"},
	)
)

// MetricsUnaryInterceptor creates a gRPC interceptor for Prometheus metrics
func MetricsUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Increment active requests
		ActiveRequests.WithLabelValues(info.FullMethod).Inc()
		defer ActiveRequests.WithLabelValues(info.FullMethod).Dec()

		// Call handler
		resp, err := handler(ctx, req)

		// Record metrics
		duration := time.Since(start).Seconds()
		code := status.Code(err).String()

		RequestsTotal.WithLabelValues(info.FullMethod, code).Inc()
		RequestDuration.WithLabelValues(info.FullMethod).Observe(duration)

		return resp, err
	}
}

// RecordDeliveryOperation records a delivery assignment operation
func RecordDeliveryOperation(operation, status string) {
	DeliveryAssignmentsTotal.WithLabelValues(status, operation).Inc()
}

// RecordDatabaseQuery records a database query with timing
func RecordDatabaseQuery(operation string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}

	DatabaseQueriesTotal.WithLabelValues(operation, status).Inc()
	DatabaseQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}
