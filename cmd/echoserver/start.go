package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sha1n/echo-server/cmd/echoserver/http"
	"github.com/sha1n/echo-server/cmd/echoserver/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Start(port int, additionalEchoPath string) {
	server := createHttpServer(port, additionalEchoPath)
	server.StartAsync()

	awaitShutdownSig()
}

func awaitShutdownSig() {
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Waiting for shutdown signal...")

	<-quitChannel
}

func createHttpServer(port int, additionalEchoPath string) http.Server {
	serverBuilder := http.
		NewServer(port).
		WithGetHandler("/echo", handleEcho).
		WithPostHandler("/echo", handleEcho)

	if additionalEchoPath != "" {
		serverBuilder.
			WithGetHandler(additionalEchoPath, handleEcho).
			WithPostHandler(additionalEchoPath, handleEcho).
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

func handleEcho(c *gin.Context) {
	err := c.Request.Write(c.Writer)
	if err != nil {
		log.Println("Failed to echo request:", err)
	}

}
