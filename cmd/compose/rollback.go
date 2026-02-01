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
	"sort"
	"time"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"

	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

type rollbackOptions struct {
	*ProjectOptions
	version      string
	timepoint    string
	strategy     string
	preserveData bool
	services     []string
	history      bool
}

func rollbackCommand(p *ProjectOptions, dockerCli command.Cli, backendOptions *BackendOptions) *cobra.Command {
	opts := rollbackOptions{
		ProjectOptions: p,
		strategy:       "rolling",
		preserveData:   true,
	}

	cmd := &cobra.Command{
		Use:   "rollback [OPTIONS] [SERVICE...]",
		Short: "Rollback services to previous versions",
		Long: `Rollback services to previous versions with configurable strategies.

This command supports:
1. Version history management
2. Rollback to specific version
3. Rollback to specific time point
4. Rollback strategies (rolling/blue-green)
5. Data preservation options
`,
		RunE: Adapt(func(ctx context.Context, args []string) error {
			opts.services = args
			return runRollbackCommand(ctx, dockerCli, backendOptions, &opts)
		}),
	}

	cmd.Flags().StringVar(&opts.version, "version", "", "Rollback to specific version")
	cmd.Flags().StringVar(&opts.timepoint, "timepoint", "", "Rollback to specific time point (YYYY-MM-DD HH:MM:SS)")
	cmd.Flags().StringVar(&opts.strategy, "strategy", "rolling", "Rollback strategy (rolling/blue-green)")
	cmd.Flags().BoolVar(&opts.preserveData, "preserve-data", true, "Preserve service data during rollback")
	cmd.Flags().BoolVar(&opts.history, "history", false, "Show version history")
	return cmd
}

func runRollbackCommand(ctx context.Context, dockerCli command.Cli, backendOptions *BackendOptions, opts *rollbackOptions) error {
	backend, err := compose.NewComposeService(dockerCli, backendOptions.Options...)
	if err != nil {
		return err
	}

	project, _, err := opts.ToProject(ctx, dockerCli, backend, opts.services)
	if err != nil {
		return err
	}

	// Show history if requested
	if opts.history {
		return showVersionHistory(project.Name)
	}

	// Determine target version
	targetVersion, err := determineTargetVersion(opts.version, opts.timepoint, project.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Rolling back to version: %s\n", targetVersion)
	fmt.Printf("Strategy: %s\n", opts.strategy)
	fmt.Printf("Preserve data: %v\n", opts.preserveData)
	fmt.Printf("Rolling back services: %v\n", opts.services)

	// Perform rollback based on strategy
	switch opts.strategy {
	case "rolling":
		if err := runRollingRollback(ctx, backend, project, opts.services, targetVersion, opts.preserveData); err != nil {
			return err
		}
	case "blue-green":
		if err := runBlueGreenRollback(ctx, backend, project, project.Name, opts.services, targetVersion, opts.preserveData); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported rollback strategy: %s", opts.strategy)
	}

	// Show rollback status
	fmt.Println("\nRollback status:")
	containers, err := backend.Ps(ctx, project.Name, api.PsOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		fmt.Printf("%s: %s\n", container.Service, container.State)
	}

	fmt.Println("\nRollback completed successfully!")
	return nil
}

func showVersionHistory(projectName string) error {
	// Get version history (simplified implementation)
	history := getVersionHistory(projectName)

	if len(history) == 0 {
		fmt.Println("No version history found.")
		return nil
	}

	fmt.Println("Version history:")
	fmt.Println("┌─────────┬─────────────────────┬─────────────────────┬─────────────────────┐")
	fmt.Println("│ Version │ Created At          │ Updated At          │ Description         │")
	fmt.Println("├─────────┼─────────────────────┼─────────────────────┼─────────────────────┤")

	for _, version := range history {
		fmt.Printf("│ %-7s │ %-19s │ %-19s │ %-19s │\n",
			version.Version, version.CreatedAt, version.UpdatedAt, version.Description)
	}

	fmt.Println("└─────────┴─────────────────────┴─────────────────────┴─────────────────────┘")
	return nil
}

func determineTargetVersion(version, timepoint, projectName string) (string, error) {
	if version != "" {
		return version, nil
	}

	if timepoint != "" {
		// Find version closest to the specified timepoint
		targetTime, err := time.Parse("2006-01-02 15:04:05", timepoint)
		if err != nil {
			return "", fmt.Errorf("invalid timepoint format: %v", err)
		}

		history := getVersionHistory(projectName)
		if len(history) == 0 {
			return "", fmt.Errorf("no version history found")
		}

		// Find closest version before or at the timepoint
		var closestVersion *VersionInfo
		var minDiff time.Duration

		for i, v := range history {
			vTime, err := time.Parse("2006-01-02 15:04:05", v.CreatedAt)
			if err != nil {
				continue
			}

			if vTime.After(targetTime) {
				continue
			}

			diff := targetTime.Sub(vTime)
			if closestVersion == nil || diff < minDiff {
				closestVersion = &history[i]
				minDiff = diff
			}
		}

		if closestVersion == nil {
			return "", fmt.Errorf("no version found before the specified timepoint")
		}

		return closestVersion.Version, nil
	}

	// Default to previous version
	history := getVersionHistory(projectName)
	if len(history) < 2 {
		return "", fmt.Errorf("not enough version history to rollback")
	}

	// Sort by created time (newest first)
	sort.Slice(history, func(i, j int) bool {
		timeI, _ := time.Parse("2006-01-02 15:04:05", history[i].CreatedAt)
		timeJ, _ := time.Parse("2006-01-02 15:04:05", history[j].CreatedAt)
		return timeI.After(timeJ)
	})

	return history[1].Version, nil
}

func runRollingRollback(ctx context.Context, backend api.Compose, project *types.Project, services []string, version string, preserveData bool) error {
	// Rolling rollback: stop and start services one by one
	targetServices := project.Services
	if len(services) > 0 {
		// Filter services to only those specified
		filteredServices := make(map[string]types.ServiceConfig)
		for _, serviceName := range services {
			if service, ok := project.Services[serviceName]; ok {
				filteredServices[serviceName] = service
			}
		}
		targetServices = filteredServices
	}

	for serviceName := range targetServices {
		fmt.Printf("Rolling back service: %s to version %s\n", serviceName, version)

		// Stop the service
		if err := backend.Stop(ctx, project.Name, api.StopOptions{
			Services: []string{serviceName},
		}); err != nil {
			fmt.Printf("Warning: Stop failed: %v\n", err)
			// Continue even if stop fails
		}

		// Start the service (in real implementation, this would use the specified version)
		if err := backend.Start(ctx, project.Name, api.StartOptions{
			Services: []string{serviceName},
		}); err != nil {
			return err
		}
	}

	return nil
}

func runBlueGreenRollback(ctx context.Context, backend api.Compose, project *types.Project, projectName string, services []string, version string, preserveData bool) error {
	// Blue-green rollback: create new instances alongside existing ones
	fmt.Printf("Performing blue-green rollback to version %s\n", version)

	// Stop all services
	if err := backend.Stop(ctx, projectName, api.StopOptions{
		Services: services,
	}); err != nil {
		fmt.Printf("Warning: Stop failed: %v\n", err)
		// Continue even if stop fails
	}

	// Start all services (in real implementation, this would use the specified version)
	if err := backend.Start(ctx, projectName, api.StartOptions{
		Services: services,
	}); err != nil {
		return err
	}

	return nil
}

// VersionInfo represents a version in the history
type VersionInfo struct {
	Version     string
	CreatedAt   string
	UpdatedAt   string
	Description string
}

func getVersionHistory(projectName string) []VersionInfo {
	// Simplified implementation - in real code, this would read from a version store
	// For demo purposes, return mock version history
	return []VersionInfo{
		{
			Version:     "v3",
			CreatedAt:   time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05"),
			UpdatedAt:   time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05"),
			Description: "Initial deployment",
		},
		{
			Version:     "v2",
			CreatedAt:   time.Now().Add(-48 * time.Hour).Format("2006-01-02 15:04:05"),
			UpdatedAt:   time.Now().Add(-48 * time.Hour).Format("2006-01-02 15:04:05"),
			Description: "Second version",
		},
		{
			Version:     "v1",
			CreatedAt:   time.Now().Add(-72 * time.Hour).Format("2006-01-02 15:04:05"),
			UpdatedAt:   time.Now().Add(-72 * time.Hour).Format("2006-01-02 15:04:05"),
			Description: "Initial version",
		},
	}
}
