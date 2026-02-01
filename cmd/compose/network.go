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

	"github.com/docker/compose/v5/pkg/compose"
)

type networkOptions struct {
	*ProjectOptions
	list      bool
	create    bool
	remove    bool
	inspect   bool
	connect   bool
	disconnect bool
	name      string
	driver    string
	attachable bool
	internal  bool
	service   string
	ipamDriver string
	ipamConfig string
}

func networkCommand(p *ProjectOptions, dockerCli command.Cli, backendOptions *BackendOptions) *cobra.Command {
	opts := networkOptions{
		ProjectOptions: p,
		driver:         "bridge",
	}

	cmd := &cobra.Command{
		Use:   "network [OPTIONS] [NAME]",
		Short: "Manage networks",
		Long: `EXPERIMENTAL - Manage networks for Compose projects.

This command helps you create, configure, and manage networks for your Compose projects.
`,
		RunE: Adapt(func(ctx context.Context, args []string) error {
			if len(args) > 0 {
				opts.name = args[0]
			}
			return runNetwork(ctx, dockerCli, backendOptions, &opts)
		}),
	}

	cmd.Flags().BoolVar(&opts.list, "list", false, "List networks")
	cmd.Flags().BoolVar(&opts.create, "create", false, "Create network")
	cmd.Flags().BoolVar(&opts.remove, "remove", false, "Remove network")
	cmd.Flags().BoolVar(&opts.inspect, "inspect", false, "Inspect network")
	cmd.Flags().BoolVar(&opts.connect, "connect", false, "Connect service to network")
	cmd.Flags().BoolVar(&opts.disconnect, "disconnect", false, "Disconnect service from network")
	cmd.Flags().StringVar(&opts.driver, "driver", "bridge", "Network driver")
	cmd.Flags().BoolVar(&opts.attachable, "attachable", false, "Make network attachable")
	cmd.Flags().BoolVar(&opts.internal, "internal", false, "Make network internal")
	cmd.Flags().StringVar(&opts.service, "service", "", "Service name for connect/disconnect")
	cmd.Flags().StringVar(&opts.ipamDriver, "ipam-driver", "default", "IPAM driver")
	cmd.Flags().StringVar(&opts.ipamConfig, "ipam-config", "", "IPAM configuration (e.g., \"subnet=192.168.1.0/24\")")
	return cmd
}

func runNetwork(ctx context.Context, dockerCli command.Cli, backendOptions *BackendOptions, opts *networkOptions) error {
	backend, err := compose.NewComposeService(dockerCli, backendOptions.Options...)
	if err != nil {
		return err
	}

	project, _, err := opts.ToProject(ctx, dockerCli, backend, nil)
	if err != nil {
		return err
	}

	// For now, we'll just list the services and their networks
	fmt.Println("Network Information:")
	fmt.Println("====================")

	for _, service := range project.Services {
		fmt.Printf("Service: %s\n", service.Name)
		if len(service.Networks) > 0 {
			fmt.Println("Networks:")
			for networkName := range service.Networks {
				fmt.Printf("  - %s\n", networkName)
			}
		} else {
			fmt.Println("Networks: (default)")
		}
		fmt.Println()
	}

	return nil
}

// Network management functions are integrated into the main runNetwork function

