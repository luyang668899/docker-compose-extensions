---
title: "docker compose monitor"
description: "docker compose monitor 命令的参考文档"
layout: reference
---

# docker compose monitor

`docker compose monitor` 命令用于实时监控服务状态，显示服务的运行状态和端点信息，帮助用户快速了解服务的健康状况。

## 用法

```bash
docker compose monitor [OPTIONS]
```

## 选项

| 选项 | 描述 |
|------|------|
| `-f`, `--file` | 指定 Compose 文件的路径（默认：docker-compose.yml） |
| `-p`, `--project-name` | 指定项目名称 |
| `--interval` | 设置状态检查间隔时间（秒），默认：5 |
| `--format` | 输出格式，支持 table、json（默认：table） |
| `--quiet`, `-q` | 安静模式，减少输出信息 |
| `--help` | 显示帮助信息并退出 |

## 使用示例

### 基本用法

```bash
docker compose monitor
```

该命令会：
1. 显示所有服务的当前状态
2. 每 5 秒更新一次状态信息
3. 显示服务的访问端点（如果有）

### 指定检查间隔

```bash
docker compose monitor --interval 10
```

### 指定输出格式为 JSON

```bash
docker compose monitor --format json
```

### 指定 Compose 文件

```bash
docker compose monitor -f docker-compose.prod.yml
```

### 指定项目名称

```bash
docker compose monitor -p my-project
```

## 输出示例

### 表格格式输出

```
┌───────────┬──────────┬────────────┬───────────┐
│ 服务名称  │ 状态     │ 容器 ID    │ 端点      │
├───────────┼──────────┼────────────┼───────────┤
│ web       │ 运行中   │ abc123     │ http://localhost:8080 │
│ api       │ 运行中   │ def456     │ http://localhost:3000 │
│ db        │ 运行中   │ ghi789     │ -         │
└───────────┴──────────┴────────────┴───────────┘
```

### JSON 格式输出

```json
[
  {
    "service": "web",
    "status": "running",
    "container_id": "abc123",
    "endpoints": ["http://localhost:8080"]
  },
  {
    "service": "api",
    "status": "running",
    "container_id": "def456",
    "endpoints": ["http://localhost:3000"]
  },
  {
    "service": "db",
    "status": "running",
    "container_id": "ghi789",
    "endpoints": []
  }
]
```

## 最佳实践

- **服务监控**：在开发和生产环境中使用该命令监控服务状态
- **问题排查**：当服务出现异常时，使用该命令快速查看服务状态和端点信息
- **自动化脚本**：结合 `--format json` 选项，在自动化脚本中解析服务状态
- **持续集成**：在 CI/CD 流程中使用该命令验证服务是否正常启动

## 注意事项

- 该命令会持续运行并输出状态信息，按 Ctrl+C 退出
- 对于大型项目，建议适当增加检查间隔时间，减少系统负载
- 端点信息仅显示已映射的端口，未映射的端口不会显示

## 相关命令

- `docker compose ps`：显示服务状态（单次）
- `docker compose logs`：查看服务日志
- `docker compose port`：显示服务端口映射