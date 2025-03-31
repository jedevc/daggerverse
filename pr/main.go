package main

import (
	"context"
	"dagger/pr/internal/dagger"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"strconv"

	"github.com/google/go-github/v70/github"
)

type Pr struct{}

type PullRequest struct {
	Owner  string
	Repo   string
	Number int
}

func (m *Pr) Open(ctx context.Context, ref string) (*PullRequest, error) {
	// parse the ref, getting the username/repo and the PR number, following a lax regexp pattern
	// e.g. "github.com/username/repo/pull/123" or "username/repo#123"
	re := regexp.MustCompile(`^(?:[^/]+//)?(?:github\.com/)?([^/]+)/([^/]+)(?:/pull/|#)(\d+)(?:/.*)?`)
	match := re.FindStringSubmatch(ref)
	if match == nil {
		return nil, fmt.Errorf("invalid ref format: %s", ref)
	}

	// Extract the username, repo, and PR number from the match
	username := match[1]
	repo := match[2]
	prStr := match[3]
	prNumber, err := strconv.Atoi(prStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PR number: %s", prStr)
	}

	pr := &PullRequest{
		Owner:  username,
		Repo:   repo,
		Number: prNumber,
	}
	_, err = pr.AsHeadRef().Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit for head ref: %w", err)
	}

	return pr, nil
}

func (pr *PullRequest) AsGit() *dagger.GitRepository {
	return dag.Git(fmt.Sprintf("https://github.com/%s/%s", pr.Owner, pr.Repo))
}

func (pr *PullRequest) AsHeadRef() *dagger.GitRef {
	return pr.AsGit().Ref(pr.headRefStr())
}
func (pr *PullRequest) headRefStr() string {
	return fmt.Sprintf("refs/pull/%d/head", pr.Number)
}

func (pr *PullRequest) AsMergeRef() *dagger.GitRef {
	return pr.AsGit().Ref(pr.mergeRefStr())
}
func (pr *PullRequest) mergeRefStr() string {
	return fmt.Sprintf("refs/pull/%d/merge", pr.Number)
}

func (pr *PullRequest) AsBaseRef(ctx context.Context) (*dagger.GitRef, error) {
	baseRef, err := pr.baseRefStr(ctx)
	if err != nil {
		return nil, err
	}
	return pr.AsGit().Ref(baseRef), nil
}
func (pr *PullRequest) baseRefStr(ctx context.Context) (string, error) {
	res, err := pr.api(ctx)
	if err != nil {
		return "", err
	}
	return *res.Base.Ref, nil
}

func (pr *PullRequest) Diff(ctx context.Context) (*dagger.File, error) {
	res, err := pr.api(ctx)
	if err != nil {
		return nil, err
	}
	if res.DiffURL == nil {
		return nil, fmt.Errorf("no diff URL found")
	}
	return pr.readGitFile(ctx, *res.DiffURL)
}

func (pr *PullRequest) Patch(ctx context.Context) (*dagger.File, error) {
	res, err := pr.api(ctx)
	if err != nil {
		return nil, err
	}
	if res.PatchURL == nil {
		return nil, fmt.Errorf("no patch URL found")
	}
	return pr.readGitFile(ctx, *res.PatchURL)
}

func (pr *PullRequest) api(ctx context.Context) (*github.PullRequest, error) {
	client := github.NewClient(nil)
	res, _, err := client.PullRequests.Get(ctx, pr.Owner, pr.Repo, pr.Number)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (pr *PullRequest) readGitFile(ctx context.Context, url string) (*dagger.File, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.text+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get file: %s", resp.Status)
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	name := path.Base(url)
	return dag.Directory().WithNewFile(name, string(content)).File(name), nil
}
