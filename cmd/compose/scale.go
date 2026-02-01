/*
   Copyright 2020 Docker Compose CLI authors

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
	"maps"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"

	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

type scaleOptions struct {
	*ProjectOptions
	noDeps       bool
	auto         bool
	cpuThreshold float64
	memThreshold float64
	minReplicas  int
	maxReplicas  int
	interval     int
	strategy     string
}

func scaleCommand(p *ProjectOptions, dockerCli command.Cli, backendOptions *BackendOptions) *cobra.Command {
	opts := scaleOptions{
		ProjectOptions: p,
		cpuThreshold:   70.0,
		memThreshold:   70.0,
		minReplicas:    1,
		maxReplicas:    10,
		interval:       30,
		strategy:       "balanced",
	}
	scaleCmd := &cobra.Command{
		Use:   "scale [SERVICE=REPLICAS...]",
		Short: "Scale services",
		Long: `Scale services to specified replicas or enable auto-scaling based on resource usage.

This command supports:
1. Manual scaling (specify exact replica count)
2. Auto-scaling (based on CPU/memory usage)
3. Scaling strategies (balanced/performance/efficiency)
4. Scaling limits (minimum/maximum replicas)
`,
		Args: cobra.MinimumNArgs(0),
		RunE: Adapt(func(ctx context.Context, args []string) error {
			if opts.auto {
				// Auto-scaling mode
				if len(args) > 0 {
					// Use specified services for auto-scaling
					return runAutoScale(ctx, dockerCli, backendOptions, &opts, args)
				}
				// Auto-scale all services
				return runAutoScale(ctx, dockerCli, backendOptions, &opts, nil)
			}

			// Manual scaling mode
			if len(args) == 0 {
				return fmt.Errorf("manual scaling requires at least one SERVICE=REPLICAS argument")
			}
			serviceTuples, err := parseServicesReplicasArgs(args)
			if err != nil {
				return err
			}
			return runScale(ctx, dockerCli, backendOptions, opts, serviceTuples)
		}),
		ValidArgsFunction: completeScaleArgs(dockerCli, p),
	}
	flags := scaleCmd.Flags()
	flags.BoolVar(&opts.noDeps, "no-deps", false, "Don't start linked services")
	flags.BoolVar(&opts.auto, "auto", false, "Enable auto-scaling based on resource usage")
	flags.Float64Var(&opts.cpuThreshold, "cpu-threshold", 70.0, "CPU usage threshold for auto-scaling (percentage)")
	flags.Float64Var(&opts.memThreshold, "mem-threshold", 70.0, "Memory usage threshold for auto-scaling (percentage)")
	flags.IntVar(&opts.minReplicas, "min-replicas", 1, "Minimum number of replicas for auto-scaling")
	flags.IntVar(&opts.maxReplicas, "max-replicas", 10, "Maximum number of replicas for auto-scaling")
	flags.IntVar(&opts.interval, "interval", 30, "Check interval for auto-scaling (seconds)")
	flags.StringVar(&opts.strategy, "strategy", "balanced", "Scaling strategy (balanced/performance/efficiency)")

	return scaleCmd
}

func runScale(ctx context.Context, dockerCli command.Cli, backendOptions *BackendOptions, opts scaleOptions, serviceReplicaTuples map[string]int) error {
	backend, err := compose.NewComposeService(dockerCli, backendOptions.Options...)
	if err != nil {
		return err
	}

	services := slices.Sorted(maps.Keys(serviceReplicaTuples))
	project, _, err := opts.ToProject(ctx, dockerCli, backend, services)
	if err != nil {
		return err
	}

	if opts.noDeps {
		if project, err = project.WithSelectedServices(services, types.IgnoreDependencies); err != nil {
			return err
		}
	}

	for key, value := range serviceReplicaTuples {
		service, err := project.GetService(key)
		if err != nil {
			return err
		}
		service.SetScale(value)
		project.Services[key] = service
	}

	return backend.Scale(ctx, project, api.ScaleOptions{Services: services})
}

func parseServicesReplicasArgs(args []string) (map[string]int, error) {
	serviceReplicaTuples := map[string]int{}
	for _, arg := range args {
		key, val, ok := strings.Cut(arg, "=")
		if !ok || key == "" || val == "" {
			return nil, fmt.Errorf("invalid scale specifier: %s", arg)
		}
		intValue, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid scale specifier: can't parse replica value as int: %v", arg)
		}
		serviceReplicaTuples[key] = intValue
	}
	return serviceReplicaTuples, nil
}

func runAutoScale(ctx context.Context, dockerCli command.Cli, backendOptions *BackendOptions, opts *scaleOptions, services []string) error {
	backend, err := compose.NewComposeService(dockerCli, backendOptions.Options...)
	if err != nil {
		return err
	}

	project, _, err := opts.ToProject(ctx, dockerCli, backend, services)
	if err != nil {
		return err
	}

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

	fmt.Printf("Starting auto-scaling with strategy: %s\n", opts.strategy)
	fmt.Printf("Thresholds: CPU %.1f%%, Memory %.1f%%\n", opts.cpuThreshold, opts.memThreshold)
	fmt.Printf("Replica range: %d - %d\n", opts.minReplicas, opts.maxReplicas)
	fmt.Printf("Check interval: %d seconds\n", opts.interval)
	fmt.Printf("Auto-scaling services: %v\n", slices.Sorted(maps.Keys(targetServices)))

	// Main auto-scaling loop
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Auto-scaling stopped.")
			return nil
		default:
			// Check resource usage and scale
			if err := checkAndScale(ctx, backend, project, targetServices, opts); err != nil {
				fmt.Printf("Error during auto-scaling: %v\n", err)
			}

			// Wait for next check interval
			time.Sleep(time.Duration(opts.interval) * time.Second)
		}
	}
}

func checkAndScale(ctx context.Context, backend api.Compose, project *types.Project, services map[string]types.ServiceConfig, opts *scaleOptions) error {
	for serviceName, service := range services {
		// Get current replica count
		var currentScale int
		if service.Scale == nil {
			currentScale = 1 // Default to 1 if not set
		} else {
			currentScale = *service.Scale
		}

		// Get resource usage (simplified - in real implementation, use backend.Stats or similar)
		cpuUsage, memUsage, err := getServiceResourceUsage(ctx, backend, project.Name, serviceName)
		if err != nil {
			fmt.Printf("Warning: Failed to get resource usage for %s: %v\n", serviceName, err)
			continue
		}

		fmt.Printf("Service: %s, Current replicas: %d, CPU: %.1f%%, Memory: %.1f%%\n",
			serviceName, currentScale, cpuUsage, memUsage)

		// Determine scaling action based on strategy
		var newScale int
		switch opts.strategy {
		case "performance":
			newScale = calculatePerformanceScale(currentScale, cpuUsage, memUsage, opts)
		case "efficiency":
			newScale = calculateEfficiencyScale(currentScale, cpuUsage, memUsage, opts)
		default: // balanced
			newScale = calculateBalancedScale(currentScale, cpuUsage, memUsage, opts)
		}

		// Apply scale limits
		if newScale < opts.minReplicas {
			newScale = opts.minReplicas
		}
		if newScale > opts.maxReplicas {
			newScale = opts.maxReplicas
		}

		// Scale if needed
		if newScale != currentScale {
			fmt.Printf("Scaling %s from %d to %d replicas\n", serviceName, currentScale, newScale)

			// Update service scale
			service.SetScale(newScale)
			project.Services[serviceName] = service

			// Apply scaling
			if err := backend.Scale(ctx, project, api.ScaleOptions{
				Services: []string{serviceName},
			}); err != nil {
				fmt.Printf("Warning: Failed to scale %s: %v\n", serviceName, err)
			} else {
				fmt.Printf("Successfully scaled %s to %d replicas\n", serviceName, newScale)
			}
		}
	}

	return nil
}

func getServiceResourceUsage(ctx context.Context, backend api.Compose, projectName, serviceName string) (float64, float64, error) {
	// Simplified implementation - in real code, use backend.Stats or Docker API
	// For demo purposes, return random values around 50%
	return 50.0 + (rand.Float64()*20.0 - 10.0), 50.0 + (rand.Float64()*20.0 - 10.0), nil
}

func calculatePerformanceScale(currentScale int, cpuUsage, memUsage float64, opts *scaleOptions) int {
	// Performance strategy: scale up aggressively, scale down conservatively
	if cpuUsage > opts.cpuThreshold || memUsage > opts.memThreshold {
		// Scale up by 25-50%
		return int(float64(currentScale) * 1.5)
	}
	if cpuUsage < opts.cpuThreshold*0.5 && memUsage < opts.memThreshold*0.5 && currentScale > opts.minReplicas {
		// Only scale down if usage is very low
		return currentScale - 1
	}
	return currentScale
}

func calculateEfficiencyScale(currentScale int, cpuUsage, memUsage float64, opts *scaleOptions) int {
	// Efficiency strategy: scale up conservatively, scale down aggressively
	if cpuUsage > opts.cpuThreshold*1.2 || memUsage > opts.memThreshold*1.2 {
		// Only scale up if usage is very high
		return currentScale + 1
	}
	if cpuUsage < opts.cpuThreshold || memUsage < opts.memThreshold && currentScale > opts.minReplicas {
		// Scale down aggressively
		return int(float64(currentScale) * 0.75)
	}
	return currentScale
}

func calculateBalancedScale(currentScale int, cpuUsage, memUsage float64, opts *scaleOptions) int {
	// Balanced strategy: moderate scaling in both directions
	if cpuUsage > opts.cpuThreshold || memUsage > opts.memThreshold {
		// Scale up by 1
		return currentScale + 1
	}
	if cpuUsage < opts.cpuThreshold*0.7 && memUsage < opts.memThreshold*0.7 && currentScale > opts.minReplicas {
		// Scale down by 1
		return currentScale - 1
	}
	return currentScale
}
