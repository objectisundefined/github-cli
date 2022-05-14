package cmd

import (
	"context"
	"os/user"
	"path"

	"github.com/spf13/cobra"

	"github.com/objectisundefined/github-cli/internal/config"
)

// shared context and config
var (
	ctx context.Context
	cfg *config.Config
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "github-cli",
		Short: "Github CLI",
	}

	rootCmd.PersistentPreRun = preRunRootHandler

	rootCmd.PersistentFlags().String("config", "", "Config File, default ~/.github-cli/config.toml")
	rootCmd.PersistentFlags().String("token", "", "Github Token")

	rootCmd.AddCommand(
		NewPullsCommand(),
		NewPullCommand(),
		NewIssuesCommand(),
		NewIssueCommand(),
		NewTrendingCommand(),
		NewEventsCommand(),
	)

	cobra.EnablePrefixMatching = true

	return rootCmd
}

func preRunRootHandler(cmd *cobra.Command, args []string) {
	usr, err := user.Current()

	if err != nil {
		panic(err)
	}

	file, _ := cmd.Flags().GetString("config")
	token, _ := cmd.Flags().GetString("token")

	if len(file) == 0 {
		file = path.Join(usr.HomeDir, ".github-cli/config.toml")

	}

	cfg_, err := config.NewConfigFromFile(file)

	if err != nil {
		panic(err)
	}

	if len(token) > 0 {
		cfg_.Token = token
	}

	ctx = context.Background()
	cfg = cfg_
}
