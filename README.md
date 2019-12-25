# ghet

A small utility for finding releases and downloading from GitHub repositories.

 * to support lookup of the latest version
 * to support resolution of dynamic version constraints
 * to support private repositories
 * to simplify multi-step `wget`, `shasum`, and `chmod` steps
 * to simplify data integrity verifications such as checksums

## Command Line

```
ghet [options] RELEASE [ASSET...]

Get something from a GitHub repository.

Global Options

  --dry-run  make no changes, but show what would happen
  --quiet    suppress status reporting
  --server   use a custom GitHub Server ($GITHUB_SERVER)
  --token    use a specific GitHub authentication token ($GITHUB_TOKEN)
  --verbose  increase status reporting

Release Options

  RELEASE  OWNER/REPOSITORY            use latest matching version
  RELEASE  OWNER/REPOSITORY@COMMITISH  matching a specific tag, branch, or commit

  --exclude-stable            exclude stable releases
  --include-drafts            include draft releases
  --include-pre-releases      include pre-releases
  --version-match=CONSTRAINT  require a version constraint

Asset Options

  ASSET  NAME            glob-friendly asset name to download
  ASSET  LOCALPATH=NAME  use a specific local path for an asset download

  --ignore-missing[=GLOB]  if an asset is not found, skip it rather than failing
  --mode=LIST              list of assets to use where LIST is:
                            * assets (default) - user-uploaded files from the release
                            * archives - GitHub-generated archives of the source files
                            * blobs - files from the release source tree

Download Options

  --cd               change to directory before downloading
  --chmod            set permission mode
  --install          equivalent to: --verify=required --chmod=+x
  --parallel=N       maximum number of parallel download operations
  --to-stdout        write files to standard out rather than disk
  --url              print download URLs instead of downloading (may be signed and/or ephemeral)
  --verify=STRATEGY  perform verification where STRATEGY is:
                      * best-effort (default) - use checksum and/or signature verification when available
                      * checksum - require checksums to be found for downloaded assets
                      * signature - require public key verification
                      * all - require checksum and signature verification
                      * required - require checksum and/or signature verification to be available
                      * none
```

### Examples

```
# download all assets of a release
$ ghet kubernetes/kubernetes
kubernetes.tar.gz: download: kubernetes.tar.gz: OK

# download the latest ssoca linux client
$ ghet dpb587/ssoca --install /usr/local/bin/ssoca=ssoca-client-*-linux-amd64
/usr/local/bin/ssoca: download: ssoca-client-0.18.1-linux-amd64: OK
/usr/local/bin/ssoca: verify: checksum: sha1: OK
/usr/local/bin/ssoca: install: mode: +x: OK

# download and extract a specific version of hugo
$ ghet gohugoio/hugo@v0.62.0 --to-stdout hugo_extended_*_Linux-64bit.tar.gz | tar -xzvf- -C /usr/local/bin hugo
hugo_extended_0.62.0_Linux-64bit.tar.gz: download: OK
hugo_extended_0.62.0_Linux-64bit.tar.gz: verify: checksum: sha256: OK
hugo

# download with more complex version requirements
$ ghet cloudfoundry/bosh-cli --version-match=^6.1 --install /usr/local/bin/bosh=bosh-cli-*-linux-amd64
/usr/local/bin/bosh: download: bosh-cli-6.1.1-linux-amd64: OK
/usr/local/bin/bosh: verify: checksum: sha256: OK
/usr/local/bin/bosh: install: mode: +x: OK

# get a file from the source tree of a version
$ ghet kubernetes/kubernetes --mode=blobs go.mod
```

## Alternatives

 * `wget`/`shasum`/`chmod` -- requires manually building commands
 * [`hub release download ...`](https://github.com/github/hub/blob/3344f0cec5672ed262ec65e5efa4d91e4a6b26db/commands/release.go#L24) -- requires an existing git working directory

## Technical Notes

### Checksum Verification

Publishing checksums are not officially supported by GitHub. The following methods seem common enough to try and support.

 * extract checksum code block from release notes (e.g. [ssoca](https://github.com/dpb587/ssoca/releases/tag/v0.18.1))
 * look for a `*checksums.txt` asset with all asset checksums (e.g. [hugo](https://github.com/gohugoio/hugo/releases/tag/v0.62.0))
 * look for matching `*.(sha1|sha256|sha512)` asset for each asset (e.g. [concourse](https://github.com/concourse/concourse/releases/tag/v5.5.7))

### Discussion

Should `--chmod` actually be supported?

 * – not very good for generic cases when downloading multiple files of different types
 * + helpful for binary installation
 * ~ perhaps it could support using it multiple times to apply to files after it

Should `--mode` become a subcommand? i.e. `ghet asset x/y z`

 * – fewer arguments keeps the command simpler
 * + may make more sense for additional types that might have other options (e.g. docker, npm)
 * + still readable when adding it
 * ~ perhaps it could be optional and default to `asset`

Support for GitHub Actions?

 * + helpful for a couple cases already
 * ~ not sure what syntax should look like (i.e. args vs inputs vs dynamic file lists vs ghet-installer)

## License

TBD
