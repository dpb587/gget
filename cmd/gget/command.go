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
	"github.com/dpb587/gget/pkg/cli/opt"
	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/github"
	"github.com/dpb587/gget/pkg/service/gitlab"
	"github.com/dpb587/gget/pkg/transfer"
	"github.com/dpb587/gget/pkg/transfer/transferutil"
	"github.com/pkg/errors"
)

type RepositoryOptions struct {
	RefStability []string           `long:"ref-stability" description:"acceptable stability level(s) for latest (values: stable, pre-release, any) (default: stable)" value-name:"STABILITY"`
	RefVersions  opt.ConstraintList `long:"ref-version" description:"version constraint(s) to require of latest (e.g. 4.x)" value-name:"CONSTRAINT"`
	Service      string             `long:"service" description:"specific git service to use (values: github, gitlab) (default: auto-detect)" value-name:"NAME"`

	// TODO(1.x) remove
	ShowRef bool `long:"show-ref" description:"show resolved repository ref instead of downloading" hidden:"true"`
}

type ResourceOptions struct {
	Exclude       opt.ResourceMatcherList `long:"exclude" description:"exclude resource(s) from download (multiple)" value-name:"RESOURCE-GLOB"`
	IgnoreMissing opt.ResourceMatcherList `long:"ignore-missing" description:"if a resource is not found, skip it rather than failing (multiple)" value-name:"[RESOURCE-GLOB]" optional:"true" optional-value:"*"`
	Type          service.ResourceType    `long:"type" description:"type of resource to get (values: asset, archive, blob)" default:"asset" value-name:"TYPE"`

	// TODO(1.x) remove
	ShowResources bool `long:"show-resources" description:"show matched resources instead of downloading" hidden:"true"`
}

type DownloadOptions struct {
	CD             string                  `long:"cd" description:"change to directory before writing files" value-name:"DIR"`
	DumpInfo       string                  `long:"dump-info" description:"write details about the download plan to file" value-name:"LOCAL-PATH"`
	Executable     opt.ResourceMatcherList `long:"executable" description:"apply executable permissions to downloads (multiple)" value-name:"[RESOURCE-GLOB]" optional:"true" optional-value:"*"`
	NoDownload     bool                    `long:"no-download" description:"do not perform any downloads"`
	NoProgress     bool                    `long:"no-progress" description:"do not show live-updating progress during downloads"`
	Parallel       int                     `long:"parallel" description:"maximum number of parallel downloads" default:"3" value-name:"INT"`
	Stdout         bool                    `long:"stdout" description:"write file contents to stdout rather than disk"`
	VerifyChecksum opt.VerifyChecksum      `long:"verify-checksum" description:"strategy for verifying checksums (values: auto, required, none, {algo}, {algo}-min)" value-name:"[METHOD]" default:"auto" optional-value:"required"`
}

type Command struct {
	*Runtime           `group:"Runtime Options"`
	*RepositoryOptions `group:"Repository Options"`
	*ResourceOptions   `group:"Resource Options"`
	*DownloadOptions   `group:"Download Options"`
	Args               CommandArgs `positional-args:"true"`
}

type CommandArgs struct {
	Ref       opt.Ref                  `positional-arg-name:"HOST/OWNER/REPOSITORY[@REF]" description:"repository reference"`
	Resources opt.ResourceTransferList `positional-arg-name:"[LOCAL-PATH=]RESOURCE-GLOB" description:"resource name(s) to download" optional:"true"`
}

func (c *Command) applySettings() {
	if v := c.Service; v != "" {
		c.Args.Ref.Service = v
	}

	if len(c.Args.Resources) == 0 {
		c.Args.Resources = opt.ResourceTransferList{
			{
				RemoteMatch: opt.ResourceMatcher("*"),
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

	{ // TODO(1.x) remove
		if c.ShowRef {
			c.NoDownload = true
		}

		if c.ShowResources {
			c.NoDownload = true
		}
	}
}

func (c *Command) RefResolver(ref service.Ref) (service.RefResolver, error) {
	var resolvers []service.ConditionalRefResolver

	if ref.Service == "" || ref.Service == github.ServiceName {
		resolvers = append(
			resolvers,
			github.NewService(c.Runtime.Logger(), github.NewClientFactory(c.Runtime.Logger(), c.Runtime.NewHTTPClient)),
		)
	}

	if ref.Service == "" || ref.Service == gitlab.ServiceName {
		resolvers = append(
			resolvers,
			gitlab.NewService(c.Runtime.Logger(), gitlab.NewClientFactory(c.Runtime.Logger(), c.Runtime.NewHTTPClient)),
		)
	}

	switch len(resolvers) {
	case 0:
		return nil, fmt.Errorf("unsupported service: %s", ref.Service)
	case 1:
		return resolvers[0], nil
	}

	return service.NewMultiRefResolver(resolvers...), nil
}

func (c *Command) Execute(_ []string) error {
	c.applySettings()

	if c.Args.Ref.Repository == "" {
		return fmt.Errorf("missing argument: repository reference")
	}

	refResolver, err := c.RefResolver(service.Ref(c.Args.Ref))
	if err != nil {
		return errors.Wrap(err, "getting ref resolver")
	}

	ctx := context.Background()

	ref, err := refResolver.ResolveRef(ctx, service.LookupRef{
		Ref:          service.Ref(c.Args.Ref),
		RefVersions:  c.RefVersions.Constraints(),
		RefStability: c.RefStability,
	})
	if err != nil {
		return errors.Wrap(err, "resolving ref")
	}

	{ // TODO(1.x) remove
		if c.ShowRef {
			metadata, err := ref.GetMetadata(ctx)
			if err != nil {
				return errors.Wrap(err, "getting metadata")
			}

			for _, metadatum := range metadata {
				fmt.Printf("%s\t%s\n", metadatum.Name, metadatum.Value)
			}
		}
	}

	resourceMap := map[string]service.ResolvedResource{}

	for _, userResource := range c.Args.Resources {
		candidateResources, err := ref.ResolveResource(ctx, c.Type, service.ResourceName(string(userResource.RemoteMatch)))
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

	{ // TODO(1.x) remove
		if c.ShowResources {
			var results []string

			for _, resource := range resourceMap {
				results = append(results, resource.GetName())
			}

			sort.Strings(results)

			for _, result := range results {
				fmt.Println(result)
			}
		}
	}

	if c.DumpInfo != "" {
		var infoW io.Writer
		var infoC io.Closer

		if c.DumpInfo == "-" {
			infoW = os.Stdout
		} else {
			fh, err := os.OpenFile(c.DumpInfo, os.O_WRONLY|os.O_CREATE, 0700)
			if err != nil {
				return errors.Wrap(err, "opening info file")
			}

			infoW = fh
			infoC = fh
		}

		{ // origin
			cref := ref.CanonicalRef()

			fmt.Fprintf(infoW, "origin\tresolved\t%s\n", cref)
			fmt.Fprintf(infoW, "origin\tservice\t%s\n", cref.Service)
			fmt.Fprintf(infoW, "origin\tserver\t%s\n", cref.Server)
			fmt.Fprintf(infoW, "origin\towner\t%s\n", cref.Owner)
			fmt.Fprintf(infoW, "origin\trepository\t%s\n", cref.Repository)
			fmt.Fprintf(infoW, "origin\tref\t%s\n", cref.Ref)
		}

		{ // metadata
			metadata, err := ref.GetMetadata(ctx)
			if err != nil {
				errors.Wrap(err, "getting metadata")
			}

			for _, metadatum := range metadata {
				fmt.Fprintf(infoW, "metadata\t%s\t%s\n", metadatum.Name, metadatum.Value)
			}
		}

		{ // resources
			var results []string

			for _, resource := range resourceMap {
				results = append(results, resource.GetName())
			}

			sort.Strings(results)

			for _, result := range results {
				fmt.Fprintf(infoW, "resource\t%s\n", result)
			}
		}

		if infoC != nil {
			err = infoC.Close()
			if err != nil {
				return errors.Wrap(err, "closing info")
			}
		}
	}

	if c.NoDownload {
		return nil
	}

	// output = stderr since everything should be progress reports
	stdout := os.Stderr
	var finalStatus io.Writer

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

		if c.NoProgress {
			finalStatus = stdout
		}
	}

	var transfers []*transfer.Transfer

	for localPath, resource := range resourceMap {
		xfer, err := transferutil.BuildTransfer(
			ctx,
			resource,
			localPath,
			transferutil.TransferOptions{
				Executable:                   !c.Executable.Match(resource.GetName()).IsEmpty(),
				ChecksumMode:                 c.VerifyChecksum.Mode(),
				ChecksumAcceptableAlgorithms: c.VerifyChecksum.AcceptableAlgorithms(),
				FinalStatus:                  finalStatus,
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

	if c.Runtime.Quiet || c.NoProgress {
		pbO = ioutil.Discard
	}

	batch := transfer.NewBatch(transfers, c.Parallel, pbO)

	return batch.Transfer(ctx)
}
