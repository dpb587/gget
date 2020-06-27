package opt

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

type Export struct {
	Mode     string
	Template *template.Template
}

func (e *Export) UnmarshalFlag(data string) error {
	if data == "json" {
		e.Mode = "json"
	} else if data == "yaml" {
		e.Mode = "yaml"
	} else if data == "plain" {
		e.Mode = "plain"
	} else if strings.HasPrefix(data, "go-template:") {
		e.Mode = "go-template"

		tmpl, err := template.New("export").Parse(strings.TrimPrefix(data, "go-template:"))
		if err != nil {
			return errors.Wrap(err, "parsing go-template")
		}

		e.Mode = "go-template"
		e.Template = tmpl
	} else {
		return fmt.Errorf("unsupported export type: %s", data)
	}

	return nil
}
