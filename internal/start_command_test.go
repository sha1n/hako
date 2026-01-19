package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"bytes"
	"encoding/json"
	"log"
	"log/slog"

	"github.com/sha1n/hako/test"
	"github.com/stretchr/testify/assert"
)

func TestStartWithDefaults(t *testing.T) {
	config, err := newConfigWith("", 0)
	assert.NoError(t, err)

	router := createGinEngine(config)

	w, req := requestWith("GET", "/echo", "default-echo")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "default-echo", w.Body.String())
}

func TestResponseContentTypeIsEchoedAsWell(t *testing.T) {
	config, err := newConfigWith("", 0)
	assert.NoError(t, err)

	router := createGinEngine(config)
	w, req := requestWith("POST", "/echo", "{'j': 'son'}")
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{'j': 'son'}", w.Body.String())
}

func TestStartWithUndefinedPath(t *testing.T) {
	config, err := newConfigWith("", 0)
	assert.NoError(t, err)

	router := createGinEngine(config)

	w, req := requestWith("GET", "/undefined", "nobody-home")
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestStartWithCustomPath(t *testing.T) {
	config, err := newConfigWith("/custom", 0)
	assert.NoError(t, err)

	router := createGinEngine(config)

	w, req := requestWith("GET", "/custom", "custom-echo")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "custom-echo", w.Body.String())
}

func TestStartWithDelay(t *testing.T) {
	var delay int32 = 100
	config, err := newConfigWith("/delay", delay)
	assert.NoError(t, err)

	router := createGinEngine(config)

	w, req := requestWith("GET", "/delay", "zzz...")
	start := time.Now()
	router.ServeHTTP(w, req)

	assert.True(t, time.Since(start) >= time.Millisecond*100)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "zzz...", w.Body.String())
}

func TestEchoEndpointSanity(t *testing.T) {
	scope := newServerTestScope()
	c, _ := newConfigWith("", 0)
	c.ServerPort = scope.port
	c.Verbose = true
	c.VerboseHeaders = true

	stop := StartAsync(c)
	defer stop()

	assert.Eventually(
		t,
		func() bool {
			if res, err := http.Get(scope.serverUrlWith("/echo")); err == nil {
				return res.StatusCode == 200
			}
			return false
		},
		time.Second*30,
		time.Millisecond*10,
	)
}

func TestConfigureLogging(t *testing.T) {
	// Capture stderr
	var buf bytes.Buffer
	originalDefault := slog.Default()
	originalLogOutput := log.Writer()
	originalLogFlags := log.Flags()
	defer func() {
		slog.SetDefault(originalDefault)
		log.SetOutput(originalLogOutput)
		log.SetFlags(originalLogFlags)
	}()

	t.Run("JSONLogging", func(t *testing.T) {
		config := Config{JSONLog: true}
		configureLoggingWithOutput(config, &buf)

		slog.Info("test message")

		var logEntry map[string]interface{}
		err := json.Unmarshal(buf.Bytes(), &logEntry)
		assert.NoError(t, err, "Output should be valid JSON")
		assert.Equal(t, "test message", logEntry["msg"])
		assert.Equal(t, "INFO", logEntry["level"])
	})

	t.Run("TextLogging", func(t *testing.T) {
		buf.Reset()
		config := Config{JSONLog: false}
		configureLoggingWithOutput(config, &buf)

		slog.Info("test message")

		assert.Contains(t, buf.String(), "level=INFO")
		assert.Contains(t, buf.String(), "msg=\"test message\"")
	})

	t.Run("StandardLogRedirection", func(t *testing.T) {
		buf.Reset()
		config := Config{JSONLog: true}
		configureLoggingWithOutput(config, &buf)

		log.Println("legacy log message")

		var logEntry map[string]interface{}
		err := json.Unmarshal(buf.Bytes(), &logEntry)
		assert.NoError(t, err, "Output should be valid JSON. Content: "+buf.String())
		assert.Equal(t, "legacy log message", logEntry["msg"])
		assert.Equal(t, "INFO", logEntry["level"])
	})
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "/"},
		{"/", "/"},
		{"path", "/path"},
		{"/path", "/path"},
		{" path ", "/path"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, normalizePath(tt.input))
	}
}

func TestNewConfigFromArgs(t *testing.T) {
	config := NewConfigFromArgs(8080, 100, "echo", true, true, true)

	assert.Equal(t, 8080, config.ServerPort)
	assert.Equal(t, int32(100), config.Delay)
	assert.Equal(t, "/echo", config.EchoPath)
	assert.True(t, config.Verbose)
	assert.True(t, config.VerboseHeaders)
	assert.True(t, config.JSONLog)
}

func newConfigWith(path string, delay int32) (config Config, err error) {
	port, err := test.RandomFreePort()

	config = Config{
		ServerPort: port,
		Delay:      delay,
		Verbose:    time.Now().Nanosecond()%2 == 0,
		EchoPath:   path,
	}

	return config, err
}

func requestWith(method string, path string, body string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, strings.NewReader(body))

	return w, req
}
