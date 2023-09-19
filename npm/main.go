package main

import (
	"context"
	"path"
)

type Npm struct{}

func (m *Npm) InstallDependencies(ctx context.Context, dir *Directory) (*Directory, error) {
	return dir.NpmInstallDependencies(ctx)
}

func (dir *Directory) NpmInstallDependencies(ctx context.Context) (*Directory, error) {
	src := dag.Directory().
		WithDirectory("/", dir, DirectoryWithDirectoryOpts{
			Include: []string{"package.json", "package-lock.json", "npm-shrinkwrap.json"},
		})

	srcPath := "/src"
	nodeModules := dag.Container().
		From("node:current-alpine").
		WithDirectory(srcPath, src).
		WithWorkdir(srcPath).
		WithEntrypoint(nil).
		WithExec([]string{"npm", "ci"}).
		Directory(path.Join(srcPath, "node_modules"))
	dir = dir.WithDirectory("node_modules", nodeModules)
	return dir, nil
}
