package giturl

import "testing"

func TestGroupProjectPathMath(t *testing.T) {
	group, err := Parse("https://gitlab.com/sofiinc")
	if err != nil {
		t.Fatal(err)
	}
	proj, err := Parse("https://gitlab.com/sofiinc/money/funds-transfer")
	if err != nil {
		t.Fatal(err)
	}
	if rel := RelativePath(group.FullPath, proj.FullPath); rel != "money/funds-transfer" {
		t.Fatalf("relative path = %q, want money/funds-transfer", rel)
	}
	if ssh, _ := proj.CloneURL("ssh"); ssh != "git@gitlab.com:sofiinc/money/funds-transfer.git" {
		t.Fatalf("ssh clone url = %q", ssh)
	}
	if https, _ := proj.CloneURL("https"); https != "https://gitlab.com/sofiinc/money/funds-transfer.git" {
		t.Fatalf("https clone url = %q", https)
	}
	if api := group.APIBaseURL(); api != "https://gitlab.com/api/v4" {
		t.Fatalf("api base url = %q", api)
	}
}
