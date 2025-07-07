// Install NPM packages.

package main

import (
	"npm/internal/dagger"
	"path"
)

type NodeModules struct{}

// Installs npm packages into the node_modules/ directory using the package files
func (npm *NodeModules) InstallNpm(dir *dagger.Directory) *dagger.Directory {
	src := dag.Directory().
		WithDirectory("/", dir, dagger.DirectoryWithDirectoryOpts{
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
	return dir.WithDirectory("node_modules", nodeModules)
}

func (npm *NodeModules) InstallYarn(dir *dagger.Directory) *dagger.Directory {
	src := dag.Directory().
		WithDirectory("/", dir, dagger.DirectoryWithDirectoryOpts{
			Include: []string{"package.json", "yarn.lock"},
		})

	srcPath := "/src"
	nodeModules := dag.Container().
		From("node:current-alpine").
		WithDirectory(srcPath, src).
		WithWorkdir(srcPath).
		WithEntrypoint(nil).
		WithExec([]string{"yarn", "install", "--frozen-lockfile"}).
		Directory(path.Join(srcPath, "node_modules"))
	return dir.WithDirectory("node_modules", nodeModules)
}
