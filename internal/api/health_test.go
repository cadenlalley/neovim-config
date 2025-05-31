package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHealth(t *testing.T) {
	status, body, err := request(http.MethodGet, "/health", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// Check response body
	assert.JSONEq(t, `{"status":"ok"}`, body.String())
}
