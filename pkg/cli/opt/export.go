package opt

import (
	"fmt"
	"strings"

	"github.com/dpb587/gget/pkg/export"
	"github.com/pkg/errors"
)

type Export struct {
	export.Exporter
}

func (e *Export) UnmarshalFlag(data string) error {
	var exporter export.Exporter
	var template string

	if data == "json" {
		exporter = export.JSONExporter{}
	} else if strings.HasPrefix(data, "jsonpath=") {
		exporter = &export.JSONPathExporter{}
		template = strings.TrimPrefix(data, "jsonpath=")
	} else if data == "jsonpath" {
		return fmt.Errorf("jsonpath must include a template (e.g. `jsonpath='{.origin.ref}'`)")
	} else if data == "yaml" {
		exporter = export.YAMLExporter{}
	} else if data == "plain" {
		exporter = export.PlainExporter{}
		// } else if strings.HasPrefix(data, "go-template=") {
		// 	exporter = &export.GoTemplateExporter{}
		// 	template = strings.TrimPrefix(data, "go-template=")
		// } else if data == "go-template" {
		// 	return fmt.Errorf("go-template must include a template (e.g. `go-template='{{.Origin.Ref}}'`")
	} else {
		return fmt.Errorf("unsupported export type: %s", data)
	}

	if template != "" {
		tex, ok := exporter.(export.TemplatedExporter)
		if !ok {
			panic(fmt.Errorf("exporter does not support templating: %T", exporter))
		}

		err := tex.ParseTemplate(template)
		if err != nil {
			return errors.Wrap(err, "parsing export template")
		}
	}

	*e = Export{Exporter: exporter}

	return nil
}
