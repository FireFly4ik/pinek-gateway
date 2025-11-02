package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	JWTValidationTotal      *prometheus.CounterVec
	RateLimitThrottledTotal *prometheus.CounterVec

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
		RateLimitThrottledTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_rate_limit_throttled_total",
				Help: "Total number of throttled requests due to rate limiting.",
			},
			[]string{"client_ip", "path"},
		),
	}
}

func MetricsMiddleware(metrics *Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics.HTTPRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), "").Inc()
		c.Next()
	}
}
