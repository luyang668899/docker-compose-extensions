/*
   Copyright 2023 Docker Compose CLI authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package compose

import (
	"context"
	"fmt"

	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"

	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

type quickOptions struct {
	*ProjectOptions
	build    bool
	pull     bool
	detach   bool
	services []string
}

func quickCommand(p *ProjectOptions, dockerCli command.Cli, backendOptions *BackendOptions) *cobra.Command {
	opts := quickOptions{
		ProjectOptions: p,
		build:          true,
		pull:           true,
		detach:         true,
	}

	cmd := &cobra.Command{
		Use:   "quick [OPTIONS] [SERVICE...]",
		Short: "Quick setup and start services with minimal steps",
		Long: `EXPERIMENTAL - Quick setup and start services with minimal steps.

This command combines multiple operations into one:
1. Pull latest images (if needed)
2. Build services (if needed)
3. Start services in detached mode
4. Show status and endpoints
`,
		RunE: Adapt(func(ctx context.Context, args []string) error {
			opts.services = args
			return runQuick(ctx, dockerCli, backendOptions, &opts)
		}),
	}

	cmd.Flags().BoolVar(&opts.build, "no-build", false, "Skip build step")
	cmd.Flags().BoolVar(&opts.pull, "no-pull", false, "Skip pull step")
	cmd.Flags().BoolVar(&opts.detach, "no-detach", false, "Do not start in detached mode")
	return cmd
}

func runQuick(ctx context.Context, dockerCli command.Cli, backendOptions *BackendOptions, opts *quickOptions) error {
	backend, err := compose.NewComposeService(dockerCli, backendOptions.Options...)
	if err != nil {
		return err
	}

	project, name, err := opts.ToProject(ctx, dockerCli, backend, nil)
	if err != nil {
		return err
	}

	// Step 1: Pull images if needed
	if opts.pull {
		fmt.Println("Pulling latest images...")
		if err := backend.Pull(ctx, project, api.PullOptions{}); err != nil {
			fmt.Printf("Warning: Pull failed: %v\n", err)
			// Continue even if pull fails
		}
	}

	// Step 2: Build services if needed
	if opts.build {
		fmt.Println("Building services...")
		if err := backend.Build(ctx, project, api.BuildOptions{}); err != nil {
			return err
		}
	}

	// Step 3: Start services
	fmt.Println("Starting services...")
	uOptions := api.UpOptions{}
	if err := backend.Up(ctx, project, uOptions); err != nil {
		return err
	}

	// Step 4: Show status and endpoints
	fmt.Println("\nServices status:")
	containers, err := backend.Ps(ctx, project.Name, api.PsOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		fmt.Printf("%s: %s\n", container.Service, container.State)
	}

	// Show ports and endpoints
	fmt.Println("\nEndpoints:")
	for _, service := range project.Services {
		if len(service.Ports) > 0 {
			fmt.Printf("%s:\n", service.Name)
			for _, port := range service.Ports {
				fmt.Printf("  %s:%s -> %s/%s\n", port.HostIP, port.Published, port.Target, port.Protocol)
			}
		}
	}

	fmt.Printf("\nProject %s is ready!\n", name)
	return nil
}
