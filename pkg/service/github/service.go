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

	var cachedRelease *github.RepositoryRelease

	canonicalRef := lookupRef.Ref
	canonicalRef.Service = ServiceName

	if canonicalRef.Ref == "" {
		release, err := s.resolveLatest(ctx, client, lookupRef)
		if err != nil {
			return nil, errors.Wrap(err, "resolving latest")
		}

		canonicalRef.Ref = release.GetTagName()
		cachedRelease = release
	}

	{ // tag
		gitref, resp, err := client.Git.GetRefs(ctx, canonicalRef.Owner, canonicalRef.Repository, path.Join("tags", canonicalRef.Ref))
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting tag resolution")
		} else if len(gitref) == 1 {
			return s.resolveTagReference(ctx, client, canonicalRef, gitref[0], cachedRelease)
		}
	}

	{ // head
		gitref, resp, err := client.Git.GetRefs(ctx, canonicalRef.Owner, canonicalRef.Repository, path.Join("heads", canonicalRef.Ref))
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting branch resolution")
		} else if len(gitref) == 1 {
			return s.resolveHeadReference(ctx, client, canonicalRef, gitref[0])
		}
	}

	if gitutil.PotentialCommitRE.MatchString(canonicalRef.Ref) { // commit
		// client.Git.GetCommit does not resolve partial commits
		commitref, resp, err := client.Repositories.GetCommit(ctx, canonicalRef.Owner, canonicalRef.Repository, canonicalRef.Ref)
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting commit resolution")
		} else {
			canonicalRef.Ref = commitref.GetSHA()

			return s.resolveCommitReference(ctx, client, canonicalRef, commitref.GetSHA())
		}
	}

	return nil, fmt.Errorf("unable to resolve as tag, branch, nor commit: %s", canonicalRef.Ref)
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

func (s Service) resolveCommitReference(ctx context.Context, client *github.Client, ref service.Ref, commitSHA string) (service.ResolvedRef, error) {
	res := &CommitRef{
		client:          client,
		ref:             ref,
		commit:          commitSHA,
		archiveFileBase: fmt.Sprintf("%s-%s", ref.Repository, commitSHA[0:9]),
		RefMetadataService: service.RefMetadataService{
			Metadata: []service.RefMetadata{
				{
					Name:  "commit",
					Value: commitSHA,
				},
			},
		},
	}

	return res, nil
}

func (s Service) resolveHeadReference(ctx context.Context, client *github.Client, ref service.Ref, headRef *github.Reference) (service.ResolvedRef, error) {
	branchName := strings.TrimPrefix(headRef.GetRef(), "refs/heads/")
	commitSHA := headRef.Object.GetSHA()

	res := &CommitRef{
		client:          client,
		ref:             ref,
		commit:          commitSHA,
		archiveFileBase: fmt.Sprintf("%s-%s", ref.Repository, path.Base(branchName)),
		RefMetadataService: service.RefMetadataService{
			Metadata: []service.RefMetadata{
				{
					Name:  "branch",
					Value: branchName,
				},
				{
					Name:  "commit",
					Value: commitSHA,
				},
			},
		},
	}

	return res, nil
}

func (s Service) resolveTagReference(ctx context.Context, client *github.Client, ref service.Ref, tagRef *github.Reference, cachedRelease *github.RepositoryRelease) (service.ResolvedRef, error) {
	var tagObj *github.Tag

	if tagRef.Object.GetType() == "tag" { // annotated tag
		var err error

		tagObj, _, err = client.Git.GetTag(ctx, ref.Owner, ref.Repository, tagRef.Object.GetSHA())
		if err != nil {
			return nil, errors.Wrap(err, "getting tag of annotated tag")
		}
	} else { // lightweight tag
		// mock to save an API call
		tagObj = &github.Tag{
			Tag: tagRef.Ref,
			Object: &github.GitObject{
				SHA: tagRef.Object.SHA,
			},
		}
	}

	tagName := strings.TrimPrefix(tagObj.GetTag(), "refs/tags/")
	commitSHA := tagObj.Object.GetSHA()

	var res service.ResolvedRef = &CommitRef{
		client:          client,
		ref:             ref,
		commit:          commitSHA,
		archiveFileBase: fmt.Sprintf("%s-%s", ref.Repository, tagName),
		RefMetadataService: service.RefMetadataService{
			Metadata: []service.RefMetadata{
				{
					Name:  "tag",
					Value: tagName,
				},
				{
					Name:  "commit",
					Value: commitSHA,
				},
			},
		},
	}

	var release *github.RepositoryRelease

	if cachedRelease != nil && cachedRelease.GetTagName() == tagName {
		release = cachedRelease
	} else {
		var resp *github.Response
		var err error

		release, resp, err = client.Repositories.GetReleaseByTag(ctx, ref.Owner, ref.Repository, tagName)
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "getting release by tag")
		}
	}

	if release != nil {
		res = &ReleaseRef{
			client:    client,
			ref:       ref,
			release:   release,
			targetRef: res,
		}
	}

	return res, nil
}
