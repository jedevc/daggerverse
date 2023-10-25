package main

import (
	"context"
	"fmt"
)

const (
	DEFAULT_HUGO_VERSION = "0.119.0"
	DEFAULT_DART_VERSION = "1.69.4"
)

type Hugo struct{}

func (h *Hugo) Build(ctx context.Context, d *Directory, hugoVersion Optional[string], dartSassVersion Optional[string], baseURL Optional[string], minify Optional[bool]) (*Directory, error) {
	srcPath := "/src"
	destPath := "/dest"
	cachePath := "/cache"
	args := []string{
		"hugo",
		"--gc",
		"--cacheDir", cachePath,
		"--destination", destPath,
	}
	url, isSet := baseURL.Get()
	if isSet {
		args = append(args, "--baseURL", url)
	}

	cache := dag.CacheVolume("hugo-cache")

	res := env(hugoVersion.GetOr(DEFAULT_HUGO_VERSION), dartVersion.GetOr(DEFUALT_DART_VERSION)).
		WithDirectory(srcPath, d).
		WithWorkdir(srcPath).
		WithMountedCache(cachePath, cache).
		WithExec(args).
		Directory(destPath)
	return res, nil
}

func env(hugoVersion string, dartVersion string) *Container {
	return dag.Container().
		From("debian:latest").
		WithExec([]string{"apt-get", "update", "-y"}).
		WithExec([]string{"apt-get", "install", "git", "-y"}).
		WithDirectory("/", hugo(opt.HugoVersion)).
		WithDirectory("/", sass(opt.DartSassVersion))
}

func hugo(version string) *Directory {
	arch := runtime.GOARCH

	download := fmt.Sprintf("https://github.com/gohugoio/hugo/releases/download/v%s/hugo_extended_%s_linux-%s.tar.gz", version, version, arch)
	file := dag.HTTP(download)
	filePath := "/mnt/hugo.tar.gz"

	hugo := dag.Container().From("alpine:latest").
		WithFile(filePath, file).
		WithExec([]string{"tar", "--extract", "--directory", "/mnt", "--file", filePath})
	return dag.Directory().WithFile("/usr/bin/hugo", hugo.File("/mnt/hugo"))
}

func sass(version string) *Directory {
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	}

	download := fmt.Sprintf("https://github.com/sass/dart-sass/releases/download/%s/dart-sass-%s-linux-%s.tar.gz", version, version, arch)
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
