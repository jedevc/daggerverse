// dockerc compiles container images into executable binary files.
//
// See https://github.com/NilsIrl/dockerc for more information.

package main

import (
	"context"
	"runtime"

	"dockerc/internal/dagger"
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
func (m *Dockerc) Compile(ctx context.Context, target *dagger.Container) (*dagger.File, error) {
	result := dag.Container().
		From("alpine").
		WithMountedFile("/usr/local/bin/dockerc", m.binary()).
		WithMountedFile("/target.tar.gz", target.AsTarball()).
		WithExec([]string{"/usr/local/bin/dockerc", "--image", "oci-archive:/target.tar.gz", "--output=/output.tar.gz"}).
		File("/output.tar.gz")
	return result, nil
}

func (m *Dockerc) binary() *dagger.File {
	var release *dagger.GithubReleaseRelease
	switch m.Version {
	case "", "latest":
		release = dag.GithubRelease().GetLatest("NilsIrl/dockerc")
	default:
		release = dag.GithubRelease().Get("NilsIrl/dockerc", m.Version)
	}

	arch := runtime.GOARCH
	if newArch, ok := releaseArch[arch]; ok {
		arch = newArch
	}
	asset := release.FindAsset([]string{"dockerc", arch}).File()
	f := dag.Directory().
		WithFile("/download", asset, dagger.DirectoryWithFileOpts{
			Permissions: 0o755,
		}).
		File("/download")
	return f
}

var releaseArch = map[string]string{
	"amd64": "x86-64",
	"arm64": "aarch64",
}
