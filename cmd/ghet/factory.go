package ghet

func NewCommand() *Command {
	return &Command{
		Runtime: &Runtime{},
	}
}
