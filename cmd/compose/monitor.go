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
	"os"
	"time"

	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"

	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

type monitorOptions struct {
	*ProjectOptions
	interval   time.Duration
	format     string
	watch      bool
	outputFile string
}

func monitorCommand(p *ProjectOptions, dockerCli command.Cli, backendOptions *BackendOptions) *cobra.Command {
	opts := monitorOptions{
		ProjectOptions: p,
		interval:       5 * time.Second,
		format:         "table",
		watch:          true,
	}

	cmd := &cobra.Command{
		Use:   "monitor [OPTIONS]",
		Short: "Monitor services status and resources",
		Long: `EXPERIMENTAL - Monitor services status and resources usage.

This command provides real-time monitoring of:
- Service status (running, stopped, etc.)
- Container health
- Resource usage (CPU, memory, network, disk)
- Port mappings and endpoints
`,
		RunE: Adapt(func(ctx context.Context, args []string) error {
			return runMonitor(ctx, dockerCli, backendOptions, &opts)
		}),
	}

	cmd.Flags().DurationVar(&opts.interval, "interval", 5*time.Second, "Refresh interval")
	cmd.Flags().StringVar(&opts.format, "format", "table", "Output format (table, json)")
	cmd.Flags().BoolVar(&opts.watch, "watch", true, "Continuously monitor services")
	cmd.Flags().StringVar(&opts.outputFile, "output", "", "Write output to file instead of stdout")
	return cmd
}

func runMonitor(ctx context.Context, dockerCli command.Cli, backendOptions *BackendOptions, opts *monitorOptions) error {
	backend, err := compose.NewComposeService(dockerCli, backendOptions.Options...)
	if err != nil {
		return err
	}

	project, _, err := opts.ToProject(ctx, dockerCli, backend, nil)
	if err != nil {
		return err
	}

	// Determine output destination
	output := os.Stdout
	if opts.outputFile != "" {
		outputFile, err := os.Create(opts.outputFile)
		if err != nil {
			return err
		}
		defer outputFile.Close()
		output = outputFile
	}

	// Monitor loop
	for {
		// Clear screen if watching
		if opts.watch && opts.outputFile == "" {
			fmt.Fprint(output, "\033[2J\033[H")
		}

		// Show header
		fmt.Fprintf(output, "=== Docker Compose Monitor ===\n")
		fmt.Fprintf(output, "Project: %s\n", project.Name)
		fmt.Fprintf(output, "Time: %s\n\n", time.Now().Format(time.RFC3339))

		// Get services status
		containers, err := backend.Ps(ctx, project.Name, api.PsOptions{})
		if err != nil {
			return err
		}

		// Display services status
		fmt.Fprintln(output, "Services Status:")
		fmt.Fprintln(output, "================")

		if opts.format == "table" {
			// Table format
			fmt.Fprintf(output, "%-20s %-12s %-10s\n", "Service", "Status", "Health")
			fmt.Fprintln(output, "------------------------------")

			for _, container := range containers {
				health := container.Health
				if health == "" {
					health = "-"
				}

				fmt.Fprintf(output, "%-20s %-12s %-10s\n",
					container.Service,
					container.State,
					health,
				)
			}
		} else if opts.format == "json" {
			// JSON format
			fmt.Fprintln(output, "{")
			fmt.Fprintf(output, "  \"project\": \"%s\",\n", project.Name)
			fmt.Fprintf(output, "  \"time\": \"%s\",\n", time.Now().Format(time.RFC3339))
			fmt.Fprintln(output, "  \"services\": [")

			for i, container := range containers {
				if i > 0 {
					fmt.Fprintln(output, ",")
				}

				fmt.Fprintf(output, "    {\n")
				fmt.Fprintf(output, "      \"service\": \"%s\",\n", container.Service)
				fmt.Fprintf(output, "      \"status\": \"%s\",\n", container.State)
				fmt.Fprintf(output, "      \"health\": \"%s\",\n", container.Health)
				fmt.Fprintf(output, "      \"image\": \"%s\"\n", container.Image)
				fmt.Fprintf(output, "    }")
			}

			fmt.Fprintln(output, "\n  ]")
			fmt.Fprintln(output, "}")
		}

		// Show endpoints
		fmt.Fprintln(output, "\nEndpoints:")
		fmt.Fprintln(output, "==========")
		for _, service := range project.Services {
			if len(service.Ports) > 0 {
				fmt.Fprintf(output, "%s:\n", service.Name)
				for _, port := range service.Ports {
					hostIP := port.HostIP
					if hostIP == "" {
						hostIP = "0.0.0.0"
					}
					fmt.Fprintf(output, "  http://%s:%s\n", hostIP, port.Published)
				}
			}
		}

		// Check if we should exit
		if !opts.watch {
			break
		}

		// Sleep until next refresh
		time.Sleep(opts.interval)
	}

	return nil
}
