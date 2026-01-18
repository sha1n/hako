package internal

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandlerEchoesBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := "hello world"
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "text/plain")

	h := echoHandlerWith(false, false, 0)
	h(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, body, w.Body.String())
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
}

func TestHandlerDynamicStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		statusHeader   string
		expectedStatus int
	}{
		{
			name:           "NoHeader",
			statusHeader:   "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ValidHeader",
			statusHeader:   "418",
			expectedStatus: http.StatusTeapot,
		},
		{
			name:           "InvalidHeader",
			statusHeader:   "invalid",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest("GET", "/", nil)
			if tt.statusHeader != "" {
				c.Request.Header.Set("X-Hako-Status", tt.statusHeader)
			}

			h := echoHandlerWith(false, false, 0)
			h(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestHandlerDelay(t *testing.T) {
    // Basic test to ensure it runs without panic, timing is flaky in unit tests
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	h := echoHandlerWith(false, false, 10)
	h(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
