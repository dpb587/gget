# gget

A small utility for getting files from git repositories.

You may find this useful for:

 * automating downloads from the command line;
 * downloading individual source files without a full clone; and
 * avoiding hard-coded versions or download URLs.

With notable support for:

 * public and private repositories;
 * user-uploaded release assets, source files, and `git export` archives;
 * tag, branch, and commit-based references;
 * convenience methods for renaming files, marking executable, and verifying checksums; and
 * [GitHub](https://github.com/) repositories (with [Bitbucket](https://bitbucket.org/) and [GitLab](https://about.gitlab.com/) coming soon).

## Command Line Usage

Provide the repository you want to download from as the first argument. By default, all user-uploaded assets of the latest release will be downloaded.

    gget github.com/gohugoio/hugo

Include a tag to download files from something other than the latest published version.

    gget github.com/gohugoio/hugo@v0.63.1

Provide file names (or globs) as additional arguments to limit the files are downloaded.

    gget github.com/gohugoio/hugo 'hugo_extended_*_Linux-ARM.deb'

Use the `--exclude=` option to avoid files with overlapping matches.

    gget github.com/gohugoio/hugo --exclude='*extended*' 'hugo_*_Linux-ARM.deb'

Prefix remote file names with a custom local file path to use an alternative download location. Use the `--exec` option to mark a download as executable.

    gget --exec github.com/stedolan/jq /usr/local/bin/jq=jq-osx-amd64

Use the `--type=` option to download files other than user-uploaded release assets. Use `archive` to access zip or tarball archives of the repository files.

    gget --type=archive github.com/stedolan/jq '*.zip'

Use the `blob` type to download individual repository files. In addition to tags, branch and commit references may be used for these types.

    gget --type=blob github.com/stedolan/jq@jq-1.5-branch README.md

Use `--help` to see all options and learn more about advanced usage.

    gget --help

## Services

The following services are supported through their APIs:

 * **GitHub** – personal access tokens may be set via `$GITHUB_TOKEN` or a `.netrc` password

## Installation

Binaries for Linux, macOS, and Windows can be downloaded from the [releases](https://github.com/dpb587/gget/releases) page.

A [Homebrew](https://brew.sh/) recipe is available for Linux and macOS.

```
brew install dpb587/tap/gget
```

Use `go get` to build the latest development version.

```
go get -u github.com/dpb587/gget
```

## Alternatives

 * `wget`/`curl` -- if you want to manually maintain download URLs and private signing
 * [`hub release download ...`](https://github.com/github/hub) -- if you already have `git` configured and a cloned GitHub repository

## License

[MIT License](LICENSE)
