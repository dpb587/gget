package ghet

type Runtime struct {
	Quiet    bool `long:"quiet" description:"suppress status reporting"`
	Parallel int  `long:"parallel" description:"maximum number of parallel operations" default:"3"`

	Help    bool `long:"help" short:"h" description:"show documentation of this tool"`
	Version bool `long:"version" description:"show version of this tool"`
}
