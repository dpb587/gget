package gget

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/dpb587/gget/pkg/downloader"
	"github.com/dpb587/gget/pkg/service"
	"github.com/dpb587/gget/pkg/service/github"
	"github.com/pkg/errors"
	"github.com/tidwall/limiter"
	"github.com/vbauerster/mpb/v4"
)

type ResourceOptions struct {
	Type          service.ResourceType `long:"type" description:"type of resource to get (e.g. asset, archive, blob)" default:"asset"`
	IgnoreMissing []ResourceNameOpt    `long:"ignore-missing" description:"if a resource is not found, skip it rather than failing (glob-friendly)" value-name:"[RESOURCE]" optional:"true" optional-value:"*"`
	Exclude       []ResourceNameOpt    `long:"exclude" description:"exclude resource(s) from download (glob-friendly)" value-name:"RESOURCE"`
}

type DownloadOptions struct {
	ShowRef       bool `long:"show-ref" description:"list matched repository ref instead of downloading"`
	ShowResources bool `long:"show-resources" description:"list matched resources instead of downloading"`

	CD         string            `long:"cd" description:"change to directory before writing files"`
	Executable []ResourceNameOpt `long:"executable" description:"apply executable permissions to downloads" value-name:"[RESOURCE]" optional:"true" optional-value:"*"`
	Stdout     bool              `long:"stdout" description:"write file contents to stdout rather than disk"`
}

type Command struct {
	*Runtime         `group:"Runtime Options"`
	*ResourceOptions `group:"Resource Options"`
	*DownloadOptions `group:"Download Options"`
	Args             CommandArgs `positional-args:"true" required:"true"`
}

type CommandArgs struct {
	Ref       RefOpt            `positional-arg-name:"HOST/OWNER/REPOSITORY[@REF]" description:"release reference"`
	Resources []ResourcePathOpt `positional-arg-name:"[LOCAL-PATH=]RESOURCE" description:"resource name(s) to download (glob-friendly)" optional:"true"`
}

func (c *Command) applySettings() {
	if len(c.Args.Resources) == 0 {
		c.Args.Resources = []ResourcePathOpt{
			{
				RemoteMatch: ResourceNameOpt("*"),
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
	userResourceMatches := make([]bool, len(c.Args.Resources))

	for userResourceIdx, userResource := range c.Args.Resources {
		candidateResources, err := ref.ResolveResource(ctx, c.Type, service.Resource(string(userResource.RemoteMatch)))
		if err != nil {
			return errors.Wrapf(err, "resolving resource %s", string(userResource.RemoteMatch))
		} else if len(candidateResources) == 0 {
			for _, ignoreMissing := range c.IgnoreMissing {
				if ignoreMissing.Match(string(userResource.RemoteMatch)) {
					userResourceMatches[userResourceIdx] = true

					break
				}
			}

			continue
		}

		for _, candidate := range candidateResources {
			{ // is it excluded?
				var excluded bool

				for _, exclude := range c.Exclude {
					if exclude.Match(candidate.GetName()) {
						excluded = true

						break
					}
				}

				if excluded {
					continue
				}
			}

			resolved, matched := userResource.Resolve(candidate.GetName())
			if !matched {
				panic("TODO should always match by now?")
			}

			if _, found := resourceMap[resolved.LocalPath]; found {
				return fmt.Errorf("target file already specified: %s", resolved.LocalPath)
			}

			userResourceMatches[userResourceIdx] = true
			resourceMap[resolved.LocalPath] = candidate
		}
	}

	{ // finally, did we find everything the user asked for?
		for userResourceIdx, userResourceMatched := range userResourceMatches {
			if userResourceMatched {
				continue
			}

			return errors.Wrap(fmt.Errorf("no resource matched: %s", c.Args.Resources[userResourceIdx].RemoteMatch), "expected matching resources")
		}
	}

	if c.ShowResources {
		for _, resource := range resourceMap {
			fmt.Println(resource.GetName())
		}

		return nil
	}

	var downloads []*downloader.Workflow

	for localPath, resource := range resourceMap {
		var steps []downloader.Step

		if localPath == "-" {
			steps = append(
				steps,
				&downloader.DownloadWriterInstaller{
					Writer: os.Stdout,
				},
			)
		} else {
			steps = append(
				steps,
				&downloader.DownloadTmpfileInstaller{
					Tmpdir: filepath.Dir(localPath),
				},
			)
		}

		if ds, ok := resource.(downloader.StepProvider); ok {
			extraSteps, err := ds.GetDownloaderSteps(ctx)
			if err != nil {
				return errors.Wrap(err, "getting download steps")
			}

			steps = append(steps, extraSteps...)
		}

		if localPath != "-" {
			for _, ResourceNameOpt := range c.Executable {
				if !ResourceNameOpt.Match(resource.GetName()) {
					continue
				}

				steps = append(
					steps,
					&downloader.DownloadExecutableInstaller{},
				)
			}

			steps = append(
				steps,
				&downloader.DownloadRenameInstaller{
					Target: localPath,
				},
			)
		}

		downloads = append(
			downloads,
			downloader.NewWorkflow(resource, steps...),
		)
	}

	sort.Slice(downloads, func(i, j int) bool {
		return downloads[i].GetSubject() < downloads[j].GetSubject()
	})

	var pbO io.Writer = os.Stderr

	if c.Runtime.Quiet {
		pbO = ioutil.Discard
	}

	pb := mpb.New(mpb.WithWidth(1), mpb.WithOutput(pbO))

	for _, d := range downloads {
		d.Prepare(pb)
	}

	l := limiter.New(c.Runtime.Parallel)

	for _, d := range downloads {
		d := d
		go func() {
			l.Begin()
			defer l.End()

			err := d.Execute(ctx)
			if err != nil {
				// TODO concurrency
				panic(errors.Wrapf(err, "downloading %s", d.GetSubject()))
			}
		}()
	}

	pb.Wait()

	return nil
}
