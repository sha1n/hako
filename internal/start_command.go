package internal

import (
	"fmt"
	"io"
	"log"
	"log/slog"
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
	cmd.Flags().BoolP("json", "", false, `Use JSON logging format`)

	return cmd
}

func doStart(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetInt("port")
	delay, _ := cmd.Flags().GetInt32("delay")
	additionalPath, _ := cmd.Flags().GetString("path")
	verbose, _ := cmd.Flags().GetBool("verbose")
	verboseHeaders, _ := cmd.Flags().GetBool("verbose-headers")
	jsonLog, _ := cmd.Flags().GetBool("json")

	config := Config{
		ServerPort:     port,
		EchoPath:       normalizePath(additionalPath),
		Verbose:        verbose,
		VerboseHeaders: verboseHeaders,
		Delay:          delay,
		JSONLog:        jsonLog,
	}

	Start(config)
}

// NewConfigFromArgs creates a Config object from parsed arguments (helper for testing)
func NewConfigFromArgs(port int, delay int32, path string, verbose, verboseHeaders, jsonLog bool) Config {
	return Config{
		ServerPort:     port,
		EchoPath:       normalizePath(path),
		Verbose:        verbose,
		VerboseHeaders: verboseHeaders,
		Delay:          delay,
		JSONLog:        jsonLog,
	}
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
	JSONLog        bool
}

// StartAsync starts an echo server in the background and returns immediately.
func StartAsync(config Config) func() {
	server := createHTTPServer(config)
	server.StartAsync()

	return server.StopAsync
}

// Start starts an echo server in the background and awaits shutdown signal.
func Start(config Config) {
	configureLogging(config)
	StartAsync(config)

	awaitShutdownSig()
}

func configureLogging(config Config) {
	configureLoggingWithOutput(config, os.Stderr)
}

func configureLoggingWithOutput(config Config, w io.Writer) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if config.JSONLog {
		handler = slog.NewJSONHandler(w, opts)
		// Redirect standard log library to write to slog
		log.SetFlags(0)
		log.SetOutput(newSlogWriter(handler))
	} else {
		handler = slog.NewTextHandler(w, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}
				return a
			},
		})
	}

	slog.SetDefault(slog.New(handler))
}

type slogWriter struct {
	logger *slog.Logger
}

func newSlogWriter(handler slog.Handler) *slogWriter {
	return &slogWriter{
		logger: slog.New(handler),
	}
}

func (w *slogWriter) Write(p []byte) (n int, err error) {
	// Remove trailing newline if present, as slog adds its own
	msg := string(p)
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}
	w.logger.Info(msg)
	return len(p), nil
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
