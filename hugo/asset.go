package main

import (
	"context"
	"fmt"
	"strings"
)

// getMatchingAsset returns the first asset whose name contains all of the
// given keywords.
//
// This is neccessary because GitHub doesn't provide a way to filter assets by
// platform, or variant, or linux distro, etc.
func getMatchingAsset(ctx context.Context, assets []GithubAsset, keywords ...string) (GithubAsset, error) {
	for _, asset := range assets {
		name, err := asset.Name(ctx)
		if err != nil {
			return GithubAsset{}, nil
		}

		matches := true
		for _, keyword := range keywords {
			if !strings.Contains(name, keyword) {
				matches = false
				break
			}
		}
		if matches {
			return asset, nil
		}
	}

	return GithubAsset{}, fmt.Errorf("could not find asset matching keywords %s", keywords)
}
