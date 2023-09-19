package main

import (
	"context"
	"fmt"
)

type HugoOpt struct {
	HugoVersion     string
	DartSassVersion string

	BaseURL string
	Minify  bool
}

type Hugo struct{}

func (h *Hugo) Build(ctx context.Context, d *Directory, opt *HugoOpt) (*Directory, error) {
	return d.HugoBuild(ctx, opt)
}

func (d *Directory) HugoBuild(ctx context.Context, opt *HugoOpt) (*Directory, error) {
	srcPath := "/src"
	destPath := "/dest"
	cachePath := "/cache"
	args := []string{
		"hugo",
		"--gc",
		"--cacheDir", cachePath,
		"--destination", destPath,
	}
	if opt.BaseURL != "" {
		args = append(args, "--baseURL", opt.BaseURL)
	}

	cache := dag.CacheVolume("hugo-cache")

	res := env(opt).
		WithDirectory(srcPath, d).
		WithWorkdir(srcPath).
		WithMountedCache(cachePath, cache).
		WithExec(args).
		Directory(destPath)
	return res, nil
}

func env(opt *HugoOpt) *Container {
	return dag.Container().
		From("debian:latest").
		WithExec([]string{"apt-get", "update", "-y"}).
		WithExec([]string{"apt-get", "install", "git", "-y"}).
		WithDirectory("/", hugo(opt.HugoVersion)).
		WithDirectory("/", sass(opt.DartSassVersion))
}

func hugo(version string) *Directory {
	download := fmt.Sprintf("https://github.com/gohugoio/hugo/releases/download/v%s/hugo_extended_%s_linux-amd64.tar.gz", version, version)
	file := dag.HTTP(download)
	filePath := "/mnt/hugo.tar.gz"

	hugo := dag.Container().From("alpine:latest").
		WithFile(filePath, file).
		WithExec([]string{"tar", "--extract", "--directory", "/mnt", "--file", filePath})
	return dag.Directory().WithFile("/usr/bin/hugo", hugo.File("/mnt/hugo"))
}

func sass(version string) *Directory {
	download := fmt.Sprintf("https://github.com/sass/dart-sass/releases/download/%s/dart-sass-%s-linux-x64.tar.gz", version, version)
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
		WithNewFile("/usr/bin/sass", script, DirectoryWithNewFileOpts{Permissions: 0o755})
}
