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
)

func Start(port int, additionalEchoPath string, verbose bool) {
	server := createHttpServer(port, additionalEchoPath, verbose)
	server.StartAsync()

	awaitShutdownSig()
}

func awaitShutdownSig() {
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Waiting for shutdown signal...")

	<-quitChannel
}

func createHttpServer(port int, additionalEchoPath string, verbose bool) http.Server {
	router := http.CreateDefaultRouter()
	registerHandlers(router, "/echo", echoHandlerWith(verbose))

	if additionalEchoPath != "" {
		registerHandlers(router, additionalEchoPath, echoHandlerWith(verbose))
	}

	server := http.NewServer(port, router)

	stopServerAsync := func() {
		server.StopAsync()
	}

	log.Println("Registering signal listeners for graceful HTTP server shutdown..")
	utils.RegisterShutdownHook(utils.NewSignalHook(syscall.SIGTERM, stopServerAsync))
	utils.RegisterShutdownHook(utils.NewSignalHook(syscall.SIGKILL, stopServerAsync))

	return server
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

func echoHandlerWith(verbose bool) func(*gin.Context) {
	if verbose {
		return func(c *gin.Context) {
			bodyBytes, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				log.Println("Failed to read request body:", err)
			} else {
				if len(bodyBytes) > 0 {
					log.Printf("Body: %s\n\r", string(bodyBytes))
				}
			}

			_, err = c.Writer.Write(bodyBytes)
			if err != nil {
				log.Println("Failed to echo request body:", err)
			}
		}
	}

	return func(c *gin.Context) {
		err := c.Request.Write(c.Writer)
		if err != nil {
			log.Println("Failed to echo request body:", err)
		}

	}

}
