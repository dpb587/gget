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

Pass the repository you want to download from. By default, all user-uploaded assets of the latest release will be downloaded.

    gget github.com/gohugoio/hugo

Include a tag to download files from something other than the latest published version.

    gget github.com/gohugoio/hugo@v0.63.1

Provide file names (or globs) to limit the files downloaded. Use `--exclude=` to avoid files with overlapping matches.

    gget github.com/gohugoio/hugo --exclude='*extended*' 'hugo_*_Linux-ARM.deb'

Prefix remote file names with a local path to use a custom location. Use `--executable` to mark it executable.

    gget --executable github.com/stedolan/jq /usr/local/bin/jq=jq-osx-amd64

The `--ref-*` options may be used when no tag/ref is passed with the repository. Use `--ref-version=` to find the latest match of a version constraint.

    gget --ref-version=1.0.x github.com/prometheus/pushgateway '*dragonfly*'

Use `--type=` to download files other than user-uploaded release assets. Use `archive` to access zip or tarball archives of the repository files.

    gget --type=archive github.com/stedolan/jq '*.zip'

Use the `blob` type to download repository source files. Branch and commit references may also be used for these types.

    gget --type=blob github.com/stedolan/jq@jq-1.5-branch README.md

Use `--help` to see additional options and learn more for advanced usage.

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

```dockerfile
FROM docker.pkg.github.com/dpb587/gget/gget as gget
RUN gget --executable github.com/cloudfoundry/bosh-cli bosh=bosh-cli-*-linux-amd64
RUN gget --executable github.com/cloudfoundry/bosh-bootloader bbl=bbl-*_linux_x86-64
RUN gget --stdout github.com/pivotal-cf/om om-linux-*.tar.gz | tar -xzf- om

FROM ubuntu
COPY --from=gget /result/* /usr/local/bin/
# ...everything else for your image...
```

## Services

The following services are supported through their APIs:

 * **GitHub** – personal access tokens may be set via `$GITHUB_TOKEN` or a `.netrc` password
 * **GitLab** – personal access tokens may be set via `$GITLAB_TOKEN` or a `.netrc` password

## Technical Notes

### Checksum Verification

When downloading files, `gget` attempts to validate checksums when they are found for files (erroring if they do not match). Since checksums are generally not an official feature for repository assets, this is a convention-based approach.

 * Algorithms: `sha512`, `sha256`, `sha1`, `md5`
 * Format:
    * `*sum` command output - `{checksum}  {file}`
 * Sources:
    * release notes - code block or code fence
    * sibling files with an algorithm suffix (case-insensitive) - `*.{algorithm}`
    * checksum list files (case-insensitive) - `checksum`, `checksums`, `*checksums.txt`, `{algorithm}sum.txt`, `{algorithm}sums.txt`

Some personal recommendations/learnings/preferences:

 * do not use `sha1` or `md5` - they are considered weak (`gget` uses the strongest checksum it finds)
 * use two spaces instead of one when generating `*sum` command output - it is more widely supported by `*sum --check` tools (although `gget` supports both)
 * include checksums in the release notes - they are then stored in a different backend than the assets being verified
 * use a single checksums file rather than one per file, per algorithm - for predictable references and avoiding individual file checksum downloads requiring API requests

## Alternatives

 * `wget`/`curl` -- if you want to manually maintain version download URLs and private signing
 * [`hub release download ...`](https://github.com/github/hub) -- if you already have `git` configured and a cloned GitHub repository

## License

[MIT License](LICENSE)
