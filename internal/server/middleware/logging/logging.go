package logging

import (
	"net/http"
	"time"

	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/pkg/logger"
)

type LoggingMiddleware struct {
	log logger.Logger
}

func NewLoggingMiddleware(log logger.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		log: log,
	}
}

func (l *LoggingMiddleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestLogger := l.log.WithFields(map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"ip":     r.RemoteAddr,
		})

		requestLogger.Info("HTTP request started")

		ctx := ctxutil.WithLogging(r.Context(), requestLogger)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		requestLogger.WithFields(map[string]interface{}{
			"took": time.Since(start).String(),
		}).Info("HTTP request completed")
	})
}
