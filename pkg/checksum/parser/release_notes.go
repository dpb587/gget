package parser

import (
	"regexp"
	"strings"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/pkg/errors"
)

var codefence = regexp.MustCompile("(?mU)```([^`]+)```")
var codeindent = regexp.MustCompile("    ([a-f0-9]{40,64})\\s+([^\\s]+)")

// var codefence = regexp.MustCompile("(?m)```")

func ParseReleaseNotes(releaseNotes string) *checksum.InMemoryManager {
	res := parseReleaseNotesCodefence(releaseNotes)
	if res != nil {
		return res
	}

	return parseReleaseNotesCodeindent(releaseNotes)
}

func parseReleaseNotesCodefence(releaseNotes string) *checksum.InMemoryManager {
	releaseNotesSubmatches := codefence.FindAllStringSubmatch(releaseNotes, -1)

	if len(releaseNotesSubmatches) == 0 {
		return nil
	}

	manager := checksum.NewInMemoryManager()

	for _, releaseNotesSubmatch := range releaseNotesSubmatches {
		checksums := strings.Split(strings.TrimSpace(releaseNotesSubmatch[1]), "\n")

		for _, checksumLine := range checksums {
			checksumSplit := strings.Fields(strings.TrimSpace(checksumLine))

			if len(checksumSplit[0]) < 16 {
				continue
			}

			checksum, err := checksum.GuessChecksum(checksumSplit[0])
			if err != nil {
				panic(errors.Wrapf(err, "unexpected checksum %s", checksumSplit[0]))
			}

			manager.SetChecksum(checksumSplit[1], checksum)
		}
	}

	return manager
}

func parseReleaseNotesCodeindent(releaseNotes string) *checksum.InMemoryManager {
	releaseNotesSubmatch := codeindent.FindAllStringSubmatch(releaseNotes, -1)

	if len(releaseNotesSubmatch) == 0 {
		// fmt.Printf("%s\n", releaseNotes)
		return nil
	}

	manager := checksum.NewInMemoryManager()

	for _, match := range releaseNotesSubmatch {
		checksum, err := checksum.GuessChecksum(match[1])
		if err != nil {
			panic("TODO")
		}

		manager.SetChecksum(match[2], checksum)
	}

	return manager
}
