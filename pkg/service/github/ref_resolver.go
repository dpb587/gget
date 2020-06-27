package github

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/dpb587/gget/pkg/service"
	"github.com/google/go-github/v29/github"
	"github.com/pkg/errors"
)

type refResolver struct {
	client       *github.Client
	lookupRef    service.LookupRef
	canonicalRef service.Ref
}

func (rr *refResolver) resolveTagWithRelease(ctx context.Context, release *github.RepositoryRelease) (service.ResolvedRef, error) {
	return rr.resolveTag(
		ctx,
		&github.Reference{
			Ref: release.TagName,
			Object: &github.GitObject{
				// Type: "commit",
				SHA: release.TargetCommitish,
			},
		},
		false,
	)
}

func (rr *refResolver) resolveCommit(ctx context.Context, commitSHA string) (service.ResolvedRef, error) {
	res := &CommitRef{
		client:          rr.client,
		ref:             rr.canonicalRef,
		commit:          commitSHA,
		archiveFileBase: fmt.Sprintf("%s-%s", rr.canonicalRef.Repository, commitSHA[0:9]),
		metadata: service.RefMetadata{
			{
				Name:  "commit",
				Value: commitSHA,
			},
		},
	}

	return res, nil
}

func (rr *refResolver) resolveHead(ctx context.Context, headRef *github.Reference) (service.ResolvedRef, error) {
	branchName := strings.TrimPrefix(headRef.GetRef(), "refs/heads/")
	commitSHA := headRef.Object.GetSHA()

	res := &CommitRef{
		client:          rr.client,
		ref:             rr.canonicalRef,
		commit:          commitSHA,
		archiveFileBase: fmt.Sprintf("%s-%s", rr.canonicalRef.Repository, path.Base(branchName)),
		metadata: service.RefMetadata{
			{
				Name:  "branch",
				Value: branchName,
			},
			{
				Name:  "commit",
				Value: commitSHA,
			},
		},
	}

	return res, nil
}

func (rr *refResolver) resolveTag(ctx context.Context, tagRef *github.Reference, attemptRelease bool) (service.ResolvedRef, error) {
	var tagObj *github.Tag

	if tagRef.Object.GetType() == "tag" { // annotated tag
		var err error

		tagObj, _, err = rr.client.Git.GetTag(ctx, rr.canonicalRef.Owner, rr.canonicalRef.Repository, tagRef.Object.GetSHA())
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
		client:          rr.client,
		ref:             rr.canonicalRef,
		commit:          commitSHA,
		archiveFileBase: fmt.Sprintf("%s-%s", rr.canonicalRef.Repository, tagName),
		metadata: service.RefMetadata{
			{
				Name:  "tag",
				Value: tagName,
			},
			{
				Name:  "commit",
				Value: commitSHA,
			},
		},
	}

	if !attemptRelease {
		return res, nil
	}

	release, resp, err := rr.client.Repositories.GetReleaseByTag(ctx, rr.canonicalRef.Owner, rr.canonicalRef.Repository, tagName)
	if resp.StatusCode == http.StatusNotFound {
		// oh well
	} else if err != nil {
		return nil, errors.Wrap(err, "getting release by tag")
	} else if release != nil {
		res = &ReleaseRef{
			refResolver: rr,
			release:     release,
			targetRef:   res,
		}
	}

	return res, nil
}
