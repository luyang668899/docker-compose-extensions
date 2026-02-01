# docker compose dev

<!---MARKER_GEN_START-->
Development environment optimized for rapid development with hot reload, code sync, and debugging support.

This command supports:
1. Hot reload: Automatically restart services on code changes
2. Code sync: Real-time sync between local files and containers
3. Debugging: Support for setting breakpoints and debugging in containers
4. IDE integration: Integration with VS Code, IntelliJ, and other IDEs
5. Custom watch paths: Specify which paths to watch for changes
6. Ignore patterns: Exclude specific paths from watching


### Options

| Name               | Type          | Default  | Description                                                    |
|:-------------------|:--------------|:---------|:---------------------------------------------------------------|
| `--debug`          | `bool`        |          | Enable debugging support                                       |
| `--debug-port`     | `int`         | `5678`   | Debugging port                                                 |
| `--dry-run`        | `bool`        |          | Execute command in dry run mode                                |
| `--hot-reload`     | `bool`        | `true`   | Enable hot reload on code changes                              |
| `--ide`            | `string`      |          | IDE integration (vscode, intellij)                             |
| `--ignore`         | `stringArray` |          | Paths to ignore for changes                                    |
| `--poll-interval`  | `int`         | `2`      | Polling interval for file changes (seconds)                    |
| `--restart-policy` | `string`      | `always` | Restart policy on code changes (always, on-failure, never)     |
| `--sync`           | `string`      |          | Sync local directory to container (format: ./local:/container) |
| `--watch`          | `stringArray` |          | Paths to watch for changes                                     |


<!---MARKER_GEN_END-->

