package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type Server interface {
	StartAsync()
	StopAsync()
	StopNow(timeout time.Duration) error
}

type server struct {
	stopChan   chan bool
	httpServer *http.Server
}

func NewServer(port int, engine *gin.Engine) Server {
	if engine == nil {
		engine = CreateDefaultRouter()
	}

	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(int(port)),
		Handler: engine,
	}

	s := &server{
		stopChan:   make(chan bool, 1),
		httpServer: httpServer,
	}

	return s
}

func (server *server) StartAsync() {
	log.Printf("Staring HTTP Server on %s", server.httpServer.Addr)

	go func() {
		err := server.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		stop := <-server.stopChan

		log.Println("Received stop signal", stop)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		if err := server.httpServer.Shutdown(ctx); err != nil {
			server.stopChan <- false
			log.Println("Server Shutdown:", err)
		}
		server.stopChan <- true
	}()
}

func (server *server) StopAsync() {
	server.stopChan <- true
}

func (server *server) StopNow(timeout time.Duration) (err error) {
	server.StopAsync()
	timer := time.NewTimer(timeout)
	select {
	case stopped := <-server.stopChan:
		if !stopped {
			err = errors.New("failed to stop server")
		}
	case <-timer.C:
		err = errors.New("timeout waiting for server to stop")
	}
	return err
}

func CreateDefaultRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.HandleMethodNotAllowed = true

	router.Use(func(c *gin.Context) {
		log.Printf(fmt.Sprintf("Handling request at %s", c.Request.RequestURI))
	})

	return router
}
