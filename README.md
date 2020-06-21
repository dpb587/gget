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

Use `--exclude=` to avoid files with overlapping matches.

    gget github.com/gohugoio/hugo --exclude='*extended*' 'hugo_*_Linux-ARM.deb'

Prefix remote file names with a custom local file path to use an alternative download location. Use the `--executable` option to mark a download as executable.

    gget --executable github.com/stedolan/jq /usr/local/bin/jq=jq-osx-amd64

The `--ref-*` options may be used when no tag/ref is passed with the repository. Use `--ref-stability=` to limit or expand the types of releases to look for.

    gget --ref-stability=pre-release github.com/prometheus/prometheus '*dragonfly*'

Use `--ref-version=` to provide a version constraint to use instead of latest.

    gget --ref-version=1.0.x github.com/prometheus/pushgateway '*dragonfly*'

Use `--type=` to download files other than user-uploaded release assets. Use `archive` to access zip or tarball archives of the repository files.

    gget --type=archive github.com/stedolan/jq '*.zip'

Use the `blob` type to download individual repository files. Branch and commit references may also be used for these types.

    gget --type=blob github.com/stedolan/jq@jq-1.5-branch README.md

Use `--help` to see all options and learn more about advanced usage.

    Usage:
      gget HOST/OWNER/REPOSITORY[@REF] [[LOCAL-PATH=]RESOURCE-GLOB...]

    Runtime Options:
      -q, --quiet                             suppress runtime status reporting
      -v, --verbose                           increase logging verbosity (multiple)
      -h, --help                              show documentation of this command
          --version=[CONSTRAINT]              show version of this command (with optional constraint to validate)

    Repository Options:
          --service=NAME                      specific git service to use (values: github, gitlab; default: auto-detect)
          --ref-version=CONSTRAINT            version constraint(s) to require of latest (e.g. 4.x)
          --ref-stability=STABILITY           acceptable stability level(s) for latest (values: stable, pre-release, any; default: stable)
          --show-ref                          show resolved repository ref instead of downloading

    Resource Options:
          --type=TYPE                         type of resource to get (values: asset, archive, blob; default: asset)
          --ignore-missing=[RESOURCE-GLOB]    if a resource is not found, skip it rather than failing (multiple)
          --exclude=RESOURCE-GLOB             exclude resource(s) from download (multiple)
          --show-resources                    show matched resources instead of downloading

    Download Options:
          --cd=DIR                            change to directory before writing files
          --executable=[RESOURCE-GLOB]        apply executable permissions to downloads (multiple)
          --stdout                            write file contents to stdout rather than disk
          --parallel=INT                      maximum number of parallel downloads (default: 3)

    Arguments:
      HOST/OWNER/REPOSITORY[@REF]:            repository reference
      [LOCAL-PATH=]RESOURCE-GLOB:             resource name(s) to download

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
 * use a single `checksums.txt` file - for predictable usage and avoiding individual file checksum downloads requiring API requests

## Alternatives

 * `wget`/`curl` -- if you want to manually maintain version download URLs and private signing
 * [`hub release download ...`](https://github.com/github/hub) -- if you already have `git` configured and a cloned GitHub repository

## License

[MIT License](LICENSE)
