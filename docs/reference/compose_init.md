---
title: "docker compose init"
description: "docker compose init 命令的参考文档"
layout: reference
---

# docker compose init

`docker compose init` 命令用于项目初始化，根据预定义模板创建 Docker Compose 项目结构，支持 web、api、microservices 和 static 四种模板类型。

## 用法

```bash
docker compose init [OPTIONS] [TEMPLATE]
```

## 选项

| 选项 | 描述 |
|------|------|
| `-f`, `--file` | 指定要创建的 Compose 文件路径（默认：docker-compose.yml） |
| `-p`, `--project-name` | 指定项目名称 |
| `--force` | 强制覆盖现有文件 |
| `--quiet`, `-q` | 安静模式，减少输出信息 |
| `--help` | 显示帮助信息并退出 |

## 模板类型

| 模板 | 描述 | 包含的服务 |
|------|------|------------|
| `web` | Web 应用模板 | Nginx, PHP, MySQL |
| `api` | API 服务模板 | Node.js, MongoDB |
| `microservices` | 微服务模板 | API 网关, 认证服务, 业务服务, 数据库 |
| `static` | 静态网站模板 | Nginx |

## 使用示例

### 使用 web 模板初始化项目

```bash
docker compose init web
```

### 使用 api 模板初始化项目

```bash
docker compose init api
```

### 使用 microservices 模板初始化项目

```bash
docker compose init microservices
```

### 使用 static 模板初始化项目

```bash
docker compose init static
```

### 指定自定义 Compose 文件路径

```bash
docker compose init web -f docker-compose.custom.yml
```

### 指定项目名称

```bash
docker compose init api -p my-api-project
```

### 强制覆盖现有文件

```bash
docker compose init web --force
```

## 生成的文件结构

根据选择的模板，命令会生成以下文件结构：

### Web 模板

```
.
├── docker-compose.yml          # Compose 配置文件
├── nginx/                      # Nginx 配置
│   └── default.conf
├── php/                        # PHP 配置
│   └── Dockerfile
└── src/                        # 应用源码目录
    └── index.php
```

### API 模板

```
.
├── docker-compose.yml          # Compose 配置文件
├── nodejs/                     # Node.js 配置
│   ├── Dockerfile
│   ├── package.json
│   └── index.js
└── .env                        # 环境变量文件
```

### 微服务模板

```
.
├── docker-compose.yml          # Compose 配置文件
├── api-gateway/                # API 网关服务
│   ├── Dockerfile
│   └── config.yml
├── auth-service/               # 认证服务
│   ├── Dockerfile
│   └── app.js
├── business-service/           # 业务服务
│   ├── Dockerfile
│   └── app.js
└── .env                        # 环境变量文件
```

### 静态网站模板

```
.
├── docker-compose.yml          # Compose 配置文件
├── nginx/                      # Nginx 配置
│   └── default.conf
└── html/                       # 静态文件目录
    └── index.html
```

## 最佳实践

- **新项目初始化**：对于新项目，使用 `docker compose init` 快速搭建基础架构
- **模板选择**：根据项目类型选择合适的模板，减少配置工作量
- **自定义配置**：生成基础配置后，根据实际需求修改 Compose 文件
- **版本控制**：将生成的配置文件纳入版本控制系统，便于团队协作

## 注意事项

- 执行命令前确保当前目录为空，避免文件冲突
- 如需自定义服务配置，可在生成后编辑 Compose 文件
- 模板仅提供基础配置，实际项目可能需要根据具体需求进行调整

## 相关命令

- `docker compose up`：启动初始化后的服务
- `docker compose build`：构建服务镜像
- `docker compose config`：验证配置文件