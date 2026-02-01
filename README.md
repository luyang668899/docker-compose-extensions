# 开发环境指南

## 系统要求

- **Go** 1.25+
- **Docker** 20.10+
- **Git**
- **Make**

## 安装依赖

1. **安装Go**
   - macOS: `brew install go`
   - Linux: 参考官方文档安装
   - Windows: 下载安装包安装

2. **安装Docker**
   - 下载并安装Docker Desktop或Docker Engine

3. **安装make**
   - macOS: `brew install make`
   - Linux: 通常已安装
   - Windows: 可通过Chocolatey或WSL安装

## 构建命令

### 基本构建
```bash
# 构建项目
make build

# 构建并安装到Docker CLI插件目录
make install

# 跨平台构建
make cross
```

### 测试命令
```bash
# 运行单元测试
make test

# 运行端到端测试
make e2e

# 运行特定测试
make e2e-compose E2E_TEST=TestName
```

### 文档命令
```bash
# 生成文档
make docs

# 验证文档
make validate-docs
```

### 代码质量命令
```bash
# 格式化代码
make fmt

# 运行lint检查
make lint

# 验证依赖
make validate-go-mod
```

## 开发工作流

1. **设置环境**
   - 克隆仓库
   - 安装依赖
   - 构建项目

2. **编写代码**
   - 在`cmd/`目录下添加新命令
   - 在`pkg/`目录下添加核心功能
   - 遵循项目代码结构

3. **测试代码**
   - 编写单元测试
   - 运行端到端测试
   - 确保测试覆盖关键功能

4. **提交代码**
   - 格式化代码
   - 运行lint检查
   - 提交符合规范的代码

## 调试

1. **启用调试模式**
   ```bash
   go run ./cmd --debug
   ```

2. **使用Delve调试器**
   ```bash
   dlv debug ./cmd
   ```

3. **Docker调试**
   ```bash
   make build
   docker-compose --verbose
   ```

## 常见问题

### 构建失败
- 检查Go版本是否符合要求
- 确保Docker守护进程正在运行
- 检查网络连接，确保能下载依赖

### 测试失败
- 检查Docker环境是否正常
- 确保端口没有被占用
- 查看详细的测试日志

### 代码风格问题
- 运行`make fmt`自动修复
- 运行`make lint`查看详细错误
- # API文档

## 命令行接口

### 基本命令
```bash
# 启动服务
docker-compose up

# 查看状态
docker-compose ps

# 停止服务
docker-compose down

# 构建服务
docker-compose build

# 查看日志
docker-compose logs
```

### 扩展命令

#### dev（开发环境优化）
```bash
# 启动开发模式
docker-compose dev

# 选项
--watch: 监视文件变化
--sync: 自动同步代码
--hot-reload: 启用热重载
```

#### test（测试自动化）
```bash
# 运行测试
docker-compose test

# 选项
--watch: 监视测试文件变化
--coverage: 生成覆盖率报告
--format: 测试报告格式（json, xml, text）
```

#### sync（代码同步工具）
```bash
# 同步代码
docker-compose sync

# 选项
--direction: 同步方向（bidirectional, one-way）
--source: 源目录
--target: 目标目录
--watch: 监视模式
```

#### perf（性能分析）
```bash
# 分析性能
docker-compose perf

# 选项
--cpu: 分析CPU使用
--memory: 分析内存使用
--output: 输出格式（json, text）
--duration: 分析持续时间
```

#### share（环境分享）
```bash
# 分享环境
docker-compose share

# 选项
--expires: 过期时间（如1h, 1d）
--access: 访问控制（public, private）
--password: 设置访问密码
```

## 配置文件格式

### docker-compose.yml 扩展字段

#### x-dev（开发环境配置）
```yaml
x-dev:
  watch: true
  sync:
    - source: ./src
      target: /app/src
  hot-reload: true
```

#### x-test（测试配置）
```yaml
x-test:
  coverage: true
  format: json
  reports: ./test-reports
```

#### x-perf（性能配置）
```yaml
x-perf:
  cpu: true
  memory: true
  duration: 30s
```

## 环境变量

### 开发相关
- `DOCKER_COMPOSE_DEV_WATCH`: 启用文件监视
- `DOCKER_COMPOSE_DEV_SYNC`: 启用代码同步
- `DOCKER_COMPOSE_HOT_RELOAD`: 启用热重载

### 测试相关
- `DOCKER_COMPOSE_TEST_COVERAGE`: 生成覆盖率报告
- `DOCKER_COMPOSE_TEST_FORMAT`: 测试报告格式
- `DOCKER_COMPOSE_TEST_REPORTS`: 测试报告目录

### 性能相关
- `DOCKER_COMPOSE_PERF_DURATION`: 性能分析持续时间
- `DOCKER_COMPOSE_PERF_OUTPUT`: 性能分析输出格式

### 分享相关
- `DOCKER_COMPOSE_SHARE_EXPIRES`: 分享链接过期时间
- `DOCKER_COMPOSE_SHARE_ACCESS`: 分享访问控制
- `DOCKER_COMPOSE_SHARE_PASSWORD`: 分享访问密码
