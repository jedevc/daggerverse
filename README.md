# Daggerverse

> **Warning**
>
> This daggerverse repo requires a build of Dagger from https://github.com/dagger/dagger/pull/5966.

## GitHub

Get details of GitHub releases.

Command-line usage:

```
$ dagger -m github call get-latest-release --repo="dagger/dagger" body
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

## Hugo

Build hugo sites.

Command-line usage:

```
$ dagger -m hugo export --export-path=./public build --target=./src --hugo-version=latest
✔ dagger download build [3.04s]
┃ Asset exported to "./public"                                                                                                                             
• Engine: b87daf4a7392 (version devel ())
⧗ 3.83s ✔ 85 ∅ 52
```

## Npm

Install NPM packages to a directory.

Command-line usage:

```
$ ls -lh ./target
total 36K
-rw-r--r-- 1 jedevc jedevc 249 Oct 26 15:11 package.json
-rw-r--r-- 1 jedevc jedevc 32K Oct 26 15:11 package-lock.json

$ dagger -m npm export --export-path=./target node-modules --dir=./target
✔ dagger download node-modules [4.12s]
┃ Asset exported to "./target"                                                                                                                                 
• Engine: b87daf4a7392 (version devel ())
⧗ 4.85s ✔ 54 ∅ 9

$ ls -lh ./target
total 40K
drwxr-xr-x 56 jedevc jedevc 4.0K Oct 26 15:12 node_modules
-rw-r--r--  1 jedevc jedevc  249 Oct 26 15:11 package.json
-rw-r--r--  1 jedevc jedevc  32K Oct 26 15:11 package-lock.json
```