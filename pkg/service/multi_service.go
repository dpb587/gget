package service

import (
	"context"
	"fmt"
)

type MultiRefResolver struct {
	resolvers []ConditionalRefResolver
}

var _ RefResolver = MultiRefResolver{}

func NewMultiRefResolver(resolvers ...ConditionalRefResolver) RefResolver {
	return MultiRefResolver{
		resolvers: resolvers,
	}
}

type ConditionalRefResolver interface {
	RefResolver

	IsRefSupported(context.Context, Ref) bool
}

func (rr MultiRefResolver) ResolveRef(ctx context.Context, ref Ref) (ResolvedRef, error) {
	for _, resolver := range rr.resolvers {
		if !resolver.IsRefSupported(ctx, ref) {
			continue
		}

		return resolver.ResolveRef(ctx, ref)
	}

	return nil, fmt.Errorf("failed to detect service: %s", ref)
}
