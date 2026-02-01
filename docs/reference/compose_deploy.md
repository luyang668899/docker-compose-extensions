# docker compose deploy

<!---MARKER_GEN_START-->
Deploy services to specified environment with automated workflow.

This command supports:
1. Multi-environment deployment (dev/test/prod)
2. Automatic build and push images
3. Deployment strategies (rolling/blue-green)
4. CI/CD integration
5. Rollback to previous versions


### Options

| Name            | Type     | Default   | Description                                  |
|:----------------|:---------|:----------|:---------------------------------------------|
| `--ci`          | `bool`   |           | CI mode for integration with CI/CD pipelines |
| `--dry-run`     | `bool`   |           | Execute command in dry run mode              |
| `--env`         | `string` | `dev`     | Environment to deploy to (dev/test/prod)     |
| `--no-build`    | `bool`   |           | Skip build step                              |
| `--push`        | `bool`   |           | Push images to registry                      |
| `--rollback`    | `bool`   |           | Rollback to previous version                 |
| `--rollback-to` | `string` |           | Rollback to specific version                 |
| `--strategy`    | `string` | `rolling` | Deployment strategy (rolling/blue-green)     |


<!---MARKER_GEN_END-->

