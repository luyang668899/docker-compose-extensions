# docker compose env

<!---MARKER_GEN_START-->
The `docker compose env` command helps you manage different environment configurations for your Compose projects. It allows you to create, activate, switch between, and manage environments such as development, testing, and production.

### Options

| Name            | Type     | Default | Description                     |
|:----------------|:---------|:--------|:--------------------------------|
| `--activate`    | `bool`   |         | Activate environment            |
| `--create`      | `bool`   |         | Create new environment          |
| `--deactivate`  | `bool`   |         | Deactivate current environment  |
| `--description` | `string` |         | Environment description         |
| `--dry-run`     | `bool`   |         | Execute command in dry run mode |
| `--export`      | `string` |         | Export environment to file      |
| `--import`      | `string` |         | Import environment from file    |
| `--list`        | `bool`   |         | List available environments     |
| `--remove`      | `bool`   |         | Remove environment              |


<!---MARKER_GEN_END-->

## Description

The `docker compose env` command helps you manage different environment configurations for your Compose projects. It allows you to create, activate, switch between, and manage environments such as development, testing, and production.

## Usage

### List available environments

```bash
docker compose env --list
```

This will display all available environments, including their descriptions and which one is currently active.

### Create a new environment

```bash
docker compose env --create --description "Development environment" dev
```

This creates a new environment named `dev` with the description "Development environment".

### Activate an environment

```bash
docker compose env --activate dev
```

This activates the `dev` environment, making it the current active environment.

### Deactivate the current environment

```bash
docker compose env --deactivate
```

This deactivates the currently active environment.

### Remove an environment

```bash
docker compose env --remove dev
```

This removes the `dev` environment. If it's currently active, it will be deactivated first.

### Import an environment from a file

```bash
docker compose env --import path/to/compose.yaml --create staging
```

This creates a new environment named `staging` and imports the configuration from the specified file.

### Export an environment to a file

```bash
docker compose env --export path/to/output.yaml production
```

This exports the `production` environment configuration to the specified file.

## Examples

### Example 1: Basic environment management

```bash
# Create environments
docker compose env --create --description "Development environment" dev
docker compose env --create --description "Testing environment" test
docker compose env --create --description "Production environment" prod

# List environments
docker compose env --list

# Activate development environment
docker compose env --activate dev

# Use the environment
docker compose up

# Switch to production environment
docker compose env --activate prod

# Use the production environment
docker compose up
```

### Example 2: Environment-specific configurations

You can create different environment configurations for different purposes:

```bash
# Create a development environment with debug settings
docker compose env --create --description "Development with debug" dev-debug

# Create a production environment with optimized settings
docker compose env --create --description "Production optimized" prod-optimized
```

## Environment Storage

Environments are stored in a platform-specific directory:

- **Linux/macOS**: `~/.docker/compose/environments/`
- **Windows**: `%USERPROFILE%\.docker\compose\environments\`

Each environment has its own directory containing:
- `compose.yaml`: The Compose file for the environment
- `.env`: Environment variables file
- `description.txt`: Environment description

## Best Practices

1. **Use descriptive names**: Name environments clearly (e.g., `dev`, `test`, `prod`)
2. **Add descriptions**: Use the `--description` flag to document the purpose of each environment
3. **Version control**: Consider keeping your environment configurations in version control
4. **Environment isolation**: Use separate environments for different stages of development and deployment
5. **Secret management**: Be careful with sensitive information in environment files

## Notes

- This command is experimental and subject to change
- Environment names should be unique within a project
- When activating an environment, the command provides instructions on how to use it
- Environment configurations are project-specific
