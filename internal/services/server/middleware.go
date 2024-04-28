package server

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
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
		var buf bytes.Buffer
		teeBody := io.TeeReader(r.Body, &buf)
		newRequest := r.Clone(r.Context())
		newRequest.Body = io.NopCloser(teeBody)
		next.ServeHTTP(&lw, newRequest)
		duration := time.Since(start)
		s.logger.Infof(
			"request method: %s; uri: %s; duration: %s; request size: %d, content: %s",
			r.Method,
			r.RequestURI,
			duration,
			r.ContentLength,
			buf.String(),
		)

		s.logger.Infof("response status: %d; size: %d;", responseData.status, responseData.size)

	})
}

func (s *Server) gzipCompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isContentTypeAllowed(r.Header.Get("Accept")) {
			next.ServeHTTP(w, r)
			return
		}

		gz := gzip.NewWriter(w)
		defer gz.Close()
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")
		w = &BufferedResponseWriter{w, gz}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) gzipDecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") &&
			isContentTypeAllowed(r.Header.Get("Content-Type")) {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request body", http.StatusInternalServerError)
				return
			}
			defer reader.Close()
			uncompressed, _ := io.ReadAll(reader)
			r.Body = io.NopCloser(bytes.NewReader(uncompressed))
		}

		next.ServeHTTP(w, r)
	})
}

type BufferedResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (w *BufferedResponseWriter) Write(data []byte) (int, error) {
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

func (s *Server) hashCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.conf.GetShaKey() == "" {
			next.ServeHTTP(w, r)
			return
		}

		hash := hmac.New(sha256.New, []byte(s.conf.GetShaKey()))

		var buf bytes.Buffer
		teeBody := io.TeeReader(r.Body, &buf)
		newRequest := r.Clone(r.Context())
		newRequest.Body = io.NopCloser(teeBody)
		calculatedHash := fmt.Sprintf("%x", hash.Sum(buf.Bytes()))

		expectedHash := r.Header.Get("HashSHA256")
		if calculatedHash != expectedHash {
			http.Error(w, "Хеш не совпадает", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, newRequest)
	})
}

func (s *Server) hashResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.conf.GetShaKey() == "" {
			next.ServeHTTP(w, r)
			return
		}

		var buf bytes.Buffer
		w = &BufferedResponseWriter{w, &buf}
		next.ServeHTTP(w, r)

		hash := hmac.New(sha256.New, []byte(s.conf.GetShaKey()))
		calculatedHash := fmt.Sprintf("%x", hash.Sum(buf.Bytes()))
		w.Header().Set("HashSHA256", calculatedHash)
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
		"html/text":              true,
		"plain/text":             true,
		"css/text":               true,
		"xml/text":               true,
	}
	return allowedContentTypes[contentType]
}
