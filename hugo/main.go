package main

import (
	"context"
	"fmt"
	"runtime"
	"strings"
)

type Hugo struct{}

func (h *Hugo) Build(
	ctx context.Context,
	d *Directory,

	hugoVersion string,
	dartSassVersion string,

	baseURL string,
) (*Directory, error) {
	srcPath := "/src"
	destPath := "/dest"
	cachePath := "/cache"
	args := []string{
		"hugo",
		"--gc",
		"--cacheDir", cachePath,
		"--destination", destPath,
	}
	if baseURL != "" {
		args = append(args, "--baseURL", baseURL)
	}

	cache := dag.CacheVolume("hugo-cache")

	env, err := env(ctx, hugoVersion, dartSassVersion)
	if err != nil {
		return nil, err
	}
	res := env.
		WithDirectory(srcPath, d).
		WithWorkdir(srcPath).
		WithMountedCache(cachePath, cache).
		WithExec(args).
		Directory(destPath)
	return res, nil
}

func env(ctx context.Context, hugoVersion string, dartSassVersion string) (*Container, error) {
	hugo, err := hugo(ctx, hugoVersion)
	if err != nil {
		return nil, err
	}
	sass, err := sass(ctx, dartSassVersion)
	if err != nil {
		return nil, err
	}
	c := dag.Container().
		From("debian:latest").
		WithExec([]string{"apt-get", "update", "-y"}).
		WithExec([]string{"apt-get", "install", "git", "-y"}).
		WithDirectory("/", hugo).
		WithDirectory("/", sass)
	return c, nil
}

func hugo(ctx context.Context, version string) (*Directory, error) {
	var release *GithubRelease
	switch version {
	case "", "latest":
		release = dag.Github().GetLatestRelease("gohugoio/hugo")
	default:
		release = dag.Github().GetRelease("gohugoio/hugo", version)
	}

	assets, err := release.Assets(ctx)
	if err != nil {
		return nil, err
	}
	asset, err := getMatchingAsset(ctx, assets, []string{"hugo", "extended", runtime.GOOS, runtime.GOARCH, ".tar"})
	if err != nil {
		return nil, err
	}
	download, err := asset.DownloadURL(ctx)
	if err != nil {
		return nil, err
	}

	file := dag.HTTP(download)
	filePath := "/mnt/hugo.tar.gz"

	hugo := dag.Container().From("alpine:latest").
		WithFile(filePath, file).
		WithExec([]string{"tar", "--extract", "--directory", "/mnt", "--file", filePath})
	return dag.Directory().WithFile("/usr/bin/hugo", hugo.File("/mnt/hugo")), nil
}

func sass(ctx context.Context, version string) (*Directory, error) {
	var release *GithubRelease
	switch version {
	case "", "latest":
		release = dag.Github().GetLatestRelease("sass/dart-sass")
	default:
		release = dag.Github().GetRelease("sass/dart-sass", version)
	}

	assets, err := release.Assets(ctx)
	if err != nil {
		return nil, err
	}
	arch := runtime.GOARCH
	if newArch, ok := sassReleaseArch[arch]; ok {
		arch = newArch
	}

	asset, err := getMatchingAsset(ctx, assets, []string{"dart", "sass", runtime.GOOS, arch, ".tar"})
	if err != nil {
		return nil, err
	}
	download, err := asset.DownloadURL(ctx)
	if err != nil {
		return nil, err
	}

	file := dag.HTTP(download)
	filePath := "/mnt/dart-sass.tar.gz"

	sass := dag.Container().From("alpine:latest").
		WithFile(filePath, file).
		WithExec([]string{"tar", "--extract", "--directory", "/mnt", "--file", filePath, "--strip", "2"})
	script := `#!/bin/sh
exec "/usr/share/dart-sass/dart" "/usr/share/dart-sass/sass.snapshot" "$@"
`
	return dag.Directory().
		WithFile("/usr/share/dart-sass/dart", sass.File("/mnt/dart")).
		WithFile("/usr/share/dart-sass/sass.snapshot", sass.File("/mnt/sass.snapshot")).
		WithNewFile("/usr/bin/sass", script, DirectoryWithNewFileOpts{Permissions: 0o755}), nil
}

func getMatchingAsset(ctx context.Context, assets []GithubAsset, keywords []string) (GithubAsset, error) {
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

var sassReleaseArch = map[string]string{
	"amd64": "x64",
	"386":   "ia32",
	"arm":   "arm",
	"arm64": "arm64",
}
