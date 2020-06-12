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

func (s Service) IsRefSupported(_ context.Context, ref service.Ref) bool {
	return ref.Server == "github.com"
}

func (s Service) ResolveRef(ctx context.Context, ref service.Ref) (service.ResolvedRef, error) {
	client, err := s.clientFactory.Get(ctx, ref)
	if err != nil {
		return nil, errors.Wrap(err, "building client")
	}

	var cachedRelease *github.RepositoryRelease

	if ref.Ref == "" {
		release, resp, err := client.Repositories.GetLatestRelease(ctx, ref.Owner, ref.Repository)
		if err != nil {
			return nil, errors.Wrap(err, "getting latest release")
		} else if resp.StatusCode == http.StatusNotFound {
			return nil, errors.New("repository not found")
		}

		ref.Ref = release.GetTagName()
		cachedRelease = release
	}

	{ // tag
		gitref, resp, err := client.Git.GetRefs(ctx, ref.Owner, ref.Repository, path.Join("tags", ref.Ref))
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting tag resolution")
		} else if len(gitref) == 1 {
			return s.resolveTagReference(ctx, client, ref, gitref[0], cachedRelease)
		}
	}

	{ // head
		gitref, resp, err := client.Git.GetRefs(ctx, ref.Owner, ref.Repository, path.Join("heads", ref.Ref))
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting branch resolution")
		} else if len(gitref) == 1 {
			return s.resolveHeadReference(ctx, client, ref, gitref[0])
		}
	}

	if gitutil.PotentialCommitRE.MatchString(ref.Ref) { // commit
		// client.Git.GetCommit does not resolve partial commits
		commitref, resp, err := client.Repositories.GetCommit(ctx, ref.Owner, ref.Repository, ref.Ref)
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting commit resolution")
		} else {
			ref.Ref = commitref.GetSHA()

			return s.resolveCommitReference(ctx, client, ref, commitref.GetSHA())
		}
	}

	return nil, fmt.Errorf("unable to resolve as tag, branch, nor commit: %s", ref.Ref)
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
			client:          client,
			ref:             ref,
			release:         release,
			targetRef:       res,
			checksumManager: NewReleaseChecksumManager(client, ref.Owner, ref.Repository, release),
		}
	}

	return res, nil
}
