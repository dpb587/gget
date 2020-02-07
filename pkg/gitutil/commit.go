package gitutil

import "regexp"

var CommitRE = regexp.MustCompile(`^[0-9a-f]{40}$`)
var PotentialCommitRE = regexp.MustCompile(`^[0-9a-f]{1,40}$`)
