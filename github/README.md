# GitHub

Get details of GitHub releases.

Command-line usage:

```
$ dagger -m github.com/jedevc/daggerverse/github call \
    get-latest-release --repo="dagger/dagger" body
## v0.10.1 - 2024-03-05


### Added
- Allow passing git URLs to `dagger call` file type args by @jedevc in https://github.com/dagger/dagger/pull/6769
- Support privileges and nesting in default terminal command by @TomChv in https://github.com/dagger/dagger/pull/6805

### Fixed
- Fix panic in Contents for massive files by @jedevc in https://github.com/dagger/dagger/pull/6772
- Dagger go modules default to the module name instead of "main" by @jedevc in https://github.com/dagger/dagger/pull/6774
- Fix a regression where secrets used with dockerBuild could error out by @jedevc in https://github.com/dagger/dagger/pull/6809
- Fix goroutine and memory leaks in engine by @sipsma in https://github.com/dagger/dagger/pull/6760
- Fix potential name clash with "Client" in Go functions by @jedevc in https://github.com/dagger/dagger/pull/6716

### What to do next?
- Read the [documentation](https://docs.dagger.io)
- Join our [Discord server](https://discord.gg/dagger-io)
- Follow us on [Twitter](https://twitter.com/dagger_io)
```

Go SDK usage:

```
$ dagger install github.com/jedevc/daggerverse/github
```

```go
func (m *MyModule) GetDaggerChecksums(
	ctx context.Context,
	// +optional
	// +default="latest"
	version string,
) (*File, error) {
	var release *GithubRelease
	switch version {
	case "latest":
		release = dag.Github().GetLatestRelease("dagger/dagger")
	default:
		release = dag.Github().GetRelease("dagger/dagger", version)
	}

	assets, err := release.Assets(ctx)
	if err != nil {
		return nil, err
	}
	for _, asset := range assets {
		name, err := asset.Name(ctx)
		if err != nil {
			return nil, err
		}
		if name == "checksums.txt" {
			download, err := asset.DownloadURL(ctx)
			if err != nil {
				return nil, err
			}
			return dag.HTTP(download), nil
		}
	}

	return nil, fmt.Errorf("checksums.txt not found")
}

```