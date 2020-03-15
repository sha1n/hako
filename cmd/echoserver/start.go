package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sha1n/echo-server/cmd/echoserver/http"
	"github.com/sha1n/echo-server/cmd/echoserver/utils"
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
	serverBuilder := http.
		NewServer(port).
		WithGetHandler("/echo", echoHandlerWith(verbose)).
		WithPostHandler("/echo", echoHandlerWith(verbose))

	if additionalEchoPath != "" {
		serverBuilder.
			WithGetHandler(additionalEchoPath, echoHandlerWith(verbose)).
			WithPostHandler(additionalEchoPath, echoHandlerWith(verbose)).
			Build()
	}

	server := serverBuilder.Build()

	stopServerAsync := func() {
		server.StopAsync()
	}

	log.Println("Registering signal listeners for graceful HTTP server shutdown..")
	utils.RegisterShutdownHook(utils.NewSignalHook(syscall.SIGTERM, stopServerAsync))
	utils.RegisterShutdownHook(utils.NewSignalHook(syscall.SIGKILL, stopServerAsync))

	return server
}

func echoHandlerWith(verbose bool) func(*gin.Context) {
	if verbose {
		return func(c *gin.Context) {
			bodyBytes, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				log.Println("Failed to read request body:", err)
			} else {
				fmt.Printf("Received: %s\n\r", string(bodyBytes))
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
