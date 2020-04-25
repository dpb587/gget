package transfer

import (
	"context"

	"github.com/vbauerster/mpb/v4/decor"
)

type Step interface {
	GetProgressParams() (int64, decor.Decorator)
	Execute(ctx context.Context, s *State) error
}
