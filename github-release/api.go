package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type release struct {
	Name string `json:"name"`
	Tag  string `json:"tag_name"`
	Body string `json:"body"`

	URL    string  `json:"url"`
	Assets []asset `json:"assets"`

	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at"`
}

type asset struct {
	Name  string `json:"name"`
	Label string `json:"label"`

	ContentType string `json:"content_type"`
	Size        int    `json:"size"`

	URL         string `json:"url"`
	DownloadURL string `json:"browser_download_url"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func getLatestRelease(ctx context.Context, repo string) (*release, error) {
	target := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	// req.Header.Add("Authorization", "Bearer <YOUR-TOKEN>")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response %s", resp.Status)
	}

	var release release
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func getReleaseByTag(ctx context.Context, repo string, tag string) (*release, error) {
	target := fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/%s", repo, tag)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	// req.Header.Add("Authorization", "Bearer <YOUR-TOKEN>")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response %s", resp.Status)
	}

	var release release
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}
