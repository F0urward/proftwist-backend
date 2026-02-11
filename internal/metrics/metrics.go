package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics interface {
	IncHTTPRequest(method, path string, statusCode int)
	ObserveHTTPDuration(method, path string, duration time.Duration)
	IncHTTPError(method, path string, statusCode int)
	IncHTTPInFlight()
	DecHTTPInFlight()
	ObserveHTTPResponseSize(method, path string, bytes float64)

	IncGRPCRequest(method, code string)
	ObserveGRPCDuration(method string, duration time.Duration)
	IncGRPCError(method, code string)
	SetGRPCInFlight(count int)

	Handler() http.Handler
}

type MetricsImpl struct {
	registry   *prometheus.Registry
	registerer prometheus.Registerer

	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration  *prometheus.HistogramVec
	httpErrorsTotal      *prometheus.CounterVec
	httpRequestsInFlight prometheus.Gauge
	httpResponseSize     *prometheus.HistogramVec

	grpcRequestsTotal    *prometheus.CounterVec
	grpcRequestDuration  *prometheus.HistogramVec
	grpcErrorsTotal      *prometheus.CounterVec
	grpcRequestsInFlight prometheus.Gauge
}

func NewMetrics(
	reg *prometheus.Registry,
	registerer prometheus.Registerer,
) Metrics {
	factory := promauto.With(reg)

	return &MetricsImpl{
		registry:   reg,
		registerer: registerer,

		httpRequestsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),

		httpRequestDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
			},
			[]string{"method", "path"},
		),

		httpErrorsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_errors_total",
				Help: "Total number of HTTP errors",
			},
			[]string{"method", "path", "status"},
		),

		httpRequestsInFlight: factory.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
		),

		httpResponseSize: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: []float64{100, 1000, 10000, 100000, 1000000, 10000000},
			},
			[]string{"method", "path"},
		),

		grpcRequestsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_requests_total",
				Help: "Total number of gRPC requests",
			},
			[]string{"method", "code"},
		),

		grpcRequestDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grpc_request_duration_seconds",
				Help:    "gRPC request duration in seconds",
				Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
			},
			[]string{"method"},
		),

		grpcErrorsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_errors_total",
				Help: "Total number of gRPC errors",
			},
			[]string{"method", "code"},
		),

		grpcRequestsInFlight: factory.NewGauge(
			prometheus.GaugeOpts{
				Name: "grpc_requests_in_flight",
				Help: "Current number of gRPC requests being processed",
			},
		),
	}
}

func (m *MetricsImpl) IncHTTPRequest(method, path string, statusCode int) {
	m.httpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(statusCode)).Inc()
}

func (m *MetricsImpl) ObserveHTTPDuration(method, path string, duration time.Duration) {
	m.httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

func (m *MetricsImpl) IncHTTPError(method, path string, statusCode int) {
	m.httpErrorsTotal.WithLabelValues(method, path, strconv.Itoa(statusCode)).Inc()
}

func (m *MetricsImpl) IncHTTPInFlight() {
	m.httpRequestsInFlight.Inc()
}

func (m *MetricsImpl) DecHTTPInFlight() {
	m.httpRequestsInFlight.Dec()
}

func (m *MetricsImpl) ObserveHTTPResponseSize(method, path string, bytes float64) {
	m.httpResponseSize.WithLabelValues(method, path).Observe(bytes)
}

func (m *MetricsImpl) IncGRPCRequest(method, code string) {
	m.grpcRequestsTotal.WithLabelValues(method, code).Inc()
}

func (m *MetricsImpl) ObserveGRPCDuration(method string, duration time.Duration) {
	m.grpcRequestDuration.WithLabelValues(method).Observe(duration.Seconds())
}

func (m *MetricsImpl) IncGRPCError(method, code string) {
	m.grpcErrorsTotal.WithLabelValues(method, code).Inc()
}

func (m *MetricsImpl) SetGRPCInFlight(count int) {
	m.grpcRequestsInFlight.Set(float64(count))
}

func (m *MetricsImpl) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
