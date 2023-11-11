package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestResponseLoggingMiddleware(t *testing.T) {
	// Create a mock HTTP handler
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, world!")
	})

	// Create a new request
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Create a new middleware handler
	handler := RequestResponseLoggingMiddleware(mockHandler)

	// Call the middleware handler
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Hello, world!\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
