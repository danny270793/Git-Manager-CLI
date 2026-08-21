package cmd

import (
	"fmt"
	"os"

	"github.com/schollz/progressbar/v3"

	"danny270793/gitmanager/internal/gitops"
)

// runWithProgress runs do for every item, rendering a progress bar labeled
// with the given description, then prints a summary and any failures.
func runWithProgress[T any](items []T, description string, describe func(T) string, do func(T) gitops.Result) error {
	bar := progressbar.NewOptions(len(items),
		progressbar.OptionSetDescription(description),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(30),
		progressbar.OptionClearOnFinish(),
	)

	var failures []gitops.Result
	for _, item := range items {
		bar.Describe(describe(item))
		res := do(item)
		if res.Err != nil {
			failures = append(failures, res)
		}
		bar.Add(1)
	}
	fmt.Println()

	for _, r := range failures {
		fmt.Fprintf(os.Stderr, "FAILED %s (%s): %v\n", r.Destination, r.Action, r.Err)
	}
	fmt.Printf("Done: %d succeeded, %d failed\n", len(items)-len(failures), len(failures))
	if len(failures) > 0 {
		return fmt.Errorf("%d repositories failed", len(failures))
	}
	return nil
}
