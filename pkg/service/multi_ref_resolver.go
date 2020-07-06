package service

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type MultiRefResolver struct {
	log       *logrus.Logger
	resolvers []ConditionalRefResolver
}

var _ RefResolver = MultiRefResolver{}

func NewMultiRefResolver(log *logrus.Logger, resolvers ...ConditionalRefResolver) RefResolver {
	return MultiRefResolver{
		log:       log,
		resolvers: resolvers,
	}
}

type ConditionalRefResolver interface {
	RefResolver

	ServiceName() string
	IsKnownServer(context.Context, LookupRef) bool
	IsDetectedServer(context.Context, LookupRef) bool
}

func (rr MultiRefResolver) ResolveRef(ctx context.Context, lookupRef LookupRef) (ResolvedRef, error) {
	if serviceName := lookupRef.Service; serviceName != "" {
		for _, resolver := range rr.resolvers {
			if resolver.ServiceName() != serviceName {
				continue
			}

			rr.log.Infof("using service based on ref: %s", resolver.ServiceName())

			return resolver.ResolveRef(ctx, lookupRef)
		}

		return nil, fmt.Errorf("service not recognized: %s", serviceName)
	}

	for _, resolver := range rr.resolvers {
		if !resolver.IsKnownServer(ctx, lookupRef) {
			continue
		}

		rr.log.Infof("using service based on known servers: %s", resolver.ServiceName())

		return resolver.ResolveRef(ctx, lookupRef)
	}

	rr.log.Debugf("attempting ref server detection (ref server not known)")

	for _, resolver := range rr.resolvers {
		if !resolver.IsDetectedServer(ctx, lookupRef) {
			continue
		}

		rr.log.Infof("using service based on server detection: %s", resolver.ServiceName())

		return resolver.ResolveRef(ctx, lookupRef)
	}

	return nil, fmt.Errorf("failed to find service for ref: %s", lookupRef)
}
