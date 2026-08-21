// Package gitlabapi lists all projects within a GitLab group, recursing
// through subgroups.
package gitlabapi

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Project is the subset of GitLab project fields gitmanager needs.
type Project struct {
	PathWithNamespace string
	SSHURLToRepo      string
	HTTPURLToRepo     string
}

// ListGroupProjectsRecursive returns every project in the given group and
// all of its subgroups. baseURL is the GitLab API base, e.g.
// "https://gitlab.com/api/v4". token may be empty for public groups.
func ListGroupProjectsRecursive(baseURL, token, groupFullPath string) ([]Project, error) {
	client, err := newClient(baseURL, token)
	if err != nil {
		return nil, err
	}

	includeSubgroups := true
	archived := false
	opts := &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
		IncludeSubGroups: &includeSubgroups,
		Archived:         &archived,
	}

	var projects []Project
	for {
		page, resp, err := client.Groups.ListGroupProjects(groupFullPath, opts)
		if err != nil {
			return nil, fmt.Errorf("list projects for group %q: %w", groupFullPath, err)
		}
		for _, p := range page {
			projects = append(projects, Project{
				PathWithNamespace: p.PathWithNamespace,
				SSHURLToRepo:      p.SSHURLToRepo,
				HTTPURLToRepo:     p.HTTPURLToRepo,
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return projects, nil
}

func newClient(baseURL, token string) (*gitlab.Client, error) {
	if token == "" {
		return gitlab.NewClient("", gitlab.WithBaseURL(baseURL))
	}
	return gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
}
