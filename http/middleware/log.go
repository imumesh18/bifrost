package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"log/slog"
)

type APILogField struct {
	RequestTime    time.Time
	RequestHeader  http.Header
	ResponseHeader http.Header
	Method         string
	URL            string
	Latency        string
	RequestBody    string
	ResponseBody   string
	ClientIP       string
	UserAgent      string
}

// RequestResponseWriter is a wrapper around http.ResponseWriter that provides
// access to the response body.
type RequestResponseWriter struct {
	http.ResponseWriter
	body bytes.Buffer
}

// Write writes the response body to the buffer and the original http.ResponseWriter.
func (w *RequestResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteHeader writes the response header to the original http.ResponseWriter.
func (w *RequestResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

// Header returns the response header map from the original http.ResponseWriter.
func (w *RequestResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func RequestResponseLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Read request body
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		// Create a new buffer to store the request body
		requestBodyBuffer := bytes.NewBuffer(requestBody)

		// Replace the request body with the buffer so it can be read again later
		r.Body = io.NopCloser(requestBodyBuffer)

		// Call the next handler
		rrw := &RequestResponseWriter{ResponseWriter: w}
		next.ServeHTTP(rrw, r)

		// Calculate latency
		latency := time.Since(start)

		// Read response body
		responseBody, err := io.ReadAll(&rrw.body)
		if err != nil {
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}

		// Create a new buffer to store the response body
		responseBodyBuffer := bytes.NewBuffer(responseBody)

		apiLogFields := APILogField{
			Method:         r.Method,
			URL:            r.URL.String(),
			Latency:        latency.String(),
			RequestBody:    requestBodyBuffer.String(),
			ResponseBody:   responseBodyBuffer.String(),
			RequestTime:    start,
			ClientIP:       r.Header.Get("X-Forwarded-For"),
			UserAgent:      r.Header.Get("User-Agent"),
			RequestHeader:  r.Header,
			ResponseHeader: rrw.Header(),
		}

		slog.Info(
			"",
			slog.String("method", apiLogFields.Method),
			slog.String("url", apiLogFields.URL),
			slog.String("latency", apiLogFields.Latency),
			slog.String("request_body", apiLogFields.RequestBody),
			slog.String("response_body", apiLogFields.ResponseBody),
			slog.Time("request_time", apiLogFields.RequestTime),
			slog.String("client_ip", apiLogFields.ClientIP),
			slog.String("user_agent", apiLogFields.UserAgent),
			slog.Any("request_header", apiLogFields.RequestHeader),
			slog.Any("response_header", apiLogFields.ResponseHeader),
		)
	})
}
