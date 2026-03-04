package app_metrics

import (
	"net/http"
	"strconv"
	"time"
)

type ResponseWithStatus struct {
	http.ResponseWriter
	status int
}

func (r ResponseWithStatus) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
func (r ResponseWithStatus) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		resp := ResponseWithStatus{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		next.ServeHTTP(resp, r)
		duration := time.Since(start).Seconds()
		HttpReqTotal.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(resp.status)).Inc()
		HttpReqDuration.WithLabelValues(r.URL.Path).Observe(duration)
	})
}
