package server

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
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
		next.ServeHTTP(&lw, r)
		duration := time.Since(start)

		s.logger.Infof(
			"request method: %s; uri: %s; duration: %s; request size: %d,",
			r.Method,
			r.RequestURI,
			duration,
			r.ContentLength,
		)

		s.logger.Infof("response status: %d; size: %d;", responseData.status, responseData.size)
	})
}

func (s *Server) gzipResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			if !isContentTypeAllowed(w.Header().Get("Content-Type")) {
				next.ServeHTTP(w, r)
				return
			}

			gz := gzip.NewWriter(w)
			defer gz.Close()
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Del("Content-Length")
			w = &gzipResponseWriter{ResponseWriter: w, Writer: gz}
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) gzipRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") &&
			isContentTypeAllowed(r.Header.Get("Content-Type")) {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request body", http.StatusBadRequest)
				return
			}
			defer reader.Close()

			uncompressed, err := io.ReadAll(reader)
			if err != nil {
				http.Error(w, "Failed to read decompressed request body", http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(uncompressed))
		}

		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func (s *Server) JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isContentTypeAllowed(contentType string) bool {
	allowedContentTypes := map[string]bool{
		"text/plain":             true,
		"text/html":              true,
		"text/css":               true,
		"text/xml":               true,
		"application/javascript": true,
		"application/json":       true,
	}
	return allowedContentTypes[contentType]
}
