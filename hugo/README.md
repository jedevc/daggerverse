## Hugo

Build hugo sites.

Command-line usage:

```
$ dagger -m github.com/jedevc/daggerverse/hugo \
    export --output=./public build --target=./src --hugo-version=latest
✔ dagger download build [3.04s]
┃ Asset exported to "./public"
• Engine: b87daf4a7392 (version devel ())
⧗ 3.83s ✔ 85 ∅ 52
```

Go SDK usage:

```
$ dagger mod use github.com/jedevc/daggerverse/hugo
```

```go
func (m *MyModule) BuildSite(ctx context.Context, dir *Directory) *Directory {
    return dag.Hugo().Build(dir, HugoBuildOpts{
        HugoVersion:     "latest",
        DartSassVersion: "latest",
    })
}
```