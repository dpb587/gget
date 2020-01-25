package ghet

type Cmd struct {
	*Global

	Asset *AssetCmd `command:"asset" description:"fetch user-uploaded files from the release"`
}

func New() Cmd {
	o := &Global{}

	return Cmd{
		Global: o,
		Asset: &AssetCmd{
			Global: o,
		},
	}
}
