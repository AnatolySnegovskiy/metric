package server

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func (s *Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseData := &responseData{
			status: 200,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		start := time.Now()

		if r.ContentLength == 0 {
			r.Body = http.NoBody
		}

		b, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		r.Body = io.NopCloser(bytes.NewBuffer(b))
		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		s.logger.Infof(
			"request method: %s; uri: %s; duration: %s; request size: %d, request body: %s",
			r.Method,
			r.RequestURI,
			duration,
			r.ContentLength,
			string(b),
		)

		s.logger.Infof("response status: %d; size: %d;", responseData.status, responseData.size)
	})
}

func (s *Server) JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "bad request", http.StatusBadRequest)
		}
		next.ServeHTTP(w, r)
	})
}
