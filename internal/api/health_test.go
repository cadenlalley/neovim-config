package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHealth(t *testing.T) {
	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call the handler directly
	testApp.API.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Optional: Check response body
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}
