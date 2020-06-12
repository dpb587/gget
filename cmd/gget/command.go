package gget

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"code.cloudfoundry.org/bytefmt"
	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/github"
	"github.com/dpb587/gget/pkg/service/gitlab"
	"github.com/dpb587/gget/pkg/transfer"
	"github.com/dpb587/gget/pkg/transfer/transferutil"
	"github.com/pkg/errors"
)

type RepositoryOptions struct {
	Service string `long:"service" description:"specific git service to use (values: github, gitlab) (default: auto-detect)" value-name:"NAME"`

	ShowRef bool `long:"show-ref" description:"show resolved repository ref instead of downloading"`
}

type ResourceOptions struct {
	Type          service.ResourceType `long:"type" description:"type of resource to get (values: asset, archive, blob)" default:"asset" value-name:"TYPE"`
	IgnoreMissing ResourceMatchers     `long:"ignore-missing" description:"if a resource is not found, skip it rather than failing (multiple)" value-name:"[RESOURCE-GLOB]" optional:"true" optional-value:"*"`
	Exclude       ResourceMatchers     `long:"exclude" description:"exclude resource(s) from download (multiple)" value-name:"RESOURCE-GLOB"`

	ShowResources bool `long:"show-resources" description:"show matched resources instead of downloading"`
}

type DownloadOptions struct {
	CD         string           `long:"cd" description:"change to directory before writing files" value-name:"DIR"`
	Executable ResourceMatchers `long:"executable" description:"apply executable permissions to downloads (multiple)" value-name:"[RESOURCE-GLOB]" optional:"true" optional-value:"*"`
	Stdout     bool             `long:"stdout" description:"write file contents to stdout rather than disk"`
	Parallel   int              `long:"parallel" description:"maximum number of parallel downloads" default:"3" value-name:"INT"`
}

type Command struct {
	*Runtime           `group:"Runtime Options"`
	*RepositoryOptions `group:"Repository Options"`
	*ResourceOptions   `group:"Resource Options"`
	*DownloadOptions   `group:"Download Options"`
	Args               CommandArgs `positional-args:"true" required:"true"`
}

type CommandArgs struct {
	Ref       RefOpt                 `positional-arg-name:"HOST/OWNER/REPOSITORY[@REF]" description:"repository reference"`
	Resources []ResourceTransferSpec `positional-arg-name:"[LOCAL-PATH=]RESOURCE-GLOB" description:"resource name(s) to download" optional:"true"`
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

func (c *Command) RefResolver() (service.RefResolver, error) {
	var resolvers []service.ConditionalRefResolver

	if c.Service == "" || c.Service == "github" {
		resolvers = append(
			resolvers,
			github.NewService(c.Runtime.Logger(), github.NewClientFactory(c.Runtime.Logger(), c.Runtime.RoundTripLogger)),
		)
	}

	if c.Service == "" || c.Service == "gitlab" {
		resolvers = append(
			resolvers,
			gitlab.NewService(c.Runtime.Logger(), gitlab.NewClientFactory(c.Runtime.Logger(), c.Runtime.RoundTripLogger)),
		)
	}

	switch len(resolvers) {
	case 0:
		return nil, fmt.Errorf("unsupported service: %s", c.Service)
	case 1:
		return resolvers[0], nil
	}

	return service.NewMultiRefResolver(resolvers...), nil
}

func (c *Command) Execute(_ []string) error {
	c.applySettings()

	if c.Args.Ref.Repository == "" {
		return fmt.Errorf("missing argument: repository")
	}

	refResolver, err := c.RefResolver()
	if err != nil {
		return errors.Wrap(err, "getting ref resolver")
	}

	ctx := context.Background()

	ref, err := refResolver.ResolveRef(ctx, service.Ref(c.Args.Ref))
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

			return fmt.Errorf("no resource matched: %s", userResource.RemoteMatch)
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
		var results []string

		for _, resource := range resourceMap {
			results = append(results, resource.GetName())
		}

		sort.Strings(results)

		for _, result := range results {
			fmt.Println(result)
		}

		return nil
	}

	// output = stderr since everything should be progress reports
	stdout := os.Stderr

	if !c.Runtime.Quiet {
		l := len(resourceMap)
		ls := ""

		if l != 1 {
			ls = "s"
		}

		var downloadSizeMissing bool
		var downloadSize int64

		for _, resource := range resourceMap {
			size := resource.GetSize()
			if size == 0 {
				downloadSizeMissing = true

				break
			}

			downloadSize += size
		}

		var extra string

		if !downloadSizeMissing {
			extra = fmt.Sprintf(" (%s)", bytefmt.ByteSize(uint64(downloadSize)))
		}

		fmt.Fprintf(stdout, "Downloading %d file%s%s from %s\n", l, ls, extra, ref.CanonicalRef())
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

	var pbO io.Writer = stdout

	if c.Runtime.Quiet {
		pbO = ioutil.Discard
	}

	batch := transfer.NewBatch(transfers, c.Parallel, pbO)

	return batch.Transfer(ctx)
}
