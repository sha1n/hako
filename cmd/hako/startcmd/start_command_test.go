package startcmd

import (
	"github.com/sha1n/hako/cmd/hako/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_StartWithDefaults(t *testing.T) {
	config, err := newConfigWith("", 0)
	assert.NoError(t, err)

	router := createGinEngine(config)

	w, req := requestWith("GET", "/echo", "default-echo")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "default-echo", w.Body.String())
}

func Test_ResponseContentTypeIsEchoedAsWell(t *testing.T) {
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

func Test_StartWithUndefinedPath(t *testing.T) {
	config, err := newConfigWith("", 0)
	assert.NoError(t, err)

	router := createGinEngine(config)

	w, req := requestWith("GET", "/undefined", "nobody-home")
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func Test_StartWithCustomPath(t *testing.T) {
	config, err := newConfigWith("/custom", 0)
	assert.NoError(t, err)

	router := createGinEngine(config)

	w, req := requestWith("GET", "/custom", "custom-echo")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "custom-echo", w.Body.String())
}

func Test_StartWithDelay(t *testing.T) {
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

func newConfigWith(path string, delay int32) (config Config, err error) {
	port, err := utils.RandomFreePort()

	config = Config{
		ServerPort: port,
		Delay:      delay,
		Verbose:    false,
		EchoPath:   path,
	}

	return config, err
}

func requestWith(method string, path string, body string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, strings.NewReader(body))

	return w, req
}
