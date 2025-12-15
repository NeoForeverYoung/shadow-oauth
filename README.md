# 🔐 Shadow IAM - 身份认证系统

一个参考 Casdoor 设计的现代化 IAM (Identity and Access Management) 系统，实现了核心的用户注册、登录和认证功能。

## ✨ 已实现功能

### 第一阶段：基础认证功能 ✅

- ✅ **用户注册**：邮箱密码注册，密码 bcrypt 加密存储
- ✅ **用户登录**：JWT Token 认证机制
- ✅ **用户信息管理**：获取和显示当前用户信息
- ✅ **路由保护**：未登录用户自动跳转到登录页
- ✅ **会话管理**：Token 自动续期和过期处理
- ✅ **响应式 UI**：现代化的前端界面设计

## 🏗️ 技术栈

### 后端
- **语言**: Golang 1.24
- **框架**: Gin (HTTP 框架)
- **ORM**: GORM
- **数据库**: SQLite (开发环境)
- **认证**: JWT (golang-jwt/jwt/v5)
- **密码加密**: bcrypt

### 前端
- **框架**: Next.js 15 (App Router)
- **语言**: TypeScript
- **样式**: Tailwind CSS
- **HTTP 客户端**: Axios
- **状态管理**: React Hooks + LocalStorage

## 📦 项目结构

```
shadow-oauth/
├── backend/                    # 后端服务
│   ├── cmd/server/            # 主入口
│   ├── internal/              # 内部包
│   │   ├── models/            # 数据模型
│   │   ├── handlers/          # HTTP 处理器
│   │   ├── middleware/        # 中间件
│   │   ├── service/           # 业务逻辑层
│   │   └── database/          # 数据库连接
│   ├── config/                # 配置管理
│   └── go.mod
└── frontend/                   # 前端应用
    ├── app/                   # Next.js 页面
    ├── components/            # React 组件
    ├── lib/                   # 工具函数
    └── package.json
```

## 🚀 快速开始

### 前置要求

- Go 1.24+
- Node.js 20+
- npm 或 yarn

### 1. 启动后端服务

```bash
cd backend
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 2. 启动前端服务

```bash
cd frontend
npm run dev
```

前端应用将在 `http://localhost:3000` 启动。

### 3. 访问应用

打开浏览器访问 `http://localhost:3000`，您将看到：

- **首页** (`/`) - 欢迎页面，提供登录和注册入口
- **注册页** (`/register`) - 创建新账户
- **登录页** (`/login`) - 登录现有账户
- **仪表盘** (`/dashboard`) - 用户个人中心（需要登录）

## 🧪 测试流程

### 完整的端到端测试

1. **注册新用户**
   - 访问 `http://localhost:3000/register`
   - 填写邮箱、密码、确认密码
   - 点击"注册"按钮
   - 注册成功后自动跳转到登录页

2. **用户登录**
   - 在登录页输入刚注册的邮箱和密码
   - 点击"登录"按钮
   - 登录成功后跳转到仪表盘

3. **访问受保护资源**
   - 登录后可以访问仪表盘查看个人信息
   - 尝试在未登录状态访问 `/dashboard`，会自动跳转到登录页

4. **退出登录**
   - 在仪表盘点击"退出登录"按钮
   - Token 被清除，跳转回登录页

## 🔧 配置说明

### 后端配置

通过环境变量配置：

```bash
# 服务器端口
export PORT=8080

# JWT 密钥（生产环境必须修改）
export JWT_SECRET=your-secret-key-change-in-production

# JWT 过期时间（小时）
export JWT_EXPIRE_HOURS=24

# 数据库文件路径
export DATABASE_PATH=./data/shadow.db
```

### 前端配置

创建 `frontend/.env.local` 文件：

```
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

## 📡 API 接口

### 公开接口

- `GET /health` - 健康检查
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录

### 受保护接口（需要 JWT Token）

- `GET /api/auth/me` - 获取当前用户信息

### 请求示例

**注册**:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User"
  }'
```

**登录**:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**获取用户信息**:
```bash
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🎯 后续扩展方向

完成第一阶段后，可以逐步添加：

- 🔐 **OAuth 2.0 授权服务器**：实现标准的 OAuth 2.0 授权码流程
- 👥 **RBAC 权限管理**：角色和权限控制系统
- 🏢 **多租户支持**：组织和租户管理
- 🔑 **多因素认证 (MFA)**：TOTP、短信验证等
- 🌐 **第三方登录**：Google、GitHub、微信等社交登录
- 📧 **邮箱验证**：注册邮箱验证功能
- 🔄 **密码重置**：忘记密码找回功能
- 📊 **审计日志**：用户操作日志记录

## 📝 代码特点

- ✅ **详细的中文注释**：每个函数、变量都有清晰的中文说明
- ✅ **清晰的分层架构**：Handler → Service → Model，职责分明
- ✅ **完善的错误处理**：统一的错误响应格式
- ✅ **安全性保障**：密码加密、JWT 认证、CORS 配置
- ✅ **现代化 UI**：响应式设计，美观易用
- ✅ **类型安全**：TypeScript 提供完整的类型定义

## 📄 许可证

本项目采用 MIT 许可证。

## 🙏 致谢

本项目参考了 [Casdoor](https://github.com/casdoor/casdoor) 的设计理念和架构。

---

**Made with ❤️ | Golang + Gin + Next.js**

