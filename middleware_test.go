package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// function to test the middleware
func TestAuth(t *testing.T) {
	// setup
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/albums/vk/", nil)

	// execute
	router.ServeHTTP(w, req)

	// assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "please provide api_key")
}
