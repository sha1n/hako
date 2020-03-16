package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sha1n/hako/cmd/hako/http"
	"github.com/sha1n/hako/cmd/hako/utils"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	ServerPort int
	EchoPath   string
	Verbose    bool
	Delay      int32
}

func StartAsync(config Config) {
	server := createHttpServer(config)
	server.StartAsync()
}

func Start(config Config) {
	StartAsync(config)

	awaitShutdownSig()
}

func awaitShutdownSig() {
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Waiting for shutdown signal...")

	<-quitChannel
}

func createHttpServer(config Config) http.Server {
	router := createGinEngine(config)

	server := http.NewServer(config.ServerPort, router)

	stopServerAsync := func() {
		server.StopAsync()
	}

	log.Println("Registering signal listeners for graceful HTTP server shutdown..")
	utils.RegisterShutdownHook(utils.NewSignalHook(syscall.SIGTERM, stopServerAsync))
	utils.RegisterShutdownHook(utils.NewSignalHook(syscall.SIGKILL, stopServerAsync))

	return server
}

func createGinEngine(config Config) *gin.Engine {
	router := http.CreateDefaultRouter()
	registerHandlers(router, "/echo", echoHandlerWith(config.Verbose, config.Delay))

	if config.EchoPath != "" {
		registerHandlers(router, config.EchoPath, echoHandlerWith(config.Verbose, config.Delay))
	}

	return router
}

func registerHandlers(router *gin.Engine, path string, handler func(ctx *gin.Context)) {
	router.GET(path, handler)
	router.POST(path, handler)
	router.PUT(path, handler)
	router.DELETE(path, handler)
	router.PATCH(path, handler)
	router.HEAD(path, handler)
	router.OPTIONS(path, handler)
}

func echoHandlerWith(verbose bool, delay int32) func(*gin.Context) {
	maybeDelay := func() {
		if delay > 0 {
			log.Printf("Delaying response in %d millis", delay)
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}

	return func(c *gin.Context) {
		// todo: in case verbose is off, we definitely don't have to read the entire body into memory.
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Failed to read request body:", err)
		} else {
			// todo: this ignore encoding and assumes the body is printable.. maybe we can do better.
			if verbose && len(bodyBytes) > 0 {
				log.Printf("Body: %s\n\r", string(bodyBytes))
			}
		}

		maybeDelay()

		_, err = c.Writer.Write(bodyBytes)
		if err != nil {
			log.Println("Failed to echo request body:", err)
		}
	}
}
