package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

// Metrics holds all Prometheus metrics for the gateway service.
type Metrics struct {
	// HTTP Metrics
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	HTTPErrorsTotal     *prometheus.CounterVec

	// gRPC Client Metrics
	GRPCClientRequestsTotal   *prometheus.CounterVec
	GRPCClientRequestDuration *prometheus.HistogramVec
	GRPCClientErrorsTotal     *prometheus.CounterVec

	// Security Metrics
	JWTValidationTotal *prometheus.CounterVec

	Reg *prometheus.Registry
}

// NewMetrics initializes and registers the Prometheus metrics.
func NewMetrics(reg *prometheus.Registry) *Metrics {
	return &Metrics{
		Reg: reg,
		// HTTP Metrics
		HTTPRequestsTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_http_requests_total",
				Help: "Total number of incoming HTTP requests.",
			},
			[]string{"method", "path", "status_code"},
		),
		HTTPRequestDuration: promauto.With(reg).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gateway_http_request_duration_seconds",
				Help:    "Duration of incoming HTTP requests.",
				Buckets: prometheus.DefBuckets, // Default buckets
			},
			[]string{"method", "path"},
		),
		HTTPErrorsTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_http_requests_errors_total",
				Help: "Total number of HTTP requests resulting in an error.",
			},
			[]string{"method", "path", "status_code"},
		),

		// gRPC Client Metrics
		GRPCClientRequestsTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_grpc_client_requests_total",
				Help: "Total number of gRPC requests made to internal services.",
			},
			[]string{"service", "method"},
		),
		GRPCClientRequestDuration: promauto.With(reg).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gateway_grpc_client_request_duration_seconds",
				Help:    "Duration of gRPC requests made to internal services.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "method"},
		),
		GRPCClientErrorsTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_grpc_client_errors_total",
				Help: "Total number of gRPC requests to internal services that resulted in an error.",
			},
			[]string{"service", "method", "grpc_status_code"},
		),

		// Security Metrics
		JWTValidationTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_jwt_validation_total",
				Help: "Total number of JWT validation attempts.",
			},
			[]string{"result"}, // e.g., "success", "token_expired", "invalid_signature"
		),
	}
}

func MetricsMiddleware(metrics *Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ignore the /metrics endpoint itself
		if c.FullPath() == "/metrics" {
			c.Next()
			return
		}

		// Start timer
		start := prometheus.NewTimer(metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()))

		c.Next()

		start.ObserveDuration()

		statusCode := c.Writer.Status()
		metrics.HTTPRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), fmt.Sprintf("%d", statusCode)).Inc()

		if statusCode >= 400 {
			metrics.HTTPErrorsTotal.WithLabelValues(c.Request.Method, c.FullPath(), fmt.Sprintf("%d", statusCode)).Inc()
		}
	}
}

func (m *Metrics) GRPCClientMetricsInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()

	service := cc.Target()

	err := invoker(ctx, method, req, reply, cc, opts...)

	duration := time.Since(start).Seconds()

	st, _ := status.FromError(err)
	grpcStatusCode := st.Code().String()

	m.GRPCClientRequestsTotal.WithLabelValues(service, method).Inc()
	m.GRPCClientRequestDuration.WithLabelValues(service, method).Observe(duration)

	if err != nil {
		m.GRPCClientErrorsTotal.WithLabelValues(service, method, grpcStatusCode).Inc()
	}

	return err
}
