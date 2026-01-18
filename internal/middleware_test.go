package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(requestIDMiddleware)
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("GeneratesRequestIDWhenMissing", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, req)

		assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
	})

	t.Run("PreservesRequestIDWhenPresent", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		expectedID := "custom-id"
		req.Header.Set("X-Request-ID", expectedID)
		router.ServeHTTP(w, req)

		assert.Equal(t, expectedID, w.Header().Get("X-Request-ID"))
	})
}
