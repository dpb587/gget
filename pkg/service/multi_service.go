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

	IsRefSupported(context.Context, LookupRef) bool
}

func (rr MultiRefResolver) ResolveRef(ctx context.Context, lookupRef LookupRef) (ResolvedRef, error) {
	for _, resolver := range rr.resolvers {
		if !resolver.IsRefSupported(ctx, lookupRef) {
			continue
		}

		return resolver.ResolveRef(ctx, lookupRef)
	}

	return nil, fmt.Errorf("failed to detect service: %s", lookupRef)
}
