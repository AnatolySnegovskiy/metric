package server

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (w *gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

type sha256ResponseWriter struct {
	http.ResponseWriter
	key string
}

func (w *sha256ResponseWriter) Write(data []byte) (int, error) {
	var buf bytes.Buffer
	buf.Write(data)
	hash := hmac.New(sha256.New, []byte(w.key))
	calculatedHash := fmt.Sprintf("%x", hash.Sum(buf.Bytes()))
	w.ResponseWriter.Header().Set("HashSHA256", calculatedHash)

	return w.ResponseWriter.Write(data)
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.responseData.size += size
	return size, err
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.responseData.status = statusCode
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
		w = &gzipResponseWriter{w, gz}

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

		expectedHash := r.Header.Get("HashSHA256")

		if expectedHash == "" {
			next.ServeHTTP(w, r)
			return
		}

		hash := hmac.New(sha256.New, []byte(s.conf.GetShaKey()))
		body, _ := io.ReadAll(r.Body)
		hash.Write(body)
		expectedHashBytes := []byte(expectedHash)
		calculatedHash := hash.Sum(nil)

		// TODO (https://github.com/AnatolySnegovskiy/metric/issues/16): FIX ME
		if hmac.Equal(expectedHashBytes, calculatedHash) {
			log.Println(expectedHash)
			log.Printf("%x", calculatedHash)
			http.Error(w, "bad hash value", http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	})
}

func (s *Server) hashResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.conf.GetShaKey() == "" {
			next.ServeHTTP(w, r)
			return
		}

		w = &sha256ResponseWriter{w, s.conf.GetShaKey()}
		next.ServeHTTP(w, r)
	})
}

// DecryptMessageMiddleware is a middleware function to decrypt the message before passing it to the next handler.
func (s *Server) DecryptMessageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cryptoKey := s.conf.GetCryptoKey()
		if cryptoKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		bodyData, _ := io.ReadAll(r.Body)
		decryptedBody, err := DecryptionFunction(bodyData, cryptoKey)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))
		next.ServeHTTP(w, r)
	})
}

// DecryptionFunction is a function to decrypt the message before passing it to the next handler.
func DecryptionFunction(data []byte, privateKeyPath string) ([]byte, error) {
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	decryptedMessage, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return nil, err
	}

	return decryptedMessage, nil
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
