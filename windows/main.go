package main

import (
	"context"
	"dagger/windows/internal/dagger"
	"strings"
)

type Windows struct {
	// +private
	Ctr *dagger.Container

	Username string
	Password string
}

func New(
	// +optional
	// +default="2022"
	version string,
	// +optional
	// +default="4G"
	ram string,
	// +optional
	// +default="2"
	cpu string,
	// +optional
	// +default="64G"
	disk string,
	// +optional
	// +default="dagger"
	username string,
	// +optional
	// +default="dagger"
	password string,
) *Windows {
	startup := strings.TrimSpace(`
dism /Online /Add-Capability /CapabilityName:OpenSSH.Server
sc config sshd start=auto
net start sshd
	`)
	return &Windows{
		Ctr: dag.Container().From("dockurr/windows").
			WithExposedPort(8006).
			// WithExposedPort(22).
			WithEnvVariable("VERSION", version).
			WithEnvVariable("RAM_SIZE", ram).
			WithEnvVariable("CPU_CORES", cpu).
			WithEnvVariable("DISK_SIZE", disk).
			WithEnvVariable("USERNAME", username).
			WithEnvVariable("PASSWORD", password).
			// WithMountedCache("/storage", dag.CacheVolume("windows-image-store")).
			WithMountedTemp("/storage").
			WithMountedDirectory("/oem", dag.Directory().WithNewFile("install.bat", startup)).
			WithExec(nil, dagger.ContainerWithExecOpts{
				InsecureRootCapabilities: true,
				UseEntrypoint:            true,
			}),
		Username: username,
		Password: password,
	}
}

func (m *Windows) Service() *dagger.Service {
	return m.Ctr.AsService()
}

func (m *Windows) SSH(ctx context.Context, svc *dagger.Service) (*SSH, error) {
	svc, err := svc.Start(ctx)
	if err != nil {
		return nil, err
	}
	return &SSH{m, svc}, nil
}

type SSH struct {
	// +private
	*Windows
	// +private
	Svc *dagger.Service
}

func (t *SSH) Execute(ctx context.Context, command string) (string, error) {
	return dag.Container().
		From("alpine").
		WithExec([]string{"apk", "add", "openssh", "sshpass"}).
		WithServiceBinding("windows", t.Svc).
		WithExec([]string{"sshpass", "-p", t.Password, "ssh", "-o", "StrictHostKeyChecking=no", "dagger@windows", command}).
		Stdout(ctx)
}
