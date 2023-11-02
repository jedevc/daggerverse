# GitHub

Get details of GitHub releases.

Command-line usage:

```
$ dagger -m github.com/jedevc/daggerverse/github call \
    get-latest-release --repo="dagger/dagger" body
✔ dagger call get-latest-release body [0.64s]
┃ ## v0.9.0 - 2023-10-20
┃
┃ ### 🔥 Breaking Changes
┃
┃ - engine: new services API with container <=> host networking, explicit start/stop by @vito in https://github.com/dagger/dagger/pull/5557
┃ - implement new conventions for IDable objects by @vito in https://github.com/dagger/dagger/pull/5881
┃
┃ ### Added
┃
┃ - engine: support multiple cache configs for upstream remote cache by @sipsma in https://github.com/dagger/dagger/pull/5730
┃
┃ ### Changed
┃
┃ - engine: reduce connection retry noise by @sipsma in https://github.com/dagger/dagger/pull/5918
┃
┃ ### Fixed
┃
┃ - engine: fix missing descriptor handlers for lazy blobs error w/ cloud cache by @sipsma in https://github.com/dagger/dagger/pull/5885
┃
┃ ### What to do next?
┃
┃ - Read the [documentation](https://docs.dagger.io)
┃ - Join our [Discord server](https://discord.gg/dagger-io)
┃ - Follow us on [Twitter](https://twitter.com/dagger_io)
• Engine: b87daf4a7392 (version devel ())
⧗ 1.56s ✔ 30 ∅ 9
```

Go SDK usage:

```
$ dagger mod use github.com/jedevc/daggerverse/github
```

```go
func (m *MyModule) GetDaggerChecksums(ctx context.Context, version Optional[string]) (*File, error) {
	var release *GithubRelease
	switch version := version.GetOr(""); version {
	case "", "latest":
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