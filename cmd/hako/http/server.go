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

// Server is a formal HTTP interface
type Server interface {
	StartAsync()
	StopAsync()
	StopNow(timeout time.Duration) error
}

type server struct {
	stopChan   chan bool
	httpServer *http.Server
}

// NewServer creates a new Server and returns it
func NewServer(port int, engine *gin.Engine) Server {
	if engine == nil {
		engine = NewDefaultEngine()
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

// StartAsync starts the server and returns immediately
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

// StopAsync sends a stop signal to the server and returns immediately
func (server *server) StopAsync() {
	server.stopChan <- true
}

// StopNow sends a stop signal to the server and waits for it to stop
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

// NewDefaultEngine creates and returns a new default Gin router
func NewDefaultEngine() *gin.Engine {
	router := gin.Default()
	router.HandleMethodNotAllowed = true

	router.Use(basicRequestLoggerMiddleware)

	return router
}

// NewSilentEngine creates and returns a new basic Gin router with no logging
func NewSilentEngine() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.HandleMethodNotAllowed = true

	router.Use(basicRequestLoggerMiddleware)

	return router
}

func basicRequestLoggerMiddleware(c *gin.Context) {
	log.Printf(fmt.Sprintf("Handling request at %s", c.Request.RequestURI))
}
