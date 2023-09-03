package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.victorsmith.dev/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	// Init a response recorder and a request
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mock http handler -> Will be passed to mw (secureHeaders)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// secureHeaders returns a http.Handler
	secureHeaders(next).ServeHTTP(rr, r)

	// Get results of test from responseRecorder
	rs := rr.Result()

	// Check that the middleware has correctly set the Content-Security-Policy header on the response.
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	// Check that the middleware has correctly set the X-Content-Type-Options header on the response.
	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	// Check that the middleware has correctly set the X-Frame-Options header on the response.
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	// Check that the middleware has correctly set the X-XSS-Protection header on the response
	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	// Check that the middleware has correctly called the next handler in line
	// and the response status code and body are as expected.
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}
	// Remove white space
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
