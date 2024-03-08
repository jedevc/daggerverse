# Npm

Install NPM packages to a directory.

Command-line usage:

```
$ ls -lh ./target
total 36K
-rw-r--r-- 1 jedevc jedevc 249 Oct 26 15:11 package.json
-rw-r--r-- 1 jedevc jedevc 32K Oct 26 15:11 package-lock.json

$ dagger -m github.com/jedevc/daggerverse/npm call \
    node-modules --dir=./target export --path=./target
true

$ ls -lh ./target
total 40K
drwxr-xr-x 56 jedevc jedevc 4.0K Oct 26 15:12 node_modules
-rw-r--r--  1 jedevc jedevc  249 Oct 26 15:11 package.json
-rw-r--r--  1 jedevc jedevc  32K Oct 26 15:11 package-lock.json
```

Go SDK usage:

```
$ dagger mod use github.com/jedevc/daggerverse/npm
```

```go
func (m *MyModule) InstallDeps(ctx context.Context, dir *Directory) *Directory {
    return dag.Npm().NodeModules(dir)
}
```