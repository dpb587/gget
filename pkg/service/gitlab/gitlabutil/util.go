package gitlabutil

import (
	"path"

	"github.com/dpb587/gget/pkg/service"
)

func GetRepositoryID(ref service.Ref) string {
	return path.Join(ref.Owner, ref.Repository)
}
