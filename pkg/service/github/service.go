package github

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/dpb587/gget/pkg/gitutil"
	"github.com/dpb587/gget/pkg/service"
	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const ServiceName = "github"

type Service struct {
	log           *logrus.Logger
	clientFactory *ClientFactory
}

func NewService(log *logrus.Logger, clientFactory *ClientFactory) *Service {
	return &Service{
		log:           log,
		clientFactory: clientFactory,
	}
}

var _ service.RefResolver = &Service{}
var _ service.ConditionalRefResolver = &Service{}

func (s Service) IsRefSupported(_ context.Context, lookupRef service.LookupRef) bool {
	return lookupRef.Ref.Service == ServiceName || lookupRef.Ref.Server == "github.com"
}

func (s Service) ResolveRef(ctx context.Context, lookupRef service.LookupRef) (service.ResolvedRef, error) {
	client, err := s.clientFactory.Get(ctx, lookupRef)
	if err != nil {
		return nil, errors.Wrap(err, "building client")
	}

	ref := lookupRef.Ref
	ref.Service = ServiceName

	rr := &refResolver{
		client:       client,
		lookupRef:    lookupRef,
		canonicalRef: ref,
	}

	if ref.Ref == "" {
		release, err := s.resolveLatest(ctx, client, lookupRef)
		if err != nil {
			return nil, errors.Wrap(err, "resolving latest")
		}

		rr.canonicalRef.Ref = release.GetTagName()

		return &ReleaseRef{
			refResolver: rr,
			release:     release,
		}, nil
	}

	{ // tag
		gitref, resp, err := client.Git.GetRefs(ctx, rr.canonicalRef.Owner, rr.canonicalRef.Repository, path.Join("tags", rr.canonicalRef.Ref))
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting tag resolution")
		} else if len(gitref) == 1 {
			return rr.resolveTag(ctx, gitref[0], true)
		}
	}

	{ // head
		gitref, resp, err := client.Git.GetRefs(ctx, rr.canonicalRef.Owner, rr.canonicalRef.Repository, path.Join("heads", rr.canonicalRef.Ref))
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting branch resolution")
		} else if len(gitref) == 1 {
			return rr.resolveHead(ctx, gitref[0])
		}
	}

	if gitutil.PotentialCommitRE.MatchString(rr.canonicalRef.Ref) { // commit
		// client.Git.GetCommit does not resolve partial commits
		commitref, resp, err := client.Repositories.GetCommit(ctx, rr.canonicalRef.Owner, rr.canonicalRef.Repository, rr.canonicalRef.Ref)
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting commit resolution")
		} else {
			rr.canonicalRef.Ref = commitref.GetSHA()

			return rr.resolveCommit(ctx, commitref.GetSHA())
		}
	}

	return nil, fmt.Errorf("unable to resolve as tag, branch, nor commit: %s", rr.canonicalRef.Ref)
}

func (s Service) resolveLatest(ctx context.Context, client *github.Client, lookupRef service.LookupRef) (*github.RepositoryRelease, error) {
	if lookupRef.IsComplexRef() {
		opts := github.ListOptions{
			PerPage: 25,
		}

		for {
			releases, resp, err := client.Repositories.ListReleases(ctx, lookupRef.Ref.Owner, lookupRef.Ref.Repository, &opts)
			if err != nil {
				return nil, errors.Wrap(err, "iterating releases")
			} else if resp.StatusCode == http.StatusNotFound {
				return nil, errors.New("repository not found")
			}

			for _, release := range releases {
				{
					var stability = "stable"

					if release.GetPrerelease() {
						stability = "pre-release"
					}

					if !lookupRef.SatisfiesStability(stability) {
						continue
					}
				}

				tagName := release.GetTagName()
				match, err := lookupRef.SatisfiesVersion(tagName)
				if err != nil {
					s.log.Debugf("skipping invalid semver tag: %s", tagName)

					continue
				} else if !match {
					continue
				}

				return release, nil
			}

			opts.Page = resp.NextPage

			if opts.Page == 0 {
				break
			}
		}

		return nil, fmt.Errorf("failed to find release matching constraints: %s", strings.Join(lookupRef.ComplexRefModes(), ", "))
	}

	release, resp, err := client.Repositories.GetLatestRelease(ctx, lookupRef.Ref.Owner, lookupRef.Ref.Repository)
	if err != nil {
		return nil, errors.Wrap(err, "getting latest release")
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("repository not found")
	}

	return release, nil
}
