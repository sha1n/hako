package internal

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	gommonsos "github.com/sha1n/gommons/pkg/os"
	"github.com/spf13/cobra"
)

// CreateStartCommand creates a new cobra Command for the start CLI command.
func CreateStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Long:  fmt.Sprintf(`starts the server`),
		Short: fmt.Sprintf(`starts the server`),
		Run:   doStart,
	}
	cmd.Flags().IntP("port", "p", 8080, `Port to listen on. Default is 8080`)
	cmd.Flags().StringP("path", "", "", `Path of incoming requests`)
	cmd.Flags().Int32P("delay", "d", 0, `The minimum delay of each response in milliseconds`)
	cmd.Flags().BoolP("verbose", "v", false, `Prints the body of every incoming request`)
	cmd.Flags().BoolP("verbose-headers", "", false, `Prints the headers of every incoming request`)

	return cmd
}

func doStart(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetInt("port")
	delay, _ := cmd.Flags().GetInt32("delay")
	additionalPath, _ := cmd.Flags().GetString("path")
	verbose, _ := cmd.Flags().GetBool("verbose")
	verboseHeaders, _ := cmd.Flags().GetBool("verbose-headers")

	config := Config{
		ServerPort:     port,
		EchoPath:       normalizePath(additionalPath),
		Verbose:        verbose,
		VerboseHeaders: verboseHeaders,
		Delay:          delay,
	}

	Start(config)
}

func normalizePath(path string) string {
	var normalizedPath = path

	normalizedPath = strings.TrimSpace(normalizedPath)

	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}

	return normalizedPath
}

// Config start command configuration
type Config struct {
	ServerPort     int
	EchoPath       string
	Verbose        bool
	VerboseHeaders bool
	Delay          int32
}

// StartAsync starts an echo server in the background and returns immediately.
func StartAsync(config Config) func() {
	server := createHTTPServer(config)
	server.StartAsync()

	return server.StopAsync
}

// Start starts an echo server in the background and awaits shutdown signal.
func Start(config Config) {
	StartAsync(config)

	awaitShutdownSig()
}

func awaitShutdownSig() {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Waiting for shutdown signal...")

	<-quitChannel
}

func createHTTPServer(config Config) Server {
	router := createGinEngine(config)

	server := NewServer(config.ServerPort, router)

	stopServerAsync := func() {
		server.StopAsync()
	}

	log.Println("Registering signal listeners for graceful HTTP server shutdown..")
	gommonsos.RegisterShutdownHook(gommonsos.NewSignalHook(syscall.SIGTERM, stopServerAsync))
	gommonsos.RegisterShutdownHook(gommonsos.NewSignalHook(syscall.SIGKILL, stopServerAsync))

	return server
}

func createGinEngine(config Config) *gin.Engine {
	var router *gin.Engine
	if config.Verbose {
		router = NewDefaultEngine()
	} else {
		router = NewSilentEngine()
	}
	registerHandlers(router, "/echo", echoHandlerWith(config.Verbose, config.VerboseHeaders, config.Delay))

	if config.EchoPath != "" {
		registerHandlers(router, config.EchoPath, echoHandlerWith(config.Verbose, config.VerboseHeaders, config.Delay))
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
