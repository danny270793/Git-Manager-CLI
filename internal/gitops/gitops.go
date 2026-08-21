// Package gitops shells out to the git CLI to clone or update repositories.
package gitops

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

// Action describes what SyncRepo ended up doing.
type Action string

const (
	ActionCloned  Action = "cloned"
	ActionUpdated Action = "updated"
)

// Result is the outcome of syncing a single repository.
type Result struct {
	Destination string
	Action      Action
	Err         error
}

// SyncRepo clones cloneURL into destPath if it doesn't already exist there,
// or, if it does, stashes any local changes, switches to main and pulls.
func SyncRepo(cloneURL, destPath string) Result {
	if isGitRepo(destPath) {
		return UpdateRepo(destPath)
	}

	if err := cloneRepo(cloneURL, destPath); err != nil {
		return Result{Destination: destPath, Action: ActionCloned, Err: err}
	}
	return Result{Destination: destPath, Action: ActionCloned}
}

// UpdateRepo stashes any local changes, switches to main and pulls the
// repository already checked out at destPath.
func UpdateRepo(destPath string) Result {
	if err := updateRepo(destPath); err != nil {
		return Result{Destination: destPath, Action: ActionUpdated, Err: err}
	}
	return Result{Destination: destPath, Action: ActionUpdated}
}

// FindRepos walks root recursively and returns the paths of every directory
// that is the top of a git working tree (i.e. contains a .git entry). It
// does not descend into a repository once found, so nested checkouts are
// treated as a single repo.
func FindRepos(root string) ([]string, error) {
	var repos []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if d.Name() == ".git" {
			return filepath.SkipDir
		}
		if isGitRepo(path) {
			repos = append(repos, path)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk %q: %w", root, err)
	}
	return repos, nil
}

func isGitRepo(destPath string) bool {
	info, err := os.Stat(filepath.Join(destPath, ".git"))
	return err == nil && (info.IsDir() || info.Mode().IsRegular())
}

func cloneRepo(cloneURL, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("create parent dir: %w", err)
	}
	return run("", "git", "clone", cloneURL, destPath)
}

func updateRepo(destPath string) error {
	if err := run(destPath, "git", "stash", "--include-untracked"); err != nil {
		return fmt.Errorf("stash: %w", err)
	}
	if err := run(destPath, "git", "checkout", "main"); err != nil {
		return fmt.Errorf("checkout main: %w", err)
	}
	if err := run(destPath, "git", "pull"); err != nil {
		return fmt.Errorf("pull: %w", err)
	}
	return nil
}

func run(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %v: %w\n%s", name, args, err, out.String())
	}
	return nil
}
