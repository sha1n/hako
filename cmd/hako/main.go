package main

import (
	"fmt"
	"github.com/sha1n/hako/cmd/hako/startcmd"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

// ProgramName : passed from build environment
var ProgramName string

// Build : passed from build environment
var Build string

// Version : passed from build environment
var Version string

func init() {
	log.SetPrefix("[HAKO] ")
}

func main() {

	var rootCmd = &cobra.Command{
		Use: ProgramName,
		Version: fmt.Sprintf(`Version: %s
Build label: %s`, Version, Build),
	}
	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)
	rootCmd.AddCommand(createStartCommand())

	_ = rootCmd.Execute()
}

func createStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Long:  fmt.Sprintf(``),
		Short: fmt.Sprintf(``),
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

	config := startcmd.Config{
		ServerPort:     port,
		EchoPath:       normalizePath(additionalPath),
		Verbose:        verbose,
		VerboseHeaders: verboseHeaders,
		Delay:          delay,
	}

	startcmd.Start(config)
}

func normalizePath(path string) string {
	var normalizedPath = path

	normalizedPath = strings.TrimSpace(normalizedPath)

	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}

	return normalizedPath
}
