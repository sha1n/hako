package main

import (
	"fmt"
	"github.com/sha1n/hako/internal"
	"github.com/spf13/cobra"
	"log"
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

	rootCmd.AddCommand(internal.CreateStartCommand())
	rootCmd.AddCommand(internal.CreateUpdateCommand(Version, ProgramName))

	_ = rootCmd.Execute()
}
