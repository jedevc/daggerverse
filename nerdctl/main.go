package main

import (
	"context"
	"dagger/nerdctl/internal/dagger"
)

type Nerdctl struct {
	Version string

	Container *dagger.Container
}

func New(
	// +optional
	version string,
) *Nerdctl {
	n := &Nerdctl{
		Version: version,
	}

	svc := n.base().AsService(dagger.ContainerAsServiceOpts{
		Args:                     []string{"containerd"},
		InsecureRootCapabilities: true,
	})
	n.Container = n.base().WithServiceBinding("containerd", svc)

	return n
}

func (n *Nerdctl) base() *dagger.Container {
	repo := dag.Git("https://github.com/containerd/nerdctl.git")
	var ref *dagger.GitRef
	if n.Version != "" {
		ref = repo.Tag(n.Version)
	} else {
		// HACK: dagger v0.18.12 cannot build HEAD, since it uses optional dockerfile secrets
		// see https://github.com/dagger/dagger/pull/10675
		// ref = repo.Head()
		ref = repo.Tag("v2.1.2")
	}
	ctr := ref.Tree().DockerBuild()

	ctr = ctr.
		WithMountedCache("/run/containerd", dag.CacheVolume("run-containerd")).
		WithMountedCache("/var/lib/containerd", dag.CacheVolume("containerd")).
		WithMountedCache("/var/lib/buildkit", dag.CacheVolume("buildkit")).
		WithMountedCache("/var/lib/containerd-stargz-grpc", dag.CacheVolume("containerd-stargz-grpc")).
		WithMountedCache("/var/lib/nerdctl", dag.CacheVolume("nerdctl"))
	return ctr
}

func (n *Nerdctl) Reset(ctx context.Context) (*Nerdctl, error) {
	var err error
	n.Container, err = n.Container.
		WithExec([]string{"sh", "-c", "nerdctl rm -f $(nerdctl ps -aq)"}, dagger.ContainerWithExecOpts{Expect: dagger.ReturnTypeAny}).
		WithExec([]string{"sh", "-c", "ctr image rm -f $(ctr image ls -q)"}, dagger.ContainerWithExecOpts{Expect: dagger.ReturnTypeAny}).
		WithExec([]string{"sh", "-c", "ctr content rm -f $(ctr content ls -q)"}, dagger.ContainerWithExecOpts{Expect: dagger.ReturnTypeAny}).
		Sync(ctx)
	return n, err
}

func (n *Nerdctl) WithImage(ctr *dagger.Container) *Nerdctl {
	target := "/tmp/import.tar"
	n.Container = n.Container.
		WithMountedFile(target, ctr.AsTarball()).
		WithExec([]string{"nerdctl", "load", "-i", target})
	return n
}

func (n *Nerdctl) WithImageRef(ref string) *Nerdctl {
	n.Container = n.Container.WithExec([]string{"nerdctl", "pull", "--", ref})
	return n
}

func (n *Nerdctl) WithoutImage(ref string) *Nerdctl {
	n.Container = n.Container.WithExec([]string{"nerdctl", "image", "rm", "--", ref})
	return n
}
