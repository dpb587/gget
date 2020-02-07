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

func (s Service) ResolveRef(ctx context.Context, ref service.Ref) (service.ResolvedRef, error) {
	client, err := s.clientFactory.Get(ctx, ref)
	if err != nil {
		return nil, errors.Wrap(err, "building client")
	}

	repo, _, err := client.Repositories.Get(ctx, ref.Owner, ref.Repository)
	if err != nil {
		return nil, errors.Wrap(err, "getting repo")
	}

	if ref.Ref == "" {
		release, _, err := client.Repositories.GetLatestRelease(ctx, repo.GetOwner().GetLogin(), repo.GetName())
		if err != nil {
			return nil, errors.Wrap(err, "getting latest release")
		}

		// TODO avoid dup api call
		ref.Ref = release.GetTagName()
		if ref.Ref == "" {
			panic("logical inconsistency")
		}

		return s.ResolveRef(ctx, ref)
	}

	{ // tag
		gitref, resp, err := client.Git.GetRefs(ctx, repo.GetOwner().GetLogin(), repo.GetName(), path.Join("tags", ref.Ref))
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting tag resolution")
		} else if len(gitref) == 1 {
			return s.resolveTagReference(ctx, client, repo, gitref[0])
		}
	}

	{ // head
		gitref, resp, err := client.Git.GetRefs(ctx, repo.GetOwner().GetLogin(), repo.GetName(), path.Join("heads", ref.Ref))
		if resp.StatusCode == http.StatusNotFound {
			// oh well
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting branch resolution")
		} else if len(gitref) == 1 {
			return s.resolveHeadReference(ctx, client, repo, gitref[0])
		}
	}

	if gitutil.PotentialCommitRE.MatchString(ref.Ref) { // commit
		// client.Git.GetCommit does not resolve partial commits
		commitref, resp, err := client.Repositories.GetCommit(ctx, repo.GetOwner().GetLogin(), repo.GetName(), ref.Ref)
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil
		} else if err != nil {
			return nil, errors.Wrap(err, "attempting commit resolution")
		}

		return s.resolveCommitReference(ctx, client, repo, commitref.GetSHA())
	}

	return nil, fmt.Errorf("unable to resolve as tag, branch, nor commit: %s", ref.Ref)
}

func (s Service) resolveCommitReference(ctx context.Context, client *github.Client, repo *github.Repository, commitSHA string) (service.ResolvedRef, error) {
	res := &CommitRef{
		client:          client,
		repo:            repo,
		commit:          commitSHA,
		archiveFileBase: fmt.Sprintf("%s-%s", repo.GetName(), commitSHA[0:9]),
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

func (s Service) resolveHeadReference(ctx context.Context, client *github.Client, repo *github.Repository, headRef *github.Reference) (service.ResolvedRef, error) {
	branchName := strings.TrimPrefix(headRef.GetRef(), "refs/heads/")

	res := &CommitRef{
		client:          client,
		repo:            repo,
		commit:          headRef.Object.GetSHA(),
		archiveFileBase: fmt.Sprintf("%s-%s", repo.GetName(), path.Base(branchName)),
		RefMetadataService: service.RefMetadataService{
			Metadata: []service.RefMetadata{
				{
					Name:  "branch",
					Value: branchName,
				},
				{
					Name:  "commit",
					Value: headRef.Object.GetSHA(),
				},
			},
		},
	}

	return res, nil
}

func (s Service) resolveTagReference(ctx context.Context, client *github.Client, repo *github.Repository, tagRef *github.Reference) (service.ResolvedRef, error) {
	var tagObj *github.Tag

	if tagRef.Object.GetType() == "tag" { // annotated tag
		var err error

		tagObj, _, err = client.Git.GetTag(ctx, repo.GetOwner().GetLogin(), repo.GetName(), tagRef.Object.GetSHA())
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

	var res service.ResolvedRef = &CommitRef{
		client:          client,
		repo:            repo,
		commit:          tagObj.Object.GetSHA(),
		archiveFileBase: fmt.Sprintf("%s-%s", repo.GetName(), tagName),
		RefMetadataService: service.RefMetadataService{
			Metadata: []service.RefMetadata{
				{
					Name:  "tag",
					Value: tagName,
				},
				{
					Name:  "commit",
					Value: tagObj.Object.GetSHA(),
				},
			},
		},
	}

	release, resp, err := client.Repositories.GetReleaseByTag(ctx, repo.GetOwner().GetLogin(), repo.GetName(), tagName)
	if resp.StatusCode == http.StatusNotFound {
		// oh well
	} else if err != nil {
		return nil, errors.Wrap(err, "getting release by tag")
	} else if release != nil {
		res = &ReleaseRef{
			client:          client,
			repo:            repo,
			release:         release,
			targetRef:       res,
			checksumManager: NewReleaseChecksumManager(client, release),
		}
	}

	return res, nil
}
