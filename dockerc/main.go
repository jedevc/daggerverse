// dockerc compiles container images into executable binary files.
//
// See https://github.com/NilsIrl/dockerc for more information.

package main

import (
	"context"
	"fmt"
)

type Dockerc struct {
	Version string // +private
}

func New(
	// +optional
	// +default="latest"
	version string,
) *Dockerc {
	return &Dockerc{
		Version: version,
	}
}

// Compile an arbitrary target container into an executable file
func (m *Dockerc) Compile(ctx context.Context, target *Container) (*File, error) {
	binary, err := m.binary(ctx)
	if err != nil {
		return nil, err
	}

	result := dag.Container().
		From("alpine").
		WithMountedFile("/usr/local/bin/dockerc", binary).
		WithMountedFile("/target.tar.gz", target.AsTarball()).
		WithExec([]string{"/usr/local/bin/dockerc", "--image", "oci-archive:/target.tar.gz", "--output=/output.tar.gz"}).
		File("/output.tar.gz")
	return result, nil
}

func (m *Dockerc) binary(ctx context.Context) (*File, error) {
	var release *GithubRelease
	switch m.Version {
	case "", "latest":
		release = dag.Github().GetLatestRelease("NilsIrl/dockerc")
	default:
		release = dag.Github().GetRelease("NilsIrl/dockerc", m.Version)
	}

	assets, err := release.Assets(ctx)
	if err != nil {
		return nil, err
	}
	for _, asset := range assets {
		name, err := asset.Name(ctx)
		if err != nil {
			return nil, err
		}
		if name == "dockerc" {
			download, err := asset.DownloadURL(ctx)
			if err != nil {
				return nil, err
			}
			f := dag.Directory().
				WithFile("/download", dag.HTTP(download), DirectoryWithFileOpts{
					Permissions: 0o755,
				}).
				File("/download")
			return f, nil
		}
	}
	return nil, fmt.Errorf("no downloadable asset found")
}
