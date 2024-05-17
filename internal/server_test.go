package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sha1n/hako/test"
	"github.com/stretchr/testify/assert"
)

type message struct {
	Value string `json:"value" binding:"required"`
}

func TestStop(t *testing.T) {
	scope := newServerTestScope()
	server := scope.newServer(NewDefaultEngine())
	server.StartAsync()
	assert.NoError(t, scope.awaitPort())

	assert.NoError(t, server.StopNow(time.Second*3))

	_, err := http.Get(scope.serverUrlWith("/"))
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "connection refused"))
}

func TestStart(t *testing.T) {
	scope := newServerTestScope()
	engine := NewDefaultEngine()
	engine.GET("/", func(ctx *gin.Context) { ctx.JSON(200, nil) })
	server := scope.newServer(engine)
	defer server.StopAsync()

	server.StartAsync()
	assert.NoError(t, scope.awaitPort())

	res, err := http.Get(scope.serverUrlWith("/"))
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestHttpServiceShouldWork(t *testing.T) {
	inputMessage := message{test.RandomStringN(10)}
	scope := newServerTestScope()
	server := scope.newServer(engineWithPostHandler("/echo", echoHandler()))
	defer server.StopAsync()

	server.StartAsync()
	assert.NoError(t, scope.awaitPort())

	res, err := http.Post(scope.serverUrlWith("/echo"), "application/json", jsonStringReaderFor(inputMessage))
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, inputMessage, jsonMessageFrom(res))
}

type scope struct {
	port int
}

func newServerTestScope() scope {
	port, _ := test.RandomFreePort()
	return scope{
		port: port,
	}
}

func (s scope) newServer(engine *gin.Engine) Server {
	return NewServer(s.port, engine)
}

func (s scope) serverUrlWith(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", s.port, path)
}

func (s scope) awaitPort() (err error) {
	attemptsLeft := 3

	tryConnect := func() (err error) {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("", strconv.Itoa(s.port)), time.Second*10)
		if err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("Error while waiting for tcp port %d. Error: %s\r\n", s.port, err))
		} else {
			_ = conn.Close()
		}

		return err
	}

	for attemptsLeft > 0 {
		attemptsLeft--
		err = tryConnect()
		time.Sleep(time.Millisecond * 10)
	}

	return err
}

func engineWithPostHandler(path string, handler func(ctx *gin.Context)) *gin.Engine {
	router := NewDefaultEngine()
	router.POST(path, handler)

	return router
}

func echoHandler() func(*gin.Context) {
	return func(ctx *gin.Context) {
		var input message
		if ctx.BindJSON(&input) == nil {
			ctx.JSON(200, message{input.Value})
		} else {
			ctx.Status(400)
		}
	}
}

func jsonMessageFrom(response *http.Response) (res message) {
	_ = json.NewDecoder(response.Body).Decode(&res)

	return res
}

func jsonStringReaderFor(o interface{}) io.Reader {
	bytes, _ := json.Marshal(o)
	return strings.NewReader(string(bytes))
}
