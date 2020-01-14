package checksums

import (
	"fmt"
	"regexp"
	"strings"
)

var codefence = regexp.MustCompile("(?m)(```[\\w\\d]*\\s*)(([a-f0-9]{40,64})\\s+([^\\s]+)\\s*)+(```)")
var codeindent = regexp.MustCompile("    ([a-f0-9]{40,64})\\s+([^\\s]+)")

// var codefence = regexp.MustCompile("(?m)```")

func ParseReleaseNotes(releaseNotes string) ReleaseAssets {
	res := parseReleaseNotesCodefence(releaseNotes)
	if len(res) > 0 {
		return res
	}

	return parseReleaseNotesCodeindent(releaseNotes)
}

func parseReleaseNotesCodefence(releaseNotes string) ReleaseAssets {
	releaseNotesSubmatch := codefence.FindAllStringSubmatch(releaseNotes, -1)

	if len(releaseNotesSubmatch) == 0 {
		return nil
	}

	checksums := strings.Split(
		strings.TrimSpace(
			strings.TrimSuffix(
				strings.TrimPrefix(
					releaseNotesSubmatch[0][0],
					releaseNotesSubmatch[0][1],
				),
				releaseNotesSubmatch[0][5],
			),
		),
		"\n",
	)

	var releaseSHAs ReleaseAssets

	for _, checksumLine := range checksums {
		checksumSplit := strings.Fields(strings.TrimSpace(checksumLine))

		releaseSHAs = append(releaseSHAs, ReleaseAsset{
			SHA:  checksumSplit[0],
			Name: checksumSplit[1],
		})
	}

	return releaseSHAs
}

func parseReleaseNotesCodeindent(releaseNotes string) ReleaseAssets {
	releaseNotesSubmatch := codeindent.FindAllStringSubmatch(releaseNotes, -1)

	if len(releaseNotesSubmatch) == 0 {
		fmt.Printf("%s\n", releaseNotes)
		return nil
	}

	var releaseSHAs ReleaseAssets

	for _, match := range releaseNotesSubmatch {
		releaseSHAs = append(releaseSHAs, ReleaseAsset{
			SHA:  match[1],
			Name: match[2],
		})
	}

	return releaseSHAs
}
