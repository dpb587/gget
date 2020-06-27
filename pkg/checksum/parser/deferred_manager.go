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
	manager            checksum.WriteableManager
	expectedAlgorithms checksum.AlgorithmList
	opener             func(context.Context) (io.ReadCloser, error)

	loaded bool
	mutex  sync.RWMutex
}

var _ checksum.Manager = &DeferredManager{}

func NewDeferredManager(manager checksum.WriteableManager, expectedAlgorithms checksum.AlgorithmList, opener func(context.Context) (io.ReadCloser, error)) checksum.Manager {
	return &DeferredManager{
		manager:            manager,
		expectedAlgorithms: expectedAlgorithms,
		opener:             opener,
	}
}

func (m *DeferredManager) GetChecksums(ctx context.Context, resource string, algos checksum.AlgorithmList) (checksum.ChecksumList, error) {
	if len(m.expectedAlgorithms) > 0 {
		if len(m.expectedAlgorithms.Intersection(algos)) == 0 {
			// avoid loading if we don't expect it to be found
			return nil, nil
		}
	}

	err := m.requireLoad(ctx)
	if err != nil {
		return nil, err
	}

	return m.manager.GetChecksums(ctx, resource, algos)
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
