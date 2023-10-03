package main

import (
	"context"
	"fmt"
)

type GithubRelease struct {
	Repository *GitRepository

	// TODO: what is doing on with the need for this json thing, why is it needed?
	// Error: make request: input:1: github.latestRelease.name Cannot return null for non-nullable field GithubRelease.name.

	Name string `json:"name"`
	Tag  string `json:"tag"`
	Body string `json:"body"`

	URL string `json:"url"`

	CreatedAt   string `json:"createdAt"`
	PublishedAt string `json:"publishedAt"`

	Assets []GithubAsset `json:"assets"`
}

func (r *GithubRelease) Ref() *GitRef {
	return r.Repository.Tag(r.Tag)
}

type GithubAsset struct {
	Name  string `json:"name"`
	Label string `json:"label"`

	ContentType string `json:"contentType"`
	Size        int    `json:"size"`

	URL         string `json:"url"`
	DownloadURL string `json:"downloadUrl"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (a *GithubAsset) Void() {}

type Github struct{}

// TODO: use GitRepository and parse the git url from there

func (r *Github) GetLatestRelease(ctx context.Context, repo string) (*GithubRelease, error) {
	tmp, err := getLatestRelease(ctx, repo)
	if err != nil {
		return nil, err
	}
	release := convertRelease(tmp)
	release.Repository = dag.Git(fmt.Sprintf("https://github.com/%s.git", repo))

	return release, nil
}

func (r *Github) GetRelease(ctx context.Context, repo string, tag string) (*GithubRelease, error) {
	tmp, err := getReleaseByTag(ctx, repo, tag)
	if err != nil {
		return nil, err
	}
	release := convertRelease(tmp)
	release.Repository = dag.Git(fmt.Sprintf("https://github.com/%s.git", repo))

	return release, nil
}

// TODO: we can remove tmp once the json problem is resolved
func convertRelease(tmp *release) *GithubRelease {
	release := GithubRelease{
		Name: tmp.Name,
		Tag:  tmp.Tag,
		Body: tmp.Body,
		URL:  tmp.URL,

		CreatedAt:   tmp.CreatedAt,
		PublishedAt: tmp.PublishedAt,
	}
	for _, tmp := range tmp.Assets {
		release.Assets = append(release.Assets, GithubAsset{
			Name:  tmp.Name,
			Label: tmp.Label,

			ContentType: tmp.ContentType,
			Size:        tmp.Size,

			URL:         tmp.URL,
			DownloadURL: tmp.DownloadURL,

			CreatedAt: tmp.CreatedAt,
			UpdatedAt: tmp.UpdatedAt,
		})
	}
	return &release
}
