package main

import (
	"fmt"
	"log"

	gommons "github.com/sha1n/gommons/pkg/cmd"
	"github.com/sha1n/hako/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ProgramName : passed from build environment
var ProgramName string

// Build : passed from build environment
var Build string

// Version : passed from build environment
var Version string

// GitHubOwner : repository owner on github
const GitHubOwner = "sha1n"

// GitHubRepoName : repository name on github
const GitHubRepoName = "hako"

func init() {
	log.SetPrefix("[HAKO] ")
	logrus.StandardLogger().SetFormatter(
		&logrus.TextFormatter{
			DisableTimestamp: true,
		},
	)
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
	rootCmd.AddCommand(gommons.CreateUpdateCommand(GitHubOwner, GitHubRepoName, Version, ProgramName))

	_ = rootCmd.Execute()
}
