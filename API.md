# API文档

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