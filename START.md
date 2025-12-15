# 🚀 启动指南

## 第一次运行

### 终端 1：启动后端服务

```bash
cd backend
go run cmd/server/main.go
```

您应该看到类似以下输出：
```
配置加载成功，服务器端口: 8080
数据库连接成功: ./data/shadow.db
数据库表结构迁移成功
🚀 服务器启动在 http://localhost:8080
```

### 终端 2：启动前端服务

```bash
cd frontend
npm run dev
```

您应该看到类似以下输出：
```
  ▲ Next.js 15.1.3
  - Local:        http://localhost:3000

 ✓ Starting...
 ✓ Ready in 2.3s
```

## 测试应用

### 方法 1：通过浏览器测试（推荐）

1. **打开应用首页**
   - 访问：http://localhost:3000
   - 您会看到一个漂亮的欢迎页面

2. **注册新用户**
   - 点击"注册新账户"按钮
   - 填写以下信息：
     - 用户名（可选）：张三
     - 邮箱：test@example.com
     - 密码：123456
     - 确认密码：123456
   - 点击"注册"
   - 注册成功后会自动跳转到登录页

3. **用户登录**
   - 输入刚才注册的邮箱和密码
   - 点击"登录"
   - 登录成功后会跳转到仪表盘

4. **查看仪表盘**
   - 在仪表盘可以看到：
     - 欢迎信息
     - 用户的详细信息（ID、邮箱、注册时间等）
     - 已实现的功能列表

5. **退出登录**
   - 点击右上角的"退出登录"按钮
   - 会清除登录状态并跳转到登录页

6. **测试路由保护**
   - 在未登录状态下，尝试直接访问：http://localhost:3000/dashboard
   - 应该会自动跳转到登录页

### 方法 2：通过 API 测试

**健康检查**：
```bash
curl http://localhost:8080/health
```

**注册用户**：
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456",
    "name": "测试用户"
  }'
```

**用户登录**（会返回 JWT Token）：
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }'
```

**获取当前用户信息**（需要替换 YOUR_JWT_TOKEN）：
```bash
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 常见问题

### 1. 后端启动失败

**问题**：提示找不到包
**解决**：
```bash
cd backend
go mod tidy
go mod download
```

### 2. 前端启动失败

**问题**：提示依赖未安装
**解决**：
```bash
cd frontend
npm install
```

### 3. 数据库连接错误

**问题**：无法创建数据库文件
**解决**：确保 `backend/data/` 目录存在或有写入权限
```bash
cd backend
mkdir -p data
```

### 4. CORS 错误

**问题**：前端无法调用后端 API
**解决**：检查后端 CORS 配置是否允许 `http://localhost:3000`

### 5. 端口被占用

**问题**：8080 或 3000 端口已被占用
**解决**：
- 后端：设置环境变量 `PORT=8081`
- 前端：修改 `frontend/.env.local` 中的端口

## 数据库位置

SQLite 数据库文件位于：`backend/data/shadow.db`

如果需要重置数据，删除此文件并重启后端服务即可。

## 下一步

恭喜！您已经成功运行了 Shadow IAM 的第一阶段功能。

接下来可以：
- 📖 查看 README.md 了解更多功能
- 🔧 自定义配置（JWT 密钥、Token 过期时间等）
- 🚀 开始实现第二阶段功能（OAuth 2.0、RBAC 等）
- 💡 根据您的需求扩展功能

祝您使用愉快！🎉

