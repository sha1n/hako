package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sha1n/hako/test"
	"github.com/sirupsen/logrus"
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
	// Reset formatter after test
	defer logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})

	config := Config{JSONLog: true}
	configureLogging(config)

	_, ok := logrus.StandardLogger().Formatter.(*logrus.JSONFormatter)
	assert.True(t, ok, "Expected JSONFormatter")
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
