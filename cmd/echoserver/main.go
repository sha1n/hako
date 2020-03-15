package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// ProgramName : passed from build environment
var ProgramName string

// Build : passed from build environment
var Build string

// Version : passed from build environment
var Version string

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
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetInt("port")
			additionalPath, _ := cmd.Flags().GetString("path")
			verbose, _ := cmd.Flags().GetBool("verbose")

			Start(port, normalizePath(additionalPath), verbose)
		},
	}
	cmd.Flags().IntP("port", "p", 8080, `Port to listen on. Default is 8080`)
	cmd.Flags().StringP("path", "", "", `Path of incoming requests`)
	cmd.Flags().BoolP("verbose", "", false, `Prints the body of every incoming request`)

	return cmd
}

func normalizePath(path string) string {
	var normalizedPath = path

	normalizedPath = strings.TrimSpace(normalizedPath)

	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}

	return normalizedPath
}
