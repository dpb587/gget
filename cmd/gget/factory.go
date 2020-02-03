package gget

func NewCommand() *Command {
	return &Command{
		Runtime: &Runtime{},
	}
}
