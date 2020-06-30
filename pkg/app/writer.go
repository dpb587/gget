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

type versionRecord struct {
	Type      string
	Component string
	Version   string
	Metadata  []string
}

func WriteVersion(w io.Writer, self string, v Version, verbosity int) {
	if verbosity == 0 {
		fmt.Fprintf(w, "%s\n", v.String())

		return
	}

	records := []versionRecord{
		{
			Type:      "app",
			Component: v.Name,
			Version:   v.Semver,
			Metadata: []string{
				fmt.Sprintf("commit %s", v.Commit),
				fmt.Sprintf("built %s", v.Built.Format(time.RFC3339)),
			},
		},
		{
			Type:      "runtime",
			Component: "go",
			Version:   runtime.Version(),
			Metadata: []string{
				fmt.Sprintf("arch %s", runtime.GOARCH),
				fmt.Sprintf("os %s", runtime.GOOS),
			},
		},
	}

	if vv, err := version.ReadExe(os.Args[0]); err == nil {
		for _, line := range strings.Split(strings.TrimSpace(vv.ModuleInfo), "\n") {
			row := strings.Split(line, "\t")
			if row[0] != "dep" {
				continue
			}

			records = append(
				records,
				versionRecord{
					Type:      "dep",
					Component: row[1],
					Version:   row[2],
					Metadata: []string{
						fmt.Sprintf("hash %s", row[3]),
					},
				},
			)
		}
	}

	var f func(versionRecord) string = func(r versionRecord) string {
		return fmt.Sprintf("%s\t%s\t%s\n", r.Type, r.Component, r.Version)
	}

	if verbosity > 1 {
		f = func(r versionRecord) string {
			return fmt.Sprintf("%s\t%s\t%s\t%s\n", r.Type, r.Component, r.Version, strings.Join(r.Metadata, "; "))
		}
	}

	for _, record := range records {
		fmt.Fprintf(w, f(record))
	}
}
