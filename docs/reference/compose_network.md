---
title: "docker compose network"
description: "docker compose network 命令的参考文档"
layout: reference
---

# docker compose network

`docker compose network` 命令用于网络管理，显示服务与网络之间的关联关系，帮助用户了解容器网络配置。

## 用法

```bash
docker compose network [OPTIONS]
```

## 选项

| 选项 | 描述 |
|------|------|
| `-f`, `--file` | 指定 Compose 文件的路径（默认：docker-compose.yml） |
| `-p`, `--project-name` | 指定项目名称 |
| `--format` | 输出格式，支持 table、json（默认：table） |
| `--quiet`, `-q` | 安静模式，减少输出信息 |
| `--help` | 显示帮助信息并退出 |

## 使用示例

### 基本用法

```bash
docker compose network
```

该命令会：
1. 显示项目中定义的所有网络
2. 显示每个服务关联的网络
3. 显示网络的驱动类型和子网信息（如果有）

### 指定输出格式为 JSON

```bash
docker compose network --format json
```

### 指定 Compose 文件

```bash
docker compose network -f docker-compose.prod.yml
```

### 指定项目名称

```bash
docker compose network -p my-project
```

## 输出示例

### 表格格式输出

```
┌───────────┬───────────┬───────────┐
│ 服务名称  │ 网络名称  │ 驱动类型  │
├───────────┼───────────┼───────────┤
│ web       │ default   │ bridge    │
│ api       │ default   │ bridge    │
│ api       │ backend   │ bridge    │
│ db        │ backend   │ bridge    │
└───────────┴───────────┴───────────┘
```

### JSON 格式输出

```json
[
  {
    "service": "web",
    "networks": [
      {
        "name": "default",
        "driver": "bridge"
      }
    ]
  },
  {
    "service": "api",
    "networks": [
      {
        "name": "default",
        "driver": "bridge"
      },
      {
        "name": "backend",
        "driver": "bridge"
      }
    ]
  },
  {
    "service": "db",
    "networks": [
      {
        "name": "backend",
        "driver": "bridge"
      }
    ]
  }
]
```

## 最佳实践

- **网络规划**：在设计多服务架构时，使用该命令查看网络关联，确保服务间通信正常
- **问题排查**：当服务间通信出现问题时，使用该命令检查网络配置
- **安全审计**：验证服务是否只加入了必要的网络，避免不必要的网络暴露
- **文档生成**：结合 `--format json` 选项，自动生成网络拓扑文档

## 注意事项

- 该命令仅显示 Compose 文件中定义的网络关联
- 对于使用默认网络的服务，会显示 "default" 网络
- 网络驱动类型取决于 Compose 文件中的配置

## 相关命令

- `docker network ls`：列出所有 Docker 网络
- `docker network inspect`：查看网络详细信息
- `docker compose config`：查看 Compose 文件配置