package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/diogocarasco/go-pharmacy-service/internal/metrics"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := NewResponseWriter(w)

		next.ServeHTTP(lw, r)

		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(lw.Status())

		metrics.RequestCounter.WithLabelValues(r.Method, r.URL.Path, statusCode).Inc()
		metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path, statusCode).Observe(duration)
	})
}

type ResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (lw *ResponseWriter) WriteHeader(code int) {
	lw.status = code
	lw.ResponseWriter.WriteHeader(code)
}

func (lw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := lw.ResponseWriter.Write(b)
	lw.size += size
	return size, err
}

func (lw *ResponseWriter) Status() int {
	return lw.status
}

func (lw *ResponseWriter) Size() int {
	return lw.size
}
