package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

func requestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, r)

			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}

			elapsed := time.Now().Sub(start).Nanoseconds()
			uri := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

			ww.Header().Set("X-Elapsed-NS", strconv.FormatInt(elapsed, 10))

			log.WithFields(log.Fields{
				"http_method":       r.Method,
				"resp_status":       ww.Status(),
				"resp_time_ms":      float64(elapsed) / 1000000.0,
				"resp_bytes_length": ww.BytesWritten(),
				"uri":               uri,
				"user_agent":        r.UserAgent(),
			}).Info("request completed")
		}
		return http.HandlerFunc(fn)
	}
}
