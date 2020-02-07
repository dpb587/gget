package gget

import "github.com/dpb587/gget/pkg/app"

func NewCommand(app app.Version) *Command {
	return &Command{
		Runtime: NewRuntime(app),
	}
}
