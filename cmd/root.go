package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitmanager",
	Short: "gitmanager clones and syncs GitLab repositories to a local folder tree",
	Long: `gitmanager clones and syncs GitLab repositories to a local folder tree.

It can walk an entire GitLab group recursively (including subgroups) and
clone or update every repository it finds, or sync a single repository.
Existing local clones are updated by stashing local changes, switching to
main and pulling.`,
	Version:       Version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(pullCmd)
}
