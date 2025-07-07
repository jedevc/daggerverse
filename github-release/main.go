// Get details of GitHub releases and associated assets.

package main

import (
	"context"
	"fmt"
	"github/internal/dagger"
	"strings"
)

type GithubRelease struct{}

// Get the latest release for a repository
func (r *GithubRelease) GetLatest(
	ctx context.Context,

	// GitHub repository in the form of "owner/repo"
	repo string,
) (*Release, error) {
	tmp, err := getLatestRelease(ctx, repo)
	if err != nil {
		return nil, err
	}
	release := convertRelease(tmp)
	release.Repository = dag.Git(fmt.Sprintf("https://github.com/%s.git", repo))
	return &release, nil
}

// Get the specified release for a repository
func (r *GithubRelease) Get(
	ctx context.Context,

	// GitHub repository in the form of "owner/repo"
	repo string,
	// Tag name of the release
	tag string,
) (*Release, error) {
	tmp, err := getReleaseByTag(ctx, repo, tag)
	if err != nil {
		return nil, err
	}
	release := convertRelease(tmp)
	release.Repository = dag.Git(fmt.Sprintf("https://github.com/%s.git", repo))
	return &release, nil
}

// A GitHub release
type Release struct {
	// Git repository that the release is from
	Repository *dagger.GitRepository

	Name string
	Tag  string
	Body string

	URL string

	CreatedAt   string
	PublishedAt string

	Assets []Asset
}

// Git reference pointed to by the release tag
func (r *Release) Ref() *dagger.GitRef {
	return r.Repository.Tag(r.Tag)
}

// Return the shortest asset whose name contains all of the given keywords.
//
// This is neccessary because GitHub doesn't provide a way to filter assets by
// platform, or variant, or linux distro, etc.
func (r *Release) FindAsset(ctx context.Context, keywords []string) (*Asset, error) {
	var best *Asset

	for _, asset := range r.Assets {
		matches := true
		for _, keyword := range keywords {
			if !strings.Contains(asset.Name, keyword) {
				matches = false
				break
			}
		}
		if !matches {
			continue
		}

		if best == nil || len(asset.Name) < len(best.Name) {
			best = &asset
		}
	}

	if best == nil {
		return nil, fmt.Errorf("could not find asset matching keywords %s", keywords)
	}
	return best, nil
}

type Asset struct {
	Name  string
	Label string

	ContentType string
	Size        int

	URL         string
	DownloadURL string

	CreatedAt string
	UpdatedAt string
}

func (a *Asset) File() *dagger.File {
	return dag.HTTP(a.DownloadURL)
}

func convertRelease(tmp *release) Release {
	release := Release{
		Name: tmp.Name,
		Tag:  tmp.Tag,
		Body: tmp.Body,
		URL:  tmp.URL,

		CreatedAt:   tmp.CreatedAt,
		PublishedAt: tmp.PublishedAt,
	}
	for _, asset := range tmp.Assets {
		release.Assets = append(release.Assets, convertAsset(&asset))
	}
	return release
}

func convertAsset(tmp *asset) Asset {
	return Asset{
		Name:  tmp.Name,
		Label: tmp.Label,

		ContentType: tmp.ContentType,
		Size:        tmp.Size,

		URL:         tmp.URL,
		DownloadURL: tmp.DownloadURL,

		CreatedAt: tmp.CreatedAt,
		UpdatedAt: tmp.UpdatedAt,
	}
}
