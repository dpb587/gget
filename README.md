# gget

An easier way to find and automate file downloads from GitHub and GitLab repositories. Learn more from the examples and documentation at [gget.io](https://gget.io/).

 * Standalone CLI - no `git` or local clones required
 * Public & Private Servers - supporting both [GitHub](https://github.com/) and [GitLab](https://gitlab.com/)
 * Public & Private Repos - API-based access for private resources
 * Tags, Branches, and Commits - also semver-based constraint matching
 * Archives, Assets, and Blobs - download any type of resource from repos
 * Built-in File Operations - easy rename, binary, and stream options
 * Checksum Verification - verify downloads by SHA (if published)
 * Open-Source Project - contribute and discuss at [dpb587/gget](https://github.com/dpb587/gget)

## Command Line Usage

Pass the repository and file globs to match and download.

```
$ gget github.com/gohugoio/hugo 'hugo_extended_*_Linux-64bit.deb'
Found 1 file (13.9M) from github.com/gohugoio/hugo@v0.73.0
√ hugo_extended_0.73.0_Linux-64bit.deb done (sha256 OK)
```

Use `--help` to see the full list of options for more advanced usage.

### Installation

Binaries for Linux, macOS, and Windows can be downloaded from the [releases](https://github.com/dpb587/gget/releases) page. A [Homebrew](https://brew.sh/) recipe is also available for Linux and macOS.

```
brew install dpb587/tap/gget
```

## Docker Usage

The `gget` image can be used as a build stage to download assets for a later stage.

```dockerfile
FROM dpb587/gget as gget
RUN gget --executable github.com/cloudfoundry/bosh-cli --ref-version=5.x bosh=bosh-cli-*-linux-amd64
RUN gget --executable github.com/cloudfoundry/bosh-bootloader bbl=bbl-*_linux_x86-64
RUN gget --stdout github.com/pivotal-cf/om om-linux-*.tar.gz | tar -xzf- om

FROM ubuntu
COPY --from=gget /result/* /usr/local/bin/
# ...everything else for your image...
```

## Technical Notes

### Checksum Verification

When downloading files, `gget` attempts to validate checksums when they are found for files (erroring if they do not match). Since checksums are not an official feature for repository assets, the following common conventions are being used.

 * Algorithms: `sha512`, `sha256`, `sha1`, `md5`
 * File Format: `shasum`/`md5sum` output (`{checksum}  {file}`)
 * Sources:
    * release notes - code block or code fence
    * sibling files with an algorithm suffix (case-insensitive) - `*.{algorithm}`
    * checksum list files (case-insensitive) - `checksum`, `checksums`, `*checksums.txt`, `{algorithm}sum.txt`, `{algorithm}sums.txt`

## Alternatives

 * `wget`/`curl` -- if you want to manually maintain version download URLs and private signing
 * [`hub release download ...`](https://github.com/github/hub) -- if you already have `git` configured and a cloned GitHub repository
 * [`fetch`](https://github.com/gruntwork-io/fetch) -- similar tool and capabilities for GitHub (discovered this after gget@0.5.2)

## License

[MIT License](LICENSE)
