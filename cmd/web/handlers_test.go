package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.victorsmith.dev/internal/assert"
)

func TestPing (t *testing.T) {
	// Init httptest ResponceRecorder 
	// => this is a proxy for ResponceWriter used for testing
	rr := httptest.NewRecorder()

	// Init a dummy http.Request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Call ping
	ping(rr, r)

	// Get result (http.Response) from rr
	rs := rr.Result()

	// Compare status codes
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// Check resp body
	defer rs.Body.Close()
	
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		// will mark the test as failed, log the error, and then completely 
		// stop execution of the current test (or sub-test).
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "Ok")
}
