// Package giturl parses GitLab HTTP(S) URLs into host and namespace path
// components, and builds clone URLs for a given method.
package giturl

import (
	"fmt"
	"net/url"
	"strings"
)

// Ref is a parsed GitLab group/project URL.
type Ref struct {
	Scheme   string // http or https
	Host     string // e.g. gitlab.com
	FullPath string // e.g. sofiinc/money/funds-transfer (no leading/trailing slash, no .git)
}

// Parse parses a GitLab group or project URL such as
// "https://gitlab.com/sofiinc/money/funds-transfer".
func Parse(raw string) (Ref, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return Ref{}, fmt.Errorf("invalid URL %q: %w", raw, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return Ref{}, fmt.Errorf("invalid URL %q: expected http:// or https:// scheme", raw)
	}
	if u.Host == "" {
		return Ref{}, fmt.Errorf("invalid URL %q: missing host", raw)
	}
	full := strings.Trim(u.Path, "/")
	full = strings.TrimSuffix(full, ".git")
	if full == "" {
		return Ref{}, fmt.Errorf("invalid URL %q: missing group/project path", raw)
	}
	return Ref{Scheme: u.Scheme, Host: u.Host, FullPath: full}, nil
}

// CloneURL builds a clone URL for the given method ("ssh" or "https").
func (r Ref) CloneURL(method string) (string, error) {
	switch method {
	case "ssh":
		return fmt.Sprintf("git@%s:%s.git", r.Host, r.FullPath), nil
	case "https", "http":
		return fmt.Sprintf("%s://%s/%s.git", r.Scheme, r.Host, r.FullPath), nil
	default:
		return "", fmt.Errorf("unsupported method %q: expected ssh or https", method)
	}
}

// APIBaseURL returns the GitLab API v4 base URL for this host.
func (r Ref) APIBaseURL() string {
	return fmt.Sprintf("%s://%s/api/v4", r.Scheme, r.Host)
}

// RelativePath returns fullPath with the given group's full path prefix
// stripped, so a project can be placed relative to the group's own
// destination folder. If fullPath does not start with groupFullPath, the
// original fullPath is returned unchanged.
func RelativePath(groupFullPath, fullPath string) string {
	prefix := strings.Trim(groupFullPath, "/") + "/"
	if strings.HasPrefix(fullPath, prefix) {
		return strings.TrimPrefix(fullPath, prefix)
	}
	return fullPath
}
