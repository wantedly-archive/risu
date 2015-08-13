package repository

import (
	"os"

	"github.com/libgit2/git2go"

	"github.com/wantedly/risu/schema"
)

// GitHubHost is base URL
// DefaultClonePath is default git clone path
const (
	GitHubHost       = "github.com/"
	DefaultClonePath = "/tmp/risu/repository/"
)

// GetRepository run "git clone <repository_URL>" and "git checkout rebision"
func GetRepository(build schema.Build, path string) error {
	if path == "" {
		path = DefaultClonePath
	}

	if _, err := os.Stat(path); err != nil {
		os.MkdirAll(path, 0755)
	}
	token := ""
	if os.Getenv("GITHUB_ACCESS_TOKEN") != "" {
		token = os.Getenv("GITHUB_ACCESS_TOKEN") + "@"
	}

	cloneURL := "https://" + token + GitHubHost + build.SourceRepo + ".git"
	clonePath := path

	_, err := git.Clone(cloneURL, clonePath, &git.CloneOptions{CheckoutBranch: build.SourceRevision})
	if err != nil {
		return err
	}
	return nil
}
