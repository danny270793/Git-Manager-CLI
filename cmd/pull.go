package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny270793/gitmanager/internal/gitops"
)

var pullDestination string

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Update every repository already cloned under a local folder",
	Long: `pull walks --destination recursively, finds every repository already
cloned there, and updates each one by stashing any local changes, switching
to main and pulling.

Unlike sync, pull never talks to GitLab: it does not discover new
repositories, it only updates the ones already present locally.`,
	Example: `  gitmanager pull --destination=.`,
	RunE:    runPull,
}

func init() {
	pullCmd.Flags().StringVar(&pullDestination, "destination", ".", "Local folder to walk for existing repositories")
}

func runPull(cmd *cobra.Command, args []string) error {
	repos, err := gitops.FindRepos(pullDestination)
	if err != nil {
		return err
	}
	if len(repos) == 0 {
		fmt.Println("No local repositories found under", pullDestination)
		return nil
	}

	return runWithProgress(repos, "pulling repositories",
		func(path string) string { return fmt.Sprintf("pulling %s", path) },
		gitops.UpdateRepo,
	)
}
