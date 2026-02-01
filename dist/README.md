# docker-compose-extensions v1.0.0

## 安装说明

### macOS 用户
1. 下载 `docker-compose-v1.0.0-darwin-amd64` 文件
2. 将文件重命名为 `docker-compose` 并移动到系统目录：
   ```bash
   mv ~/Downloads/docker-compose-v1.0.0-darwin-amd64 /usr/local/bin/docker-compose
   chmod +x /usr/local/bin/docker-compose
   ```
3. 验证安装：
   ```bash
   docker-compose --version
   ```

### Linux 用户
1. 克隆仓库并构建：
   ```bash
   git clone https://github.com/luyang668899/docker-compose-extensions.git
   cd docker-compose-extensions
   GOOS=linux GOARCH=amd64 go build -o docker-compose ./cmd
   sudo mv docker-compose /usr/local/bin/
   ```

### Windows 用户
1. 克隆仓库并构建：
   ```bash
   git clone https://github.com/luyang668899/docker-compose-extensions.git
   cd docker-compose-extensions
   GOOS=windows GOARCH=amd64 go build -o docker-compose.exe ./cmd
   ```
2. 将 `docker-compose.exe` 文件添加到系统 PATH 中

## 使用说明

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

### 扩展功能
- **dev**: 开发环境优化，支持热重载和代码同步
- **test**: 测试自动化，支持测试发现和覆盖率分析
- **sync**: 代码同步工具，支持双向/单向同步
- **perf**: 性能分析工具，分析资源使用并生成优化建议
- **share**: 环境分享工具，支持生成可分享链接

## 系统要求
- Docker 引擎 20.10.0 或更高版本
- macOS 10.15+ / Linux / Windows 10+

## 许可证
Apache-2.0 许可证
