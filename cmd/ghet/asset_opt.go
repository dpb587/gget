package ghet

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type AssetNameOpt string

func (o *AssetNameOpt) Match(remote string) bool {
	match, _ := filepath.Match(string(*o), remote)

	return match
}

func (o *AssetNameOpt) Validate() error {
	_, err := filepath.Match(string(*o), "test")
	if err != nil {
		return errors.Wrap(err, "expected valid asset matcher")
	}

	return nil
}

type AssetPathOpt struct {
	RemoteMatch AssetNameOpt
	LocalPath   string
}

func (o *AssetPathOpt) Resolve(remote string) (AssetPathOpt, bool) {
	match := o.RemoteMatch.Match(remote)
	if !match {
		return AssetPathOpt{}, false
	}

	res := AssetPathOpt{
		RemoteMatch: AssetNameOpt(remote),
		LocalPath:   o.LocalPath,
	}

	if res.LocalPath == "" {
		res.LocalPath = string(res.RemoteMatch)
	} else if strings.HasSuffix(res.LocalPath, "/") {
		res.LocalPath = filepath.Join(res.LocalPath, string(res.RemoteMatch))
	}

	return res, true
}

func (o *AssetPathOpt) UnmarshalFlag(data string) error {
	dataSplit := strings.SplitN(data, "=", 2)

	if len(dataSplit) == 2 {
		o.RemoteMatch = AssetNameOpt(dataSplit[1])
		o.LocalPath = dataSplit[0]
	} else {
		o.RemoteMatch = AssetNameOpt(dataSplit[0])
		o.LocalPath = ""
	}

	if err := o.RemoteMatch.Validate(); err != nil {
		return err
	}

	return nil
}
