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
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

type envOptions struct {
	*ProjectOptions
	name        string
	list        bool
	activate    bool
	deactivate  bool
	create      bool
	remove      bool
	importFile  string
	exportFile  string
	description string
}

func envCommand(p *ProjectOptions, dockerCli command.Cli, backendOptions *BackendOptions) *cobra.Command {
	opts := envOptions{
		ProjectOptions: p,
	}

	cmd := &cobra.Command{
		Use:   "env [OPTIONS]",
		Short: "Manage environment configurations",
		Long: `EXPERIMENTAL - Manage environment configurations for Compose projects.

This command helps you create and manage different environment configurations
(development, testing, production) and easily switch between them.
`,
		RunE: Adapt(func(ctx context.Context, args []string) error {
			if len(args) > 0 {
				opts.name = args[0]
			}
			return runEnv(ctx, dockerCli, backendOptions, &opts)
		}),
	}

	cmd.Flags().BoolVar(&opts.list, "list", false, "List available environments")
	cmd.Flags().BoolVar(&opts.activate, "activate", false, "Activate environment")
	cmd.Flags().BoolVar(&opts.deactivate, "deactivate", false, "Deactivate current environment")
	cmd.Flags().BoolVar(&opts.create, "create", false, "Create new environment")
	cmd.Flags().BoolVar(&opts.remove, "remove", false, "Remove environment")
	cmd.Flags().StringVar(&opts.importFile, "import", "", "Import environment from file")
	cmd.Flags().StringVar(&opts.exportFile, "export", "", "Export environment to file")
	cmd.Flags().StringVar(&opts.description, "description", "", "Environment description")
	return cmd
}

func runEnv(ctx context.Context, dockerCli command.Cli, backendOptions *BackendOptions, opts *envOptions) error {
	// Get environments directory
	envsDir := getEnvironmentsDir()
	if err := os.MkdirAll(envsDir, 0755); err != nil {
		return fmt.Errorf("failed to create environments directory: %v", err)
	}

	// List environments
	if opts.list {
		return listEnvironments(envsDir)
	}

	// Create environment
	if opts.create {
		if opts.name == "" {
			return fmt.Errorf("environment name is required")
		}
		return createEnvironment(envsDir, opts.name, opts.description)
	}

	// Remove environment
	if opts.remove {
		if opts.name == "" {
			return fmt.Errorf("environment name is required")
		}
		return removeEnvironment(envsDir, opts.name)
	}

	// Activate environment
	if opts.activate {
		if opts.name == "" {
			return fmt.Errorf("environment name is required")
		}
		return activateEnvironment(envsDir, opts.name)
	}

	// Deactivate environment
	if opts.deactivate {
		return deactivateEnvironment(envsDir)
	}

	// Import environment
	if opts.importFile != "" {
		if opts.name == "" {
			return fmt.Errorf("environment name is required")
		}
		return importEnvironment(envsDir, opts.name, opts.importFile)
	}

	// Export environment
	if opts.exportFile != "" {
		if opts.name == "" {
			return fmt.Errorf("environment name is required")
		}
		return exportEnvironment(envsDir, opts.name, opts.exportFile)
	}

	// Show current environment
	return showCurrentEnvironment(envsDir)
}

func getEnvironmentsDir() string {
	// Get user config directory based on platform
	var configDir string
	switch {
	case os.Getenv("HOME") != "":
		// Unix-like systems
		configDir = filepath.Join(os.Getenv("HOME"), ".docker", "compose", "environments")
	case os.Getenv("USERPROFILE") != "":
		// Windows
		configDir = filepath.Join(os.Getenv("USERPROFILE"), ".docker", "compose", "environments")
	default:
		// Fallback
		configDir = ".docker-compose-environments"
	}
	return configDir
}

func listEnvironments(envsDir string) error {
	files, err := os.ReadDir(envsDir)
	if err != nil {
		return err
	}

	fmt.Println("Available environments:")
	fmt.Println("=====================")

	// Get current environment
	currentEnv, _ := getCurrentEnvironment(envsDir)

	for _, file := range files {
		if file.IsDir() {
			status := ""
			if file.Name() == currentEnv {
				status = " [ACTIVE]"
			}
			
			// Read description
			descFile := filepath.Join(envsDir, file.Name(), "description.txt")
			desc, err := os.ReadFile(descFile)
			description := ""
			if err == nil {
				description = strings.TrimSpace(string(desc))
			}
			
			fmt.Printf("%s%s\n", file.Name(), status)
			if description != "" {
				fmt.Printf("  Description: %s\n", description)
			}
		}
	}

	if len(files) == 0 {
		fmt.Println("No environments found. Use 'docker compose env --create' to create one.")
	}

	return nil
}

func createEnvironment(envsDir, name, description string) error {
	envDir := filepath.Join(envsDir, name)
	if _, err := os.Stat(envDir); err == nil {
		return fmt.Errorf("environment %q already exists", name)
	}

	// Create environment directory
	if err := os.MkdirAll(envDir, 0755); err != nil {
		return fmt.Errorf("failed to create environment directory: %v", err)
	}

	// Create description file
	if description != "" {
		descFile := filepath.Join(envDir, "description.txt")
		if err := os.WriteFile(descFile, []byte(description), 0644); err != nil {
			return fmt.Errorf("failed to write description: %v", err)
		}
	}

	// Create default compose.yaml template
	composeFile := filepath.Join(envDir, "compose.yaml")
	defaultCompose := `# Environment: ` + name + `
# Generated by docker compose env

services:
  # Add your services here
`
	if err := os.WriteFile(composeFile, []byte(defaultCompose), 0644); err != nil {
		return fmt.Errorf("failed to create compose.yaml: %v", err)
	}

	// Create .env file
	envFile := filepath.Join(envDir, ".env")
	defaultEnv := `# Environment variables for ` + name + `
# Generated by docker compose env
`
	if err := os.WriteFile(envFile, []byte(defaultEnv), 0644); err != nil {
		return fmt.Errorf("failed to create .env file: %v", err)
	}

	fmt.Printf("Environment %q created successfully!\n", name)
	fmt.Printf("Location: %s\n", envDir)
	return nil
}

func removeEnvironment(envsDir, name string) error {
	envDir := filepath.Join(envsDir, name)
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		return fmt.Errorf("environment %q does not exist", name)
	}

	// Check if it's the current active environment
	currentEnv, _ := getCurrentEnvironment(envsDir)
	if currentEnv == name {
		// Deactivate first
		if err := deactivateEnvironment(envsDir); err != nil {
			return err
		}
	}

	// Remove environment directory
	if err := os.RemoveAll(envDir); err != nil {
		return fmt.Errorf("failed to remove environment: %v", err)
	}

	fmt.Printf("Environment %q removed successfully!\n", name)
	return nil
}

func activateEnvironment(envsDir, name string) error {
	envDir := filepath.Join(envsDir, name)
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		return fmt.Errorf("environment %q does not exist", name)
	}

	// Write current environment
	currentEnvFile := filepath.Join(envsDir, "current")
	if err := os.WriteFile(currentEnvFile, []byte(name), 0644); err != nil {
		return fmt.Errorf("failed to activate environment: %v", err)
	}

	fmt.Printf("Environment %q activated successfully!\n", name)
	fmt.Printf("To use this environment, run: docker compose --env-file %s/.env up\n", envDir)
	return nil
}

func deactivateEnvironment(envsDir string) error {
	currentEnvFile := filepath.Join(envsDir, "current")
	if _, err := os.Stat(currentEnvFile); os.IsNotExist(err) {
		return fmt.Errorf("no active environment")
	}

	if err := os.Remove(currentEnvFile); err != nil {
		return fmt.Errorf("failed to deactivate environment: %v", err)
	}

	fmt.Println("Environment deactivated successfully!")
	return nil
}

func importEnvironment(envsDir, name, importFile string) error {
	// Check if import file exists
	if _, err := os.Stat(importFile); os.IsNotExist(err) {
		return fmt.Errorf("import file %q does not exist", importFile)
	}

	// Create environment
	if err := createEnvironment(envsDir, name, "Imported environment"); err != nil {
		return err
	}

	// Copy import file
	envDir := filepath.Join(envsDir, name)
	destFile := filepath.Join(envDir, "compose.yaml")
	content, err := os.ReadFile(importFile)
	if err != nil {
		return fmt.Errorf("failed to read import file: %v", err)
	}

	if err := os.WriteFile(destFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write compose.yaml: %v", err)
	}

	fmt.Printf("Environment %q imported successfully from %q!\n", name, importFile)
	return nil
}

func exportEnvironment(envsDir, name, exportFile string) error {
	envDir := filepath.Join(envsDir, name)
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		return fmt.Errorf("environment %q does not exist", name)
	}

	// Read compose.yaml
	composeFile := filepath.Join(envDir, "compose.yaml")
	content, err := os.ReadFile(composeFile)
	if err != nil {
		return fmt.Errorf("failed to read compose.yaml: %v", err)
	}

	// Write to export file
	if err := os.WriteFile(exportFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %v", err)
	}

	fmt.Printf("Environment %q exported successfully to %q!\n", name, exportFile)
	return nil
}

func showCurrentEnvironment(envsDir string) error {
	currentEnv, err := getCurrentEnvironment(envsDir)
	if err != nil {
		fmt.Println("No active environment")
		fmt.Println("Use 'docker compose env --activate' to activate an environment")
		return nil
	}

	envDir := filepath.Join(envsDir, currentEnv)
	
	// Read description
	descFile := filepath.Join(envDir, "description.txt")
	desc, err := os.ReadFile(descFile)
	description := ""
	if err == nil {
		description = strings.TrimSpace(string(desc))
	}

	fmt.Println("Current environment:")
	fmt.Println("==================")
	fmt.Printf("Name: %s\n", currentEnv)
	if description != "" {
		fmt.Printf("Description: %s\n", description)
	}
	fmt.Printf("Location: %s\n", envDir)
	fmt.Printf("\nTo use this environment:\n")
	fmt.Printf("  docker compose --env-file %s/.env up\n", envDir)

	return nil
}

func getCurrentEnvironment(envsDir string) (string, error) {
	currentEnvFile := filepath.Join(envsDir, "current")
	content, err := os.ReadFile(currentEnvFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}
