package api

import (
	"context"
	"fmt"
)

func (api *API) GetLatestRelease(ctx context.Context, repo string) (*Release, error) {
	var release Release
	err := api.makeRequest(ctx, &release,
		withPath(fmt.Sprintf("/repos/%s/releases/latest", repo)),
	)
	if err != nil {
		return nil, err
	}
	return &release, err
}

func (api *API) GetReleaseByTag(ctx context.Context, repo string, tag string) (*Release, error) {
	var release Release
	err := api.makeRequest(ctx, &release,
		withPath(fmt.Sprintf("/repos/%s/releases/tags/%s", repo, tag)),
	)
	if err != nil {
		return nil, err
	}
	return &release, err
}

type Release struct {
	Name string `json:"name"`
	Tag  string `json:"tag_name"`
	Body string `json:"body"`

	URL    string  `json:"url"`
	Assets []Asset `json:"assets"`

	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at"`
}

type Asset struct {
	Name  string `json:"name"`
	Label string `json:"label"`

	ContentType string `json:"content_type"`
	Size        int    `json:"size"`

	URL         string `json:"url"`
	DownloadURL string `json:"browser_download_url"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
