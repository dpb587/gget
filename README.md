# gget

A small utility for downloading files from git repositories.

You may find this useful for:

 * automating downloads from the command line;
 * downloading individual source files without a full clone; and
 * avoiding hard-coded versions or download URLs.

With notable support for:

 * public and private repositories;
 * user-managed release assets, source files, and `git export` archives;
 * tag, branch, and commit-based references;
 * convenience methods for renaming files, marking executable, and verifying checksums; and
 * [GitHub](https://github.com/) and [GitLab](https://gitlab.com/) repositories (with hopes of adding [Bitbucket](https://bitbucket.org/)).

## Command Line Usage

Provide the repository you want to download from as the first argument. By default, all user-uploaded assets of the latest release will be downloaded.

    gget github.com/gohugoio/hugo

Include a tag to download files from something other than the latest published version.

    gget github.com/gohugoio/hugo@v0.63.1

Provide file names (or globs) as additional arguments to limit the files are downloaded.

    gget github.com/gohugoio/hugo 'hugo_extended_*_Linux-ARM.deb'

Use the `--exclude=` option to avoid files with overlapping matches.

    gget github.com/gohugoio/hugo --exclude='*extended*' 'hugo_*_Linux-ARM.deb'

Prefix remote file names with a custom local file path to use an alternative download location. Use the `--executable` option to mark a download as executable.

    gget --executable github.com/stedolan/jq /usr/local/bin/jq=jq-osx-amd64

Use the `--type=` option to download files other than user-uploaded release assets. Use `archive` to access zip or tarball archives of the repository files.

    gget --type=archive github.com/stedolan/jq '*.zip'

Use the `blob` type to download individual repository files. Branch and commit references may also be used for these types.

    gget --type=blob github.com/stedolan/jq@jq-1.5-branch README.md

Use `--help` to see all options and learn more about advanced usage.

    Usage:
      gget HOST/OWNER/REPOSITORY[@REF] [LOCAL-PATH=]RESOURCE...

    Runtime Options:
      -q, --quiet                           suppress runtime status reporting
      -v, --verbose                         increase logging verbosity
          --parallel=                       maximum number of parallel operations (default: 3)
          --service=                        specific git service to use (i.e. github, gitlab; default: auto-detect)
      -h, --help                            show documentation of this command
          --version                         show version of this command

    Resource Options:
          --type=                           type of resource to get (e.g. asset, archive, blob) (default: asset)
          --ignore-missing=[RESOURCE]       if a resource is not found, skip it rather than failing (glob-friendly; multiple)
          --exclude=[RESOURCE]              exclude resource(s) from download (glob-friendly; multiple)

    Download Options:
          --show-ref                        list matched repository ref instead of downloading
          --show-resources                  list matched resources instead of downloading
          --cd=                             change to directory before writing files
          --executable=[RESOURCE]           apply executable permissions to downloads (glob-friendly; multiple)
          --stdout                          write file contents to stdout rather than disk

    Arguments:
      HOST/OWNER/REPOSITORY[@REF]:          release reference
      [LOCAL-PATH=]RESOURCE:                resource name(s) to download (glob-friendly)

### Installation

Binaries for Linux, macOS, and Windows can be downloaded from the [releases](https://github.com/dpb587/gget/releases) page.

A [Homebrew](https://brew.sh/) recipe is available for Linux and macOS.

```
brew install dpb587/tap/gget
```

Use `go get` to build the latest development version.

```
go get -u github.com/dpb587/gget
```

## Docker Usage

The `gget` image can be used as a build stage to download assets for a later stage.

```
FROM docker.pkg.github.com/dpb587/gget/gget as gget
RUN gget --executable github.com/cloudfoundry/bosh-cli bosh=bosh-cli-*-linux-amd64
RUN gget --executable github.com/cloudfoundry/bosh-bootloader bbl=bbl-*_linux_x86-64
RUN gget --stdout github.com/pivotal-cf/om om-linux-*.tar.gz | tar -xzf- om

FROM ubuntu
COPY --from=gget /result/* /usr/local/bin/
...everything else for your image...
```

## Services

The following services are supported through their APIs:

 * **GitHub** – personal access tokens may be set via `$GITHUB_TOKEN` or a `.netrc` password
 * **GitLab** – personal access tokens may be set via `$GITLAB_TOKEN` or a `.netrc` password

## Alternatives

 * `wget`/`curl` -- if you want to manually maintain version download URLs and private signing
 * [`hub release download ...`](https://github.com/github/hub) -- if you already have `git` configured and a cloned GitHub repository

## License

[MIT License](LICENSE)
