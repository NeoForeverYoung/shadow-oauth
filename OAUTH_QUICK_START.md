# ⚡ OAuth 2.0 快速开始

## 🎯 5分钟体验 OAuth 流程

### 步骤 1: 初始化测试客户端

```bash
cd backend
go run cmd/init_oauth_client/main.go
```

你会看到：
```
✅ 测试客户端创建成功！
Client ID: test_client_123
Client Secret: test_secret_456
Redirect URI: http://localhost:3000/oauth/test-client/callback
```

### 步骤 2: 启动服务

**终端 1 - 启动后端**：
```bash
cd backend
go run cmd/server/main.go
```

**终端 2 - 启动前端**：
```bash
cd frontend
npm run dev
```

### 步骤 3: 测试 OAuth 流程

1. **打开测试客户端页面**
   - 访问：http://localhost:3000/oauth/test-client

2. **点击"开始授权"**
   - 会跳转到授权服务器
   - 如果未登录，先登录

3. **同意授权**
   - 在授权页面点击"同意授权"
   - 会自动重定向回测试客户端

4. **交换 Access Token**
   - 看到授权码后，点击"交换 Access Token"
   - 会获得 Access Token

5. **获取用户信息**
   - 点击"获取用户信息"
   - 会显示你的用户信息

## 📊 完整流程演示

```
用户 → 测试客户端 → 授权服务器 → 用户登录 → 授权同意
  ↓
授权码返回 → 测试客户端 → 交换 Token → 获取用户信息
```

## 🔍 查看代码

### 后端关键文件

1. **数据模型**：
   - `backend/internal/models/oauth_client.go` - OAuth 客户端
   - `backend/internal/models/authorization_code.go` - 授权码
   - `backend/internal/models/access_token.go` - Access Token

2. **服务层**：
   - `backend/internal/service/oauth_service.go` - OAuth 业务逻辑

3. **处理器**：
   - `backend/internal/handlers/oauth.go` - OAuth API 端点

### 前端关键文件

1. **测试客户端**：
   - `frontend/app/oauth/test-client/page.tsx` - 完整的 OAuth 流程演示

2. **授权页面**：
   - `frontend/app/oauth/authorize/page.tsx` - 用户授权界面

## 🎓 学习要点

### 1. OAuth 2.0 授权码流程

```
客户端 → 授权服务器（用户登录+授权）→ 授权码
  ↓
客户端（用授权码+密钥）→ Token 端点 → Access Token
  ↓
客户端（用 Access Token）→ 资源服务器 → 用户信息
```

### 2. 关键概念

- **授权码**：临时的一次性令牌（10分钟有效）
- **Access Token**：用于访问资源的令牌（24小时有效）
- **客户端密钥**：验证客户端身份（保密）
- **重定向 URI**：授权后跳转的地址（必须匹配）

### 3. 安全机制

- ✅ 授权码只能使用一次
- ✅ 授权码短期有效（10分钟）
- ✅ 重定向 URI 必须匹配
- ✅ 客户端密钥验证
- ✅ JWT Token 签名验证

## 🐛 常见问题

### Q: 授权后没有返回授权码？

**A**: 检查：
1. 是否已登录
2. 客户端 ID 是否正确
3. 重定向 URI 是否匹配

### Q: 交换 Token 失败？

**A**: 检查：
1. 授权码是否过期（10分钟）
2. 授权码是否已使用
3. 客户端密钥是否正确
4. 重定向 URI 是否匹配

### Q: 获取用户信息失败？

**A**: 检查：
1. Access Token 是否有效
2. Token 是否过期
3. 是否正确传递 Token

## 📚 深入学习

查看 `OAUTH_GUIDE.md` 了解：
- OAuth 2.0 详细原理
- 代码实现细节
- 安全考虑
- 扩展方向

---

**现在就开始体验吧！** 🚀

