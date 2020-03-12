package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sha1n/echo-server/cmd/echoserver/utils"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type message struct {
	Value string `json:"value" binding:"required"`
}

func Test_Stop(t *testing.T) {
	scope := newServerTestScope()
	server := scope.newServerBuilder().Build()
	server.StartAsync()
	scope.awaitPort()

	server.StopNow(time.Second * 3)

	_, err := http.Get(scope.serverUrlWith("/"))
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "connection refused"))
}

func Test_Start(t *testing.T) {
	scope := newServerTestScope()
	server := scope.newServerBuilder().Build()
	defer server.StopAsync()

	server.StartAsync()
	scope.awaitPort()

	res, err := http.Get(scope.serverUrlWith("/"))
	assert.NoError(t, err)
	assert.Equal(t, 404, res.StatusCode)
}

func Test_ServerShouldReturnCode405ForUnmappedMethodsOfExistingResource(t *testing.T) {
	scope := newServerTestScope()
	server := scope.newServerBuilder().
		WithPostHandler("/post", func(ctx *gin.Context) {
			ctx.Status(200)
		}).
		Build()
	defer server.StopAsync()

	server.StartAsync()
	scope.awaitPort()

	res, err := http.Get(scope.serverUrlWith("/post"))
	assert.NoError(t, err)
	assert.Equal(t, 405, res.StatusCode)
}

func Test_GetHandlerShouldWork(t *testing.T) {
	scope := newServerTestScope()
	server := scope.newServerBuilder().
		WithGetHandler("/get", func(ctx *gin.Context) {
			ctx.Status(200)
		}).
		Build()
	defer server.StopAsync()

	server.StartAsync()
	scope.awaitPort()

	res, err := http.Get(scope.serverUrlWith("/get"))
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func Test_EchoHandlerShouldReturnInputMessage(t *testing.T) {
	inputMessage := message{utils.RandomString(10)}
	scope := newServerTestScope()
	server := scope.newServerBuilder().
		WithPostHandler("/echo", echoHandler()).
		Build()
	defer server.StopAsync()

	server.StartAsync()
	scope.awaitPort()

	res, err := http.Post(scope.serverUrlWith("/echo"), "application/json", utils.JsonStringReaderFor(inputMessage))
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, inputMessage, jsonMessageFrom(res))
}

func Test_EchoHandlerShouldFailIfContentTypeIsNotSupported(t *testing.T) {
	scope := newServerTestScope()
	server := scope.newServerBuilder().
		WithPostHandler("/echo", echoHandler()).
		Build()
	defer server.StopAsync()

	server.StartAsync()
	scope.awaitPort()

	res, err := http.Post(scope.serverUrlWith("/echo"), "text/plain", strings.NewReader(utils.RandomString(10)))
	assert.NoError(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

type scope struct {
	port int
}

func newServerTestScope() scope {
	port, _ := utils.RandomFreePort()
	return scope{
		port: port,
	}
}

func (s scope) newServerBuilder() ServerBuilder {
	return NewServer(s.port)
}

func (s scope) serverUrlWith(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", s.port, path)
}

func (s scope) awaitPort() (err error) {
	attemptsLeft := 3

	tryConnect := func() (err error) {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("", strconv.Itoa(s.port)), time.Second*10)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Error while waiting for tcp port %d. Error: %s\r\n", s.port, err))
		} else {
			conn.Close()
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
	json.NewDecoder(response.Body).Decode(&res)

	return res
}
