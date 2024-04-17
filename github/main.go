// Get details of GitHub releases and associated assets.

package main

import (
	"fmt"
	"strings"

	"github/api"
	"github/internal/dagger"
)

type Github struct {
	Token *Secret // +private
}

func (gh *Github) WithTokenAuth(token *Secret) *Github {
	ghClone := *gh
	ghClone.Token = token
	return &ghClone
}

func (gh *Github) Repo(repo string) (*Repo, error) {
	owner, name, ok := strings.Cut(repo, "/")
	if !ok {
		return nil, fmt.Errorf("%s is not a valid repo", repo)
	}
	return &Repo{
		Github: gh,
		Owner:  owner,
		Name:   name,
	}, nil
}

func (gh *Github) api() *api.API {
	return &api.API{Token: gh.Token.Plaintext}
}

type Repo struct {
	*Github // +private

	Owner string
	Name  string
}

func (r *Repo) String() string {
	return r.Owner + "/" + r.Name
}

func (r *Repo) Git() *dagger.GitRepository {
	return dag.Git(fmt.Sprintf("https://github.com/%s/%s.git", r.Owner, r.Name))
}
