package metrics

import (
	"net/http"
	"time"

	"github.com/F0urward/proftwist-backend/internal/metrics"
	"github.com/gorilla/mux"
)

type MetricsMiddleware struct {
	metrics metrics.Metrics
}

func NewMetricsMiddleware(m metrics.Metrics) *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: m,
	}
}

func (m *MetricsMiddleware) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.metrics.IncHTTPInFlight()
		defer m.metrics.DecHTTPInFlight()

		route := mux.CurrentRoute(r)
		path := r.URL.Path
		if route != nil {
			if template, err := route.GetPathTemplate(); err == nil {
				path = template
			}
		}

		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		statusCode := rw.statusCode

		m.metrics.IncHTTPRequest(r.Method, path, statusCode)
		m.metrics.ObserveHTTPDuration(r.Method, path, duration)
		m.metrics.ObserveHTTPResponseSize(r.Method, path, float64(rw.size))

		if statusCode >= 400 {
			m.metrics.IncHTTPError(r.Method, path, statusCode)
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}
