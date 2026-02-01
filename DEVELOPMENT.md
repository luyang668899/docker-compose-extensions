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