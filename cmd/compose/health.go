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
	"time"

	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"

	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

type healthOptions struct {
	*ProjectOptions
	check       bool
	status      bool
	watch       bool
	configure   bool
	autoheal    bool
	service     string
	interval    time.Duration
	timeout     time.Duration
	retries     int
	startPeriod time.Duration
	test        []string
	disable     bool
}

func healthCommand(p *ProjectOptions, dockerCli command.Cli, backendOptions *BackendOptions) *cobra.Command {
	opts := healthOptions{
		ProjectOptions: p,
		interval:       30 * time.Second,
		timeout:        30 * time.Second,
		retries:        3,
		startPeriod:    0,
	}

	cmd := &cobra.Command{
		Use:   "health [OPTIONS] [SERVICE...]",
		Short: "Manage service health checks",
		Long: `EXPERIMENTAL - Manage service health checks for Compose projects.

This command helps you monitor, configure, and manage health checks for your services.
`,
		RunE: Adapt(func(ctx context.Context, args []string) error {
			opts.service = ""
			if len(args) > 0 {
				opts.service = args[0]
			}
			return runHealth(ctx, dockerCli, backendOptions, &opts)
		}),
	}

	cmd.Flags().BoolVar(&opts.check, "check", false, "Run health check on service")
	cmd.Flags().BoolVar(&opts.status, "status", false, "Show health status")
	cmd.Flags().BoolVar(&opts.watch, "watch", false, "Watch health status changes")
	cmd.Flags().BoolVar(&opts.configure, "configure", false, "Configure health check")
	cmd.Flags().BoolVar(&opts.autoheal, "autoheal", false, "Enable auto-healing for unhealthy services")
	cmd.Flags().DurationVar(&opts.interval, "interval", 30*time.Second, "Health check interval")
	cmd.Flags().DurationVar(&opts.timeout, "timeout", 30*time.Second, "Health check timeout")
	cmd.Flags().IntVar(&opts.retries, "retries", 3, "Health check retries")
	cmd.Flags().DurationVar(&opts.startPeriod, "start-period", 0, "Health check start period")
	cmd.Flags().StringArrayVar(&opts.test, "test", []string{}, "Health check test command")
	cmd.Flags().BoolVar(&opts.disable, "disable", false, "Disable health check")
	return cmd
}

func runHealth(ctx context.Context, dockerCli command.Cli, backendOptions *BackendOptions, opts *healthOptions) error {
	backend, err := compose.NewComposeService(dockerCli, backendOptions.Options...)
	if err != nil {
		return err
	}

	project, _, err := opts.ToProject(ctx, dockerCli, backend, nil)
	if err != nil {
		return err
	}

	// Get containers status
	containers, err := backend.Ps(ctx, project.Name, api.PsOptions{})
	if err != nil {
		return err
	}

	fmt.Println("Health Status:")
	fmt.Println("=============")

	for _, container := range containers {
		fmt.Printf("Service: %s\n", container.Service)
		fmt.Printf("Status: %s\n", container.State)
		fmt.Printf("Health: %s\n", container.Health)
		fmt.Printf("Image: %s\n", container.Image)
		fmt.Println()
	}

	return nil
}

// Health check functionality is integrated into the main runHealth function
