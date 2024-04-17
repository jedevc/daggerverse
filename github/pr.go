package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github/api"
	"net/mail"
	"strings"
)

type PR struct {
	Repo *Repo

	URL    string
	Number int

	Title string
	Body  string
	State string
	Draft bool

	CreatedAt string
}

func (repo *Repo) MakePullRequest(
	ctx context.Context,

	dir *Directory,

	// +optional
	title string,
	// +optional
	body string,
	// +optional
	draft bool,

	// +optional
	// +default="[automated]"
	commitMessage string,
	// +optional
	signoff bool,
	// +optional
	author string,
	// +optional
	// +default="Automated <automated@example.com>"
	committer string,

	// +optional
	// +default="automated"
	branch string,
	// +optional
	deleteBranch bool,
	// +optional
	base string,
) (string, error) {
	committerAddr, err := mail.ParseAddress(committer)
	if err != nil {
		return "", err
	}
	git := dag.Container().
		From("alpine").
		WithExec([]string{"apk", "add", "git"}).
		WithWorkdir("/work").
		WithExec([]string{"git", "clone", fmt.Sprintf("https://github.com/%s/%s.git", repo.Owner, repo.Name), "."}).
		WithExec([]string{"git", "config", "user.name", committerAddr.Name}).
		WithExec([]string{"git", "config", "user.email", committerAddr.Address})

	if base != "" {
		git = git.WithExec([]string{"git", "checkout", base})
	} else {
		base, err = git.WithExec([]string{"git", "rev-parse", "--abbrev-ref", "HEAD"}).Stdout(ctx)
		if err != nil {
			return "", err
		}
		base = strings.TrimSpace(base)
	}
	git = git.WithExec([]string{"git", "checkout", "-b", branch})

	git = git.WithDirectory("/work", dir).
		WithExec([]string{"git", "add", "."})

	commitFlags := []string{"-m", commitMessage}
	if signoff {
		commitFlags = append(commitFlags, "-s")
	}
	git = git.WithExec(append([]string{"git", "commit"}, commitFlags...))

	if repo.Github.Token == nil {
		return "", fmt.Errorf("token is required")
	}

	githubTokenRaw, err := repo.Github.Token.Plaintext(ctx)
	if err != nil {
		return "", err
	}
	git = git.
		WithEnvVariable("GIT_CONFIG_COUNT", "1").
		WithEnvVariable("GIT_CONFIG_KEY_0", "http.https://github.com/.extraheader").
		WithSecretVariable("GIT_CONFIG_VALUE_0", dag.SetSecret("GITHUB_HEADER", fmt.Sprintf("AUTHORIZATION: Basic %s", base64.URLEncoding.EncodeToString([]byte("pat:"+githubTokenRaw)))))

	git = git.WithExec([]string{"git", "push", "origin", branch})
	_, err = git.Sync(ctx)
	if err != nil {
		return "", err
	}

	if title == "" {
		title = commitMessage
	}
	pr, err := repo.api().MakePR(ctx, repo.String(), map[string]string{
		"title": title,
		"body":  body,
		"head":  branch,
		"base":  base,
	})
	if err != nil {
		return "", err
	}

	res := convert[api.PR, PR](pr)
	res.Repo = repo
	return "", nil
}
