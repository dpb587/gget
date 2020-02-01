package github

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/dpb587/ghet/pkg/service"
	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type Service struct {
	client *github.Client
}

func NewService(client *github.Client) *Service {
	return &Service{
		client: client,
	}
}

var _ service.RefResolver = &Service{}

func (s Service) ResolveRef(ctx context.Context, ref service.Ref) (service.ResolvedRef, error) {
	repo, _, err := s.client.Repositories.Get(ctx, ref.Owner, ref.Repository)
	if err != nil {
		return nil, errors.Wrap(err, "getting repo")
	}

	if ref.Ref == "" {
		release, _, err := s.client.Repositories.GetLatestRelease(ctx, repo.GetOwner().GetLogin(), repo.GetName())
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

	// tag
	gitref, resp, err := s.client.Git.GetRefs(ctx, repo.GetOwner().GetLogin(), repo.GetName(), path.Join("tags", ref.Ref))
	if resp.StatusCode == http.StatusNotFound {
		// oh well
	} else if err != nil {
		return nil, errors.Wrap(err, "attempting tag resolution")
	} else if len(gitref) == 1 {
		return s.resolveTagReference(ctx, repo, gitref[0])
	}

	// head
	gitref, _, err = s.client.Git.GetRefs(ctx, repo.GetOwner().GetLogin(), repo.GetName(), path.Join("heads", ref.Ref))
	if err != nil {
		return nil, errors.Wrap(err, "attempting branch resolution")
	} else if len(gitref) == 1 {
		return s.resolveHeadReference(ctx, repo, gitref[0])
	}

	// commit
	commitref, _, err := s.client.Git.GetCommit(ctx, repo.GetOwner().GetLogin(), repo.GetName(), ref.Ref)
	if err != nil {
		return nil, errors.Wrap(err, "attempting commit resolution")
	} else if gitref != nil {
		return s.resolveCommitReference(ctx, repo, commitref)
	}

	return nil, fmt.Errorf("unable to resolve: %s", ref.Ref)
}

func (s Service) resolveCommitReference(ctx context.Context, repo *github.Repository, commitRef *github.Commit) (service.ResolvedRef, error) {
	res := &CommitRef{
		client:          s.client,
		repo:            repo,
		commit:          commitRef.GetSHA(),
		archiveFileBase: fmt.Sprintf("%s-%s", repo.GetName(), commitRef.GetSHA()[0:9]),
		RefMetadataService: service.RefMetadataService{
			Metadata: []service.RefMetadata{
				{
					Name:  "commit",
					Value: commitRef.GetSHA(),
				},
			},
		},
	}

	return res, nil
}

func (s Service) resolveHeadReference(ctx context.Context, repo *github.Repository, headRef *github.Reference) (service.ResolvedRef, error) {
	branchName := strings.TrimPrefix(headRef.GetRef(), "refs/heads/")

	res := &CommitRef{
		client:          s.client,
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

func (s Service) resolveTagReference(ctx context.Context, repo *github.Repository, tagRef *github.Reference) (service.ResolvedRef, error) {
	var tagObj *github.Tag

	if tagRef.Object.GetType() == "tag" {
		// annotated tag
		var err error

		tagObj, _, err = s.client.Git.GetTag(ctx, repo.GetOwner().GetLogin(), repo.GetName(), tagRef.Object.GetSHA())
		if err != nil {
			return nil, errors.Wrap(err, "getting tag of annotated tag")
		}
	} else { // lightweight
		// stub to save an API call
		tagObj = &github.Tag{
			Tag: tagRef.Ref,
			Object: &github.GitObject{
				SHA: tagRef.Object.SHA,
			},
		}
	}

	tagName := strings.TrimPrefix(tagObj.GetTag(), "refs/tags/")

	var res service.ResolvedRef = &CommitRef{
		client:          s.client,
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

	release, resp, err := s.client.Repositories.GetReleaseByTag(ctx, repo.GetOwner().GetLogin(), repo.GetName(), tagName)
	if resp.StatusCode == http.StatusNotFound {
		// oh well
	} else if err != nil {
		return nil, errors.Wrap(err, "getting release by tag")
	} else if release != nil {
		res = &ReleaseRef{
			client:    s.client,
			repo:      repo,
			release:   release,
			targetRef: res,
		}
	}

	return res, nil
}
