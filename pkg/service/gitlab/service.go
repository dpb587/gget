package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/dpb587/gget/pkg/gitutil"
	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/gitlab/gitlabutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

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
	return lookupRef.Ref.Server == "gitlab.com"
}

func (s Service) ResolveRef(ctx context.Context, lookupRef service.LookupRef) (service.ResolvedRef, error) {
	client, err := s.clientFactory.Get(ctx, lookupRef)
	if err != nil {
		return nil, errors.Wrap(err, "building client")
	}

	canonicalRef := lookupRef.Ref

	var cachedRelease *gitlab.Release

	if canonicalRef.Ref == "" {
		release, err := s.resolveLatest(ctx, client, lookupRef)
		if err != nil {
			return nil, errors.Wrap(err, "resolving latest")
		}

		canonicalRef.Ref = release.TagName
		cachedRelease = release
	}

	idPath := gitlabutil.GetRepositoryID(canonicalRef)

	{ // tag
		tag, resp, err := client.Tags.GetTag(idPath, canonicalRef.Ref)
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting tag resolution")
		} else if tag != nil {
			return s.resolveTagReference(ctx, client, canonicalRef, tag, cachedRelease)
		}
	}

	{ // head
		branch, resp, err := client.Branches.GetBranch(idPath, canonicalRef.Ref)
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting branch resolution")
		} else if branch != nil {
			return s.resolveHeadReference(ctx, client, canonicalRef, branch)
		}
	}

	if gitutil.PotentialCommitRE.MatchString(canonicalRef.Ref) { // commit
		commit, resp, err := client.Commits.GetCommit(idPath, canonicalRef.Ref)
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting commit resolution")
		} else {
			canonicalRef.Ref = commit.ID

			return s.resolveCommitReference(ctx, client, canonicalRef, commit.ID)
		}
	}

	return nil, fmt.Errorf("unable to resolve as tag, branch, nor commit: %s", canonicalRef.Ref)
}

func (s Service) resolveLatest(ctx context.Context, client *gitlab.Client, lookupRef service.LookupRef) (*gitlab.Release, error) {
	idPath := gitlabutil.GetRepositoryID(lookupRef.Ref)

	opts := gitlab.ListReleasesOptions{
		PerPage: 25,
	}

	for {
		releases, resp, err := client.Releases.ListReleases(idPath, &opts)
		if err != nil {
			return nil, errors.Wrap(err, "getting releases")
		} else if resp.StatusCode == http.StatusNotFound {
			return nil, errors.New("repository not found")
		}

		for _, release := range releases {
			if !lookupRef.SatisfiesStability("stable") {
				continue
			}

			tagName := release.TagName
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

	if lookupRef.IsComplexRef() {
		return nil, fmt.Errorf("failed to find release matching constraints: %s", strings.Join(lookupRef.ComplexRefModes(), ", "))
	}

	return nil, errors.New("no latest release found")
}

func (s Service) resolveCommitReference(ctx context.Context, client *gitlab.Client, ref service.Ref, commitSHA string) (service.ResolvedRef, error) {
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

func (s Service) resolveHeadReference(ctx context.Context, client *gitlab.Client, ref service.Ref, headRef *gitlab.Branch) (service.ResolvedRef, error) {
	branchName := headRef.Name
	commitSHA := headRef.Commit.ID

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

func (s Service) resolveTagReference(ctx context.Context, client *gitlab.Client, ref service.Ref, tagRef *gitlab.Tag, cachedRelease *gitlab.Release) (service.ResolvedRef, error) {
	tagName := tagRef.Name
	commitSHA := tagRef.Commit.ID

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

	var release *gitlab.Release

	if cachedRelease != nil && cachedRelease.TagName == tagName {
		release = cachedRelease
	} else {
		var resp *gitlab.Response
		var err error

		release, resp, err = client.Releases.GetRelease(gitlabutil.GetRepositoryID(ref), tagName)
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
			// checksumManager: NewReleaseChecksumManager(client, ref.Owner, ref.Repository, release), // TODO
		}
	}

	return res, nil
}
