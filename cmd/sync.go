package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"danny270793/gitmanager/internal/gitlabapi"
	"danny270793/gitmanager/internal/gitops"
	"danny270793/gitmanager/internal/giturl"
)

var (
	syncGroup       string
	syncRepo        string
	syncMethod      string
	syncDestination string
	syncToken       string
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Clone or update repositories from a GitLab group or a single repository",
	Long: `sync clones every repository under a GitLab group (recursing through
subgroups) into a local folder tree, or clones/updates a single repository.

Repositories that are already cloned at the destination are updated by
stashing any local changes, switching to main and pulling.`,
	Example: `  gitmanager sync --group=https://gitlab.com/sofiinc --method=ssh --destination=.
  gitmanager sync --repo=https://gitlab.com/sofiinc/money/funds-transfer --method=ssh --destination=.`,
	RunE: runSync,
}

func init() {
	syncCmd.Flags().StringVar(&syncGroup, "group", "", "GitLab group URL to sync recursively, e.g. https://gitlab.com/sofiinc")
	syncCmd.Flags().StringVar(&syncRepo, "repo", "", "Single GitLab repository URL to sync, e.g. https://gitlab.com/sofiinc/money/funds-transfer")
	syncCmd.Flags().StringVar(&syncMethod, "method", "ssh", "Clone method: ssh or https")
	syncCmd.Flags().StringVar(&syncDestination, "destination", ".", "Local destination folder")
	syncCmd.Flags().StringVar(&syncToken, "token", "", "GitLab personal access token (falls back to GITLAB_TOKEN env var; only needed for --group on private groups)")
}

func runSync(cmd *cobra.Command, args []string) error {
	if syncGroup == "" && syncRepo == "" {
		return fmt.Errorf("one of --group or --repo is required")
	}
	if syncGroup != "" && syncRepo != "" {
		return fmt.Errorf("--group and --repo are mutually exclusive")
	}
	if syncMethod != "ssh" && syncMethod != "https" {
		return fmt.Errorf("--method must be ssh or https")
	}

	token := syncToken
	if token == "" {
		token = os.Getenv("GITLAB_TOKEN")
	}

	if syncRepo != "" {
		return syncSingleRepo(syncRepo, syncMethod, syncDestination)
	}
	return syncGroupRepos(syncGroup, syncMethod, syncDestination, token)
}

type job struct {
	destPath string
	cloneURL string
}

func syncSingleRepo(repoURL, method, destination string) error {
	ref, err := giturl.Parse(repoURL)
	if err != nil {
		return err
	}
	cloneURL, err := ref.CloneURL(method)
	if err != nil {
		return err
	}

	relPath := ref.FullPath
	if syncGroup != "" {
		groupRef, err := giturl.Parse(syncGroup)
		if err != nil {
			return err
		}
		relPath = giturl.RelativePath(groupRef.FullPath, ref.FullPath)
	}

	j := job{
		destPath: filepath.Join(destination, filepath.FromSlash(relPath)),
		cloneURL: cloneURL,
	}
	return runJobs([]job{j})
}

func syncGroupRepos(groupURL, method, destination, token string) error {
	groupRef, err := giturl.Parse(groupURL)
	if err != nil {
		return err
	}

	fmt.Println("Discovering repositories in", groupURL, "...")
	projects, err := gitlabapi.ListGroupProjectsRecursive(groupRef.APIBaseURL(), token, groupRef.FullPath)
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		fmt.Println("No repositories found.")
		return nil
	}

	jobs := make([]job, 0, len(projects))
	for _, p := range projects {
		cloneURL := p.SSHURLToRepo
		if method == "https" {
			cloneURL = p.HTTPURLToRepo
		}
		relPath := giturl.RelativePath(groupRef.FullPath, p.PathWithNamespace)
		jobs = append(jobs, job{
			destPath: filepath.Join(destination, filepath.FromSlash(relPath)),
			cloneURL: cloneURL,
		})
	}
	return runJobs(jobs)
}

func runJobs(jobs []job) error {
	return runWithProgress(jobs, "syncing repositories",
		func(j job) string { return fmt.Sprintf("syncing %s", j.destPath) },
		func(j job) gitops.Result { return gitops.SyncRepo(j.cloneURL, j.destPath) },
	)
}
