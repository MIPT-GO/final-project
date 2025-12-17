package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ResponseWriterWithCaptureCode struct {
	http.ResponseWriter
	statusCode int
}

func (writter *ResponseWriterWithCaptureCode) WriteHeader(code int) {
	writter.statusCode = code
	writter.ResponseWriter.WriteHeader(code)
}

type HttpMetricsCollector struct {
	duration    *prometheus.HistogramVec
	codeCounter *prometheus.CounterVec
}

func (collector *HttpMetricsCollector) New() {
	collector.duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status_code"},
	)

	collector.codeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by status code",
		},
		[]string{"method", "path", "status_code"},
	)

	prometheus.MustRegister(collector.duration)
	prometheus.MustRegister(collector.codeCounter)
}

func (collector *HttpMetricsCollector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		writter := &ResponseWriterWithCaptureCode{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(writter, r)

		duration := time.Since(start).Seconds()
		status := http.StatusText(writter.statusCode)

		collector.duration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
		collector.codeCounter.WithLabelValues(r.Method, r.URL.Path, status).Inc()
	})
}
