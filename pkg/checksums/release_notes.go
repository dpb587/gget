package checksums

import (
	"fmt"
	"regexp"
	"strings"
)

var codeblock = regexp.MustCompile("(?m)(```[\\w\\d]*\\s*)(([a-f0-9]{40,64})\\s*([^\\s]+)\\s*)+(```)")

// var codeblock = regexp.MustCompile("(?m)```")

func ParseReleaseNotes(releaseNotes string) ReleaseAssets {
	releaseNotesSubmatch := codeblock.FindAllStringSubmatch(releaseNotes, -1)

	if len(releaseNotesSubmatch) == 0 {
		fmt.Printf("---\n%s\n---", releaseNotes)

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
