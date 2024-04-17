package main

import (
	"context"
	"github/api"
)

// A GitHub release
type Release struct {
	Repo *Repo

	Name string
	Tag  string
	Body string

	URL string

	CreatedAt   string
	PublishedAt string

	Assets []Asset
}

// Git reference pointed to by the release tag
func (r *Release) Git() *GitRef {
	return r.Repo.Git().Tag(r.Tag)
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

// Get the latest release for a repository
func (r *Repo) GetLatestRelease(ctx context.Context) (*Release, error) {
	release, err := r.api().GetLatestRelease(ctx, r.String())
	if err != nil {
		return nil, err
	}
	res := convert[api.Release, Release](release)
	res.Repo = r
	return res, nil
}

// Get the specified release for a repository
func (r *Repo) GetRelease(
	ctx context.Context,

	// Tag name of the release
	tag string,
) (*Release, error) {
	release, err := r.api().GetReleaseByTag(ctx, r.String(), tag)
	if err != nil {
		return nil, err
	}
	res := convert[api.Release, Release](release)
	res.Repo = r
	return res, nil
}
