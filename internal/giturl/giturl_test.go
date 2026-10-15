package giturl

import "testing"

func TestGroupProjectPathMath(t *testing.T) {
	group, err := Parse("https://gitlab.com/my-org")
	if err != nil {
		t.Fatal(err)
	}
	proj, err := Parse("https://gitlab.com/my-org/group/project")
	if err != nil {
		t.Fatal(err)
	}
	if rel := RelativePath(group.FullPath, proj.FullPath); rel != "group/project" {
		t.Fatalf("relative path = %q, want group/project", rel)
	}
	if ssh, _ := proj.CloneURL("ssh"); ssh != "git@gitlab.com:my-org/group/project.git" {
		t.Fatalf("ssh clone url = %q", ssh)
	}
	if https, _ := proj.CloneURL("https"); https != "https://gitlab.com/my-org/group/project.git" {
		t.Fatalf("https clone url = %q", https)
	}
	if api := group.APIBaseURL(); api != "https://gitlab.com/api/v4" {
		t.Fatalf("api base url = %q", api)
	}
}
