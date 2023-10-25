package main

import (
	"path"
)

type Npm struct{}

func (npm *Npm) NodeModules(dir *Directory) *Directory {
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
	return dir.WithDirectory("node_modules", nodeModules)
}
