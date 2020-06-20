package parser

import (
	"context"
	"io"
	"io/ioutil"
	"sync"

	"github.com/dpb587/gget/pkg/checksum"
	"github.com/pkg/errors"
)

type DeferredManager struct {
	manager checksum.WriteableManager
	opener  func(context.Context) (io.ReadCloser, error)

	loaded bool
	mutex  sync.RWMutex
}

var _ checksum.Manager = &DeferredManager{}

func NewDeferredManager(manager checksum.WriteableManager, opener func(context.Context) (io.ReadCloser, error)) checksum.Manager {
	return &DeferredManager{
		manager: manager,
		opener:  opener,
	}
}

func (m *DeferredManager) GetChecksum(ctx context.Context, resource string) (checksum.Checksum, error) {
	err := m.requireLoad(ctx)
	if err != nil {
		return nil, err
	}

	return m.manager.GetChecksum(ctx, resource)
}

func (m *DeferredManager) requireLoad(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.loaded {
		return nil
	}

	fh, err := m.opener(ctx)
	if err != nil {
		return errors.Wrap(err, "opening")
	}

	defer fh.Close()

	buf, err := ioutil.ReadAll(fh)
	if err != nil {
		return errors.Wrap(err, "reading")
	}

	ImportLines(m.manager, buf)

	m.loaded = true

	return nil
}
