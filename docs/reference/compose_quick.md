---
title: "docker compose quick"
description: "docker compose quick 命令的参考文档"
layout: reference
---

# docker compose quick

`docker compose quick` 命令用于简化服务操作，将 pull、build、start 和状态显示组合在一个命令中，减少用户操作步骤。

## 用法

```bash
docker compose quick [OPTIONS]
```

## 选项

| 选项 | 描述 |
|------|------|
| `-f`, `--file` | 指定 Compose 文件的路径（默认：docker-compose.yml） |
| `-p`, `--project-name` | 指定项目名称 |
| `--build` | 在启动前构建服务 |
| `--pull` | 在启动前拉取最新镜像 |
| `--quiet`, `-q` | 安静模式，减少输出信息 |
| `--help` | 显示帮助信息并退出 |

## 使用示例

### 基本用法

```bash
docker compose quick
```

该命令会执行以下操作：
1. 拉取服务所需的镜像
2. 构建需要构建的服务
3. 启动所有服务
4. 显示服务的运行状态

### 指定 Compose 文件

```bash
docker compose quick -f docker-compose.prod.yml
```

### 指定项目名称

```bash
docker compose quick -p my-project
```

### 仅构建和启动

```bash
docker compose quick --build
```

### 仅拉取和启动

```bash
docker compose quick --pull
```

## 最佳实践

- **开发环境**：使用 `docker compose quick` 快速启动开发环境，自动处理依赖和构建步骤
- **持续集成**：在 CI/CD 流程中使用该命令，简化部署步骤
- **快速测试**：在测试新功能时，使用该命令快速部署和验证服务状态

## 注意事项

- 该命令会依次执行多个操作，可能需要较长时间完成
- 建议在网络连接良好的环境下使用，以确保镜像拉取顺利
- 对于大型项目，可能需要调整超时设置以避免操作失败

## 相关命令

- `docker compose up`：启动服务
- `docker compose build`：构建服务
- `docker compose pull`：拉取镜像
- `docker compose ps`：显示服务状态