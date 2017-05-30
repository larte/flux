package main

import (
	"fmt"
	"os/user"

	"github.com/spf13/cobra"

	"github.com/weaveworks/flux"
	"github.com/weaveworks/flux/update"
)

type serviceReleaseOpts struct {
	*serviceOpts
	services    []string
	allServices bool
	image       string
	allImages   bool
	exclude     []string
	dryRun      bool
	user        string
	message     string
	outputOpts
}

func newServiceRelease(parent *serviceOpts) *serviceReleaseOpts {
	return &serviceReleaseOpts{serviceOpts: parent}
}

func (opts *serviceReleaseOpts) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Release a new version of a service.",
		Example: makeExample(
			"fluxctl release --service=default/foo --update-image=library/hello:v2",
			"fluxctl release --all --update-image=library/hello:v2",
			"fluxctl release --service=default/foo --update-all-images",
		),
		RunE: opts.RunE,
	}

	username := ""
	user, err := user.Current()
	if err == nil {
		username = user.Username
	}

	OutputFlags(cmd, &opts.outputOpts)
	cmd.Flags().StringSliceVarP(&opts.services, "service", "s", []string{}, "service to release")
	cmd.Flags().BoolVar(&opts.allServices, "all", false, "release all services")
	cmd.Flags().StringVarP(&opts.image, "update-image", "i", "", "update a specific image")
	cmd.Flags().BoolVar(&opts.allImages, "update-all-images", false, "update all images to latest versions")
	cmd.Flags().StringSliceVar(&opts.exclude, "exclude", []string{}, "exclude a service")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "do not release anything; just report back what would have been done")
	cmd.Flags().StringVarP(&opts.message, "message", "m", "", "attach a message to the release job")
	cmd.Flags().StringVar(&opts.user, "user", username, "override the user reported as initating the release job")
	return cmd
}

func (opts *serviceReleaseOpts) RunE(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return errorWantedNoArgs
	}

	if err := checkExactlyOne("--update-image=<image> or --update-all-images", opts.image != "", opts.allImages); err != nil {
		return err
	}

	if len(opts.services) <= 0 && !opts.allServices {
		return newUsageError("please supply either --all, or at least one --service=<service>")
	}

	var services []update.ServiceSpec
	if opts.allServices {
		services = []update.ServiceSpec{update.ServiceSpecAll}
	} else {
		for _, service := range opts.services {
			if _, err := flux.ParseServiceID(service); err != nil {
				return err
			}
			services = append(services, update.ServiceSpec(service))
		}
	}

	var (
		image update.ImageSpec
		err   error
	)
	switch {
	case opts.image != "":
		image, err = update.ParseImageSpec(opts.image)
		if err != nil {
			return err
		}
	case opts.allImages:
		image = update.ImageSpecLatest
	}

	var kind update.ReleaseKind = update.ReleaseKindExecute
	if opts.dryRun {
		kind = update.ReleaseKindPlan
	}

	var excludes []flux.ServiceID
	for _, exclude := range opts.exclude {
		s, err := flux.ParseServiceID(exclude)
		if err != nil {
			return err
		}
		excludes = append(excludes, s)
	}

	if opts.dryRun {
		fmt.Fprintf(cmd.OutOrStderr(), "Requesting dry-run release...\n")
	} else {
		fmt.Fprintf(cmd.OutOrStderr(), "Requesting release ...\n")
	}

	jobID, err := opts.API.UpdateImages(noInstanceID, update.ReleaseSpec{
		ServiceSpecs: services,
		ImageSpec:    image,
		Kind:         kind,
		Excludes:     excludes,
		Cause: update.ReleaseCause{
			User:    opts.user,
			Message: opts.message,
		},
	})
	if err != nil {
		return err
	}

	return await(cmd.OutOrStdout(), cmd.OutOrStderr(), opts.API, jobID, !opts.dryRun, opts.verbose)
}
