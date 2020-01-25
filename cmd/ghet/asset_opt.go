package ghet

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type AssetOpt struct {
	RemoteMatch string
	LocalPath   string
}

func (o *AssetOpt) Resolve(remote string) (AssetOpt, bool) {
	match, _ := filepath.Match(o.RemoteMatch, remote)
	if !match {
		return AssetOpt{}, false
	}

	res := AssetOpt{
		RemoteMatch: remote,
		LocalPath:   o.LocalPath,
	}

	if res.LocalPath == "" {
		res.LocalPath = res.RemoteMatch
	} else if strings.HasSuffix(res.LocalPath, "/") {
		res.LocalPath = filepath.Join(res.LocalPath, res.RemoteMatch)
	}

	return res, true
}

func (o *AssetOpt) UnmarshalFlag(data string) error {
	dataSplit := strings.SplitN(data, "=", 2)

	if len(dataSplit) == 2 {
		o.RemoteMatch = dataSplit[1]
		o.LocalPath = dataSplit[0]
	} else {
		o.RemoteMatch = dataSplit[0]
		o.LocalPath = ""
	}

	_, err := filepath.Match(o.RemoteMatch, "test")
	if err != nil {
		return errors.Wrap(err, "expected valid asset matcher")
	}

	return nil
}
