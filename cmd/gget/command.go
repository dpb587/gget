package gget

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/github"
	"github.com/dpb587/gget/pkg/transfer"
	"github.com/dpb587/gget/pkg/transfer/transferutil"
	"github.com/pkg/errors"
)

type ResourceOptions struct {
	Type          service.ResourceType `long:"type" description:"type of resource to get (e.g. asset, archive, blob)" default:"asset"`
	IgnoreMissing ResourceMatchers     `long:"ignore-missing" description:"if a resource is not found, skip it rather than failing (glob-friendly; multiple)" value-name:"[RESOURCE]" optional:"true" optional-value:"*"`
	Exclude       ResourceMatchers     `long:"exclude" description:"exclude resource(s) from download (glob-friendly; multiple)" value-name:"RESOURCE"`
}

type DownloadOptions struct {
	ShowRef       bool `long:"show-ref" description:"list matched repository ref instead of downloading"`
	ShowResources bool `long:"show-resources" description:"list matched resources instead of downloading"`

	CD         string           `long:"cd" description:"change to directory before writing files"`
	Executable ResourceMatchers `long:"executable" description:"apply executable permissions to downloads (glob-friendly; multiple)" value-name:"[RESOURCE]" optional:"true" optional-value:"*"`
	Stdout     bool             `long:"stdout" description:"write file contents to stdout rather than disk"`
}

type Command struct {
	*Runtime         `group:"Runtime Options"`
	*ResourceOptions `group:"Resource Options"`
	*DownloadOptions `group:"Download Options"`
	Args             CommandArgs `positional-args:"true" required:"true"`
}

type CommandArgs struct {
	Ref       RefOpt                 `positional-arg-name:"HOST/OWNER/REPOSITORY[@REF]" description:"release reference"`
	Resources []ResourceTransferSpec `positional-arg-name:"[LOCAL-PATH=]RESOURCE" description:"resource name(s) to download (glob-friendly)" optional:"true"`
}

func (c *Command) applySettings() {
	if len(c.Args.Resources) == 0 {
		c.Args.Resources = []ResourceTransferSpec{
			{
				RemoteMatch: ResourceMatcher("*"),
			},
		}
	}

	if c.Stdout {
		for resourceIdx, resource := range c.Args.Resources {
			if resource.LocalPath != "" {
				continue
			}

			c.Args.Resources[resourceIdx].LocalPath = "-"
		}
	}

	if c.CD != "" {
		for resourceIdx, resource := range c.Args.Resources {
			if resource.LocalPath == "-" {
				continue
			}

			c.Args.Resources[resourceIdx].LocalPath = filepath.Join(c.CD, resource.LocalPath)
		}
	}
}

func (c *Command) Execute(_ []string) error {
	c.applySettings()

	if c.Args.Ref.Repository == "" {
		return fmt.Errorf("missing argument: repository")
	}

	ctx := context.Background()
	svc := github.NewService(c.Runtime.Logger(), &github.ClientFactory{RoundTripFactory: c.Runtime.RoundTripLogger})

	ref, err := svc.ResolveRef(ctx, service.Ref(c.Args.Ref))
	if err != nil {
		return errors.Wrap(err, "resolving ref")
	}

	if c.ShowRef {
		for _, metadata := range ref.GetMetadata() {
			fmt.Printf("%s\t%s\n", metadata.Name, metadata.Value)
		}

		if !c.ShowResources {
			// exit early unless they also want to see resources
			return nil
		}
	}

	resourceMap := map[string]service.ResolvedResource{}

	for _, userResource := range c.Args.Resources {
		candidateResources, err := ref.ResolveResource(ctx, c.Type, service.Resource(string(userResource.RemoteMatch)))
		if err != nil {
			return errors.Wrapf(err, "resolving resource %s", string(userResource.RemoteMatch))
		} else if len(candidateResources) == 0 {
			if !c.IgnoreMissing.Match(string(userResource.RemoteMatch)).IsEmpty() {
				continue
			}

			return errors.Wrap(fmt.Errorf("no resource matched: %s", userResource.RemoteMatch), "expected matching resources")
		}

		for _, candidate := range candidateResources {
			if !c.Exclude.Match(candidate.GetName()).IsEmpty() {
				continue
			}

			resolved, matched := userResource.Resolve(candidate.GetName())
			if !matched {
				panic("TODO should always match by now?")
			}

			if _, found := resourceMap[resolved.LocalPath]; found {
				return fmt.Errorf("target file already specified: %s", resolved.LocalPath)
			}

			resourceMap[resolved.LocalPath] = candidate
		}
	}

	if c.ShowResources {
		for _, resource := range resourceMap {
			fmt.Println(resource.GetName())
		}

		return nil
	}

	var transfers []*transfer.Transfer

	for localPath, resource := range resourceMap {
		xfer, err := transferutil.BuildTransfer(
			ctx,
			resource,
			localPath,
			transferutil.TransferOptions{
				Executable: !c.Executable.Match(resource.GetName()).IsEmpty(),
			},
		)
		if err != nil {
			return errors.Wrapf(err, "preparing transfer of %s", resource.GetName())
		}

		transfers = append(transfers, xfer)
	}

	sort.Slice(transfers, func(i, j int) bool {
		// TODO first order by user arg order
		return transfers[i].GetSubject() < transfers[j].GetSubject()
	})

	var pbO io.Writer = os.Stderr

	if c.Runtime.Quiet {
		pbO = ioutil.Discard
	}

	batch := transfer.NewBatch(transfers, c.Runtime.Parallel, pbO)

	return batch.Transfer(ctx)
}
