# 🔐 OAuth 2.0 简化版实现指南

## 📚 什么是 OAuth 2.0？

OAuth 2.0 是一个授权框架，允许第三方应用在用户授权后访问用户的资源，而无需获取用户的密码。

### 核心概念

1. **资源所有者（Resource Owner）**：用户
2. **客户端（Client）**：第三方应用
3. **授权服务器（Authorization Server）**：我们的 Shadow IAM 系统
4. **资源服务器（Resource Server）**：提供用户数据的服务器

### 授权码流程（Authorization Code Flow）

这是最安全的 OAuth 流程，我们实现的就是这个：

```
1. 客户端引导用户到授权服务器
   ↓
2. 用户登录并授权
   ↓
3. 授权服务器返回授权码（通过重定向）
   ↓
4. 客户端用授权码 + 客户端密钥交换 Access Token
   ↓
5. 客户端使用 Access Token 访问用户资源
```

## 🏗️ 实现架构

### 数据模型

1. **OAuthClient**：OAuth 客户端（第三方应用）
   - `client_id`：客户端标识符（公开）
   - `client_secret`：客户端密钥（保密）
   - `redirect_uri`：重定向地址

2. **AuthorizationCode**：授权码（临时令牌）
   - `code`：授权码字符串
   - `client_id`：关联的客户端
   - `user_id`：授权用户
   - `expires_at`：过期时间（10分钟）
   - `used`：是否已使用（授权码只能使用一次）

3. **AccessToken**：访问令牌
   - `token`：JWT Token
   - `client_id`：关联的客户端
   - `user_id`：授权用户
   - `expires_at`：过期时间

### API 端点

1. **GET /oauth/authorize**：授权端点
   - 参数：`client_id`, `redirect_uri`, `response_type=code`, `state`
   - 功能：生成授权码并重定向回客户端

2. **POST /oauth/token**：Token 端点
   - 参数：`grant_type=authorization_code`, `code`, `client_id`, `client_secret`, `redirect_uri`
   - 功能：用授权码交换 Access Token

3. **GET /oauth/userinfo**：用户信息端点
   - 参数：`access_token`（查询参数或 Header）
   - 功能：使用 Access Token 获取用户信息

## 🚀 快速开始

### 1. 初始化测试客户端

```bash
cd backend
go run cmd/init_oauth_client/main.go
```

这会创建一个测试客户端：
- Client ID: `test_client_123`
- Client Secret: `test_secret_456`
- Redirect URI: `http://localhost:3000/oauth/test-client/callback`

### 2. 启动服务

**后端**：
```bash
cd backend
go run cmd/server/main.go
```

**前端**：
```bash
cd frontend
npm run dev
```

### 3. 测试 OAuth 流程

1. 访问测试客户端页面：http://localhost:3000/oauth/test-client
2. 点击"开始授权"按钮
3. 如果未登录，会跳转到登录页
4. 登录后会自动跳转到授权页面
5. 点击"同意授权"
6. 会重定向回测试客户端，显示授权码
7. 点击"交换 Access Token"
8. 点击"获取用户信息"
9. 完成！✅

## 📖 代码学习路径

### 1. 理解数据模型

**文件**：
- `backend/internal/models/oauth_client.go` - OAuth 客户端模型
- `backend/internal/models/authorization_code.go` - 授权码模型
- `backend/internal/models/access_token.go` - Access Token 模型

**学习重点**：
- 为什么授权码要设置过期时间？
- 为什么授权码只能使用一次？
- Access Token 和授权码的区别是什么？

### 2. 理解服务层

**文件**：`backend/internal/service/oauth_service.go`

**关键方法**：
- `GenerateAuthorizationCode()` - 生成授权码
- `ExchangeAuthorizationCode()` - 用授权码交换 Token
- `ValidateAccessToken()` - 验证 Access Token

**学习重点**：
- 授权码是如何生成的？（随机字符串）
- Token 交换时做了哪些验证？
- 为什么需要验证 `redirect_uri`？

### 3. 理解处理器

**文件**：`backend/internal/handlers/oauth.go`

**关键端点**：
- `Authorize()` - 处理授权请求
- `Token()` - 处理 Token 交换
- `UserInfo()` - 返回用户信息

**学习重点**：
- 授权端点如何生成授权码？
- Token 端点如何验证授权码？
- 为什么需要验证客户端密钥？

### 4. 理解前端流程

**文件**：
- `frontend/app/oauth/test-client/page.tsx` - 测试客户端
- `frontend/app/oauth/authorize/page.tsx` - 授权页面

**学习重点**：
- 客户端如何构建授权 URL？
- 如何处理授权回调？
- 如何用授权码交换 Token？

## 🔍 关键概念详解

### 1. 授权码（Authorization Code）

**作用**：临时的一次性令牌，用于交换 Access Token

**特点**：
- 有效期短（10分钟）
- 只能使用一次
- 通过 HTTPS 传输
- 不包含用户信息

**为什么需要授权码？**
- 安全性：客户端密钥不会暴露给浏览器
- 一次性：防止重放攻击

### 2. Access Token

**作用**：用于访问受保护资源的令牌

**特点**：
- 有效期较长（24小时）
- 可以多次使用
- 包含用户和客户端信息（JWT）
- 需要妥善保管

### 3. 重定向 URI 验证

**为什么重要？**
- 防止授权码被劫持
- 确保授权码只能被正确的客户端使用

**验证逻辑**：
- 授权时的 `redirect_uri` 必须与客户端注册的 `redirect_uri` 一致
- Token 交换时也要验证 `redirect_uri`

### 4. State 参数

**作用**：防止 CSRF 攻击

**使用方式**：
- 客户端生成随机字符串
- 在授权请求中传递
- 授权服务器原样返回
- 客户端验证是否一致

## 🧪 测试场景

### 场景 1：正常流程

1. ✅ 用户已登录
2. ✅ 客户端信息正确
3. ✅ 重定向 URI 匹配
4. ✅ 授权成功

### 场景 2：未登录

1. ❌ 用户未登录
2. ✅ 自动跳转到登录页
3. ✅ 登录后返回授权页面

### 场景 3：无效客户端

1. ❌ Client ID 不存在
2. ✅ 返回错误信息

### 场景 4：授权码过期

1. ❌ 授权码超过 10 分钟
2. ✅ 返回"授权码已过期"

### 场景 5：授权码重复使用

1. ❌ 同一个授权码使用两次
2. ✅ 第二次返回"授权码已被使用"

## 🔐 安全考虑

### 已实现的安全措施

1. ✅ **HTTPS 传输**（生产环境必须）
2. ✅ **授权码一次性使用**
3. ✅ **授权码短期有效（10分钟）**
4. ✅ **重定向 URI 验证**
5. ✅ **客户端密钥验证**
6. ✅ **JWT Token 签名验证**

### 可以改进的地方

1. 🔄 **Refresh Token**：延长登录时间
2. 🔄 **PKCE**：增强移动端安全性
3. 🔄 **Scope 权限控制**：细粒度权限
4. 🔄 **Token 撤销机制**：允许用户撤销授权

## 📝 常见问题

### Q1: 为什么授权码不能直接访问资源？

**A**: 授权码是临时的一次性令牌，必须通过服务器端交换成 Access Token。这样可以：
- 保护客户端密钥（不会暴露给浏览器）
- 验证客户端身份
- 记录 Token 使用情况

### Q2: Access Token 和 JWT Token 有什么区别？

**A**: 在这个实现中，Access Token 就是 JWT Token。JWT 是一种 Token 格式，包含：
- Header：算法信息
- Payload：用户和客户端信息
- Signature：签名验证

### Q3: 为什么需要验证 redirect_uri？

**A**: 防止授权码被劫持。如果攻击者知道授权码，但没有正确的 redirect_uri，也无法交换 Token。

### Q4: State 参数是必须的吗？

**A**: 不是必须的，但强烈推荐。用于防止 CSRF 攻击。

## 🎯 下一步学习

完成这个简化版后，可以学习：

1. **Refresh Token**：延长登录时间
2. **Scope 权限**：细粒度权限控制
3. **PKCE**：增强移动端安全性
4. **Token 撤销**：允许用户撤销授权
5. **OpenID Connect**：在 OAuth 基础上添加身份信息

## 📚 参考资源

- [OAuth 2.0 官方规范](https://oauth.net/2/)
- [RFC 6749](https://tools.ietf.org/html/rfc6749)
- [JWT.io](https://jwt.io/) - JWT 调试工具

---

**祝你学习愉快！有问题随时查看代码注释。** 🎉

