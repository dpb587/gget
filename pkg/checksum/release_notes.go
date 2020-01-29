package checksum

import (
	"regexp"
	"strings"

	"github.com/dpb587/ghet/pkg/model"
	"github.com/pkg/errors"
)

var codefence = regexp.MustCompile("(?mU)```([^`]+)```")
var codeindent = regexp.MustCompile("    ([a-f0-9]{40,64})\\s+([^\\s]+)")

// var codefence = regexp.MustCompile("(?m)```")

func ParseReleaseNotes(releaseNotes string) model.ChecksumMap {
	res := parseReleaseNotesCodefence(releaseNotes)
	if len(res) > 0 {
		return res
	}

	return parseReleaseNotesCodeindent(releaseNotes)
}

func parseReleaseNotesCodefence(releaseNotes string) model.ChecksumMap {
	releaseNotesSubmatches := codefence.FindAllStringSubmatch(releaseNotes, -1)

	if len(releaseNotesSubmatches) == 0 {
		return nil
	}

	releaseSHAs := model.ChecksumMap{}

	for _, releaseNotesSubmatch := range releaseNotesSubmatches {
		checksums := strings.Split(strings.TrimSpace(releaseNotesSubmatch[1]), "\n")

		for _, checksumLine := range checksums {
			checksumSplit := strings.Fields(strings.TrimSpace(checksumLine))

			if len(checksumSplit[0]) < 16 {
				continue
			}

			checksum, err := GuessChecksum(checksumSplit[0])
			if err != nil {
				panic(errors.Wrapf(err, "unexpected checksum %s", checksumSplit[0]))
			}

			releaseSHAs[checksumSplit[1]] = checksum
		}
	}

	return releaseSHAs
}

func parseReleaseNotesCodeindent(releaseNotes string) model.ChecksumMap {
	releaseNotesSubmatch := codeindent.FindAllStringSubmatch(releaseNotes, -1)

	if len(releaseNotesSubmatch) == 0 {
		// fmt.Printf("%s\n", releaseNotes)
		return nil
	}

	releaseSHAs := model.ChecksumMap{}

	for _, match := range releaseNotesSubmatch {

		checksum, err := GuessChecksum(match[1])
		if err != nil {
			panic("TODO")
		}

		releaseSHAs[match[2]] = checksum
	}

	return releaseSHAs
}
