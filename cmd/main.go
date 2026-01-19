package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"

	gommons "github.com/sha1n/gommons/pkg/cmd"
	"github.com/sha1n/hako/internal"
	"github.com/spf13/cobra"
)

var (
	// ProgramName : passed from build environment
	ProgramName string

	// Build : passed from build environment
	Build string

	// Version : passed from build environment
	Version string

	// DisableSelfUpdate : passed from build environment to specify that self-update should be disabled
	DisableSelfUpdate string = ""
)

// GitHubOwner : repository owner on github
const GitHubOwner = "sha1n"

// GitHubRepoName : repository name on github
const GitHubRepoName = "hako"

func init() {
	log.SetPrefix("[HAKO] ")
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})))
}

func main() {

	var rootCmd = &cobra.Command{
		Use: ProgramName,
		Version: fmt.Sprintf(`Version: %s
Build label: %s`, Version, Build),
	}
	rootCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)

	rootCmd.AddCommand(internal.CreateStartCommand())
	rootCmd.AddCommand(gommons.CreateShellCompletionScriptGenCommand())
	if enableSelfUpdate() {
		rootCmd.AddCommand(gommons.CreateUpdateCommand(GitHubOwner, GitHubRepoName, Version, ProgramName))
	}

	_ = rootCmd.Execute()
}

func enableSelfUpdate() bool {
	return DisableSelfUpdate != "true" && runtime.GOOS != "windows"
}
