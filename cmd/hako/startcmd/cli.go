package startcmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// CreateCommand creates a new cobra Command for the start CLI command.
func CreateCommand() *cobra.Command {
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
