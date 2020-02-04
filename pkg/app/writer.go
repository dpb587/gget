package app

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"rsc.io/goversion/version"
)

func WriteVersion(w io.Writer, v Version) {
	fmt.Fprintf(w, "%s\n", v.String())

	return
}

func WriteVersionVerbose(w io.Writer, v Version, self string) {
	fmt.Fprintf(w, "app\t%s\t%s\t(commit %s; built %s)\n", v.Name, v.Semver, v.Commit, v.Built.Format(time.RFC3339))
	fmt.Fprintf(w, "build\t%s\t%s\n", "go", runtime.Version())

	vv, err := version.ReadExe(os.Args[0])
	if err != nil {
		return
	}

	for _, line := range strings.Split(strings.TrimSpace(vv.ModuleInfo), "\n") {
		row := strings.Split(line, "\t")
		if row[0] != "dep" {
			continue
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t(hash %s)\n", row[0], row[1], row[2], row[3])
	}
}
