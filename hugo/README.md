## Hugo

Build hugo sites.

Command-line usage:

```
$ dagger -m github.com/jedevc/daggerverse/hugo call \
    build --target=./src --hugo-version=latest export --path=./public
true
```

Go SDK usage:

```
$ dagger install github.com/jedevc/daggerverse/hugo
```

```go
func (m *MyModule) BuildSite(ctx context.Context, dir *Directory) *Directory {
    return dag.Hugo().Build(dir, HugoBuildOpts{
        HugoVersion:     "latest",
        DartSassVersion: "latest",
    })
}
```