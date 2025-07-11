package main

import (
	"cmp"
	"context"
	"crypto/rand"
	"fmt"

	"dagger/podman/internal/dagger"
)

type Podman struct {
	Version string

	Container *dagger.Container
}

func New(
	ctx context.Context,

	// +optional
	version string,
) (*Podman, error) {
	n := &Podman{
		Version: version,
	}

	port := 4000
	svc := n.base().
		WithExposedPort(port).
		AsService(dagger.ContainerAsServiceOpts{
			Args:                     []string{"podman", "system", "service", fmt.Sprintf("tcp://0.0.0.0:%d", port), "--time=0"},
			InsecureRootCapabilities: true,
		})

	dockerHost, err := svc.Endpoint(ctx, dagger.ServiceEndpointOpts{
		Scheme: "tcp",
	})
	if err != nil {
		return nil, err
	}
	n.Container = n.base().
		WithServiceBinding("podman", svc).
		WithEnvVariable("CONTAINER_HOST", dockerHost)

	return n, nil
}

func (p *Podman) base() *dagger.Container {
	return dag.Container().From("quay.io/podman/stable:" + cmp.Or(p.Version, "latest"))
}

func (p *Podman) Reset(ctx context.Context) (*Podman, error) {
	opts := dagger.ContainerWithExecOpts{
		Expect: dagger.ReturnTypeAny, // ignore failures (from non-existent objects)
	}
	_, err := p.Container.
		WithEnvVariable("CACHEBUSTER", rand.Text()).
		WithExec([]string{"sh", "-c", "podman rm -f $(podman ps -aq)"}, opts).
		WithExec([]string{"sh", "-c", "podman system prune -a"}, opts).
		Sync(ctx)
	return p, err
}
