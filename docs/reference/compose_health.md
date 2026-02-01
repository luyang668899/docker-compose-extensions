---
title: "docker compose health"
description: "docker compose health 命令的参考文档"
layout: reference
---

# docker compose health

`docker compose health` 命令用于健康检查管理，显示服务的健康状态，帮助用户了解服务是否正常运行。

## 用法

```bash
docker compose health [OPTIONS]
```

## 选项

| 选项 | 描述 |
|------|------|
| `-f`, `--file` | 指定 Compose 文件的路径（默认：docker-compose.yml） |
| `-p`, `--project-name` | 指定项目名称 |
| `--interval` | 设置健康检查间隔时间（秒），默认：5 |
| `--format` | 输出格式，支持 table、json（默认：table） |
| `--quiet`, `-q` | 安静模式，减少输出信息 |
| `--help` | 显示帮助信息并退出 |

## 使用示例

### 基本用法

```bash
docker compose health
```

该命令会：
1. 显示所有服务的当前健康状态
2. 每 5 秒更新一次健康状态信息
3. 显示健康检查的详细结果（如果有）

### 指定检查间隔

```bash
docker compose health --interval 10
```

### 指定输出格式为 JSON

```bash
docker compose health --format json
```

### 指定 Compose 文件

```bash
docker compose health -f docker-compose.prod.yml
```

### 指定项目名称

```bash
docker compose health -p my-project
```

## 输出示例

### 表格格式输出

```
┌───────────┬──────────┬────────────┬───────────┐
│ 服务名称  │ 健康状态 │ 容器 ID    │ 检查结果  │
├───────────┼──────────┼────────────┼───────────┤
│ web       │ 健康     │ abc123     │ HTTP 200  │
│ api       │ 健康     │ def456     │ HTTP 200  │
│ db        │ 健康     │ ghi789     │ TCP 连接成功 │
└───────────┴──────────┴────────────┴───────────┘
```

### JSON 格式输出

```json
[
  {
    "service": "web",
    "status": "healthy",
    "container_id": "abc123",
    "check_result": "HTTP 200"
  },
  {
    "service": "api",
    "status": "healthy",
    "container_id": "def456",
    "check_result": "HTTP 200"
  },
  {
    "service": "db",
    "status": "healthy",
    "container_id": "ghi789",
    "check_result": "TCP 连接成功"
  }
]
```

## 健康状态说明

| 状态 | 描述 |
|------|------|
| `healthy` | 服务健康，所有健康检查通过 |
| `unhealthy` | 服务不健康，健康检查失败 |
| `starting` | 服务正在启动，健康检查尚未完成 |
| `no_healthcheck` | 服务未配置健康检查 |
| `unknown` | 健康状态未知 |

## 最佳实践

- **服务监控**：在生产环境中使用该命令监控服务健康状态
- **问题排查**：当服务出现异常时，使用该命令快速查看健康检查结果
- **自动化脚本**：结合 `--format json` 选项，在自动化脚本中解析健康状态
- **持续集成**：在 CI/CD 流程中使用该命令验证服务是否正常运行

## 注意事项

- 该命令会持续运行并输出健康状态信息，按 Ctrl+C 退出
- 只有在 Compose 文件中配置了健康检查的服务才会显示详细的健康状态
- 对于未配置健康检查的服务，会显示 "no_healthcheck" 状态

## 相关命令

- `docker compose ps`：显示服务状态
- `docker compose logs`：查看服务日志
- `docker compose up`：启动服务并等待健康检查通过