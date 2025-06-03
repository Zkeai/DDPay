# DDPay

DDPay 是一个支付处理系统。

## Docker 部署说明

### 前提条件

- 安装 Docker 和 Docker Compose
- 确保配置文件正确（etc/config.yaml）

### 环境变量配置

在部署前，您需要设置环境变量。我们提供了两种方式：

1. 使用配置向导（推荐）：

```bash
./setup-env.sh
```

按照提示输入相关配置信息，脚本会自动生成`.env`文件。

2. 手动配置：

```bash
cp .env.example .env
```

然后编辑`.env`文件，填写相关配置信息。

### 快速部署

1. 构建并启动所有服务

```bash
docker-compose up -d
```

2. 仅启动 HTTP 服务

```bash
docker-compose up -d http-server
```

3. 仅启动定时任务服务

```bash
docker-compose up -d cron-server
```

4. 仅启动钱包服务

```bash
docker-compose up -d wallet-server
```

### 查看日志

```bash
# 查看HTTP服务日志
docker-compose logs -f http-server

# 查看定时任务服务日志
docker-compose logs -f cron-server

# 查看钱包服务日志
docker-compose logs -f wallet-server
```

### 停止服务

```bash
docker-compose down
```

### 配置说明

配置文件位于`etc/config.yaml`，请确保以下配置正确：

- 数据库连接信息
- Redis 连接信息
- 钱包助记词
- Telegram Bot Token

### 注意事项

1. 配置文件中的敏感信息（如数据库密码、助记词等）建议使用环境变量替代
2. 生产环境部署时，请确保网络安全设置
3. 定期备份数据
4. 数据库初始化 SQL 脚本请放在`sql/`目录下，将在首次启动 MySQL 容器时自动执行

## JWT 鉴权和用户模块

本项目实现了基于 JWT 的鉴权功能和完整的用户模块，包括：

### JWT 鉴权功能

- 基于 JWT 的身份验证
- 使用 Redis 存储令牌，支持令牌撤销
- 令牌过期自动处理
- 支持令牌刷新

### 用户模块

- 邮箱验证码注册
- 邮箱密码登录
- OAuth 第三方登录（支持 Auth.js）
- 密码重置
- 用户信息管理
- 登录日志记录

### 数据库表结构

- users：用户表
- verification_codes：验证码表
- oauth_accounts：OAuth 账号关联表
- login_logs：登录日志表

### 配置说明

在`etc/config.yaml`中添加以下配置：

```yaml
jwt:
  secret: "your-secret-key-change-this-in-production"
  issuer: "DDPay"
  expirationTime: 24h
  refreshTokenTime: 168h
```

### API 接口

- POST /api/v1/user/register - 用户注册
- POST /api/v1/user/login - 用户登录
- POST /api/v1/user/oauth/login - OAuth 登录
- POST /api/v1/user/send-code - 发送验证码
- POST /api/v1/user/reset-password - 重置密码
- GET /api/v1/user/profile - 获取用户信息（需要认证）
- PUT /api/v1/user/profile - 更新用户信息（需要认证）
- POST /api/v1/user/logout - 注销登录（需要认证）
