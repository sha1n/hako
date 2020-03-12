package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type server struct {
	stopChan   chan bool
	httpServer *http.Server
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
