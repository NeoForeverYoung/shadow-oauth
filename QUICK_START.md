# ⚡ 5分钟快速上手

## 🎯 目标
快速理解这个项目是做什么的，以及代码是如何组织的。

## 📍 你现在在哪里？

你点击了首页的"注册新账户"按钮，这个按钮的代码在：
```
frontend/app/page.tsx (第 20-25 行)
```

## 🔍 让我们追踪一下点击后的流程

### 第1步：点击按钮 → 跳转到注册页

**文件**: `frontend/app/page.tsx`
```tsx
<Link href="/register">注册新账户</Link>
```
👉 这行代码让浏览器跳转到 `/register` 页面

---

### 第2步：注册页面显示

**文件**: `frontend/app/register/page.tsx`

这个页面做了什么？
1. 显示标题"创建新账户"
2. 显示注册表单组件 `<RegisterForm />`

---

### 第3步：填写表单并提交

**文件**: `frontend/components/RegisterForm.tsx`

**关键代码**（第 60-75 行）：
```tsx
const handleSubmit = async (e) => {
  e.preventDefault();  // 阻止页面刷新
  
  // 调用注册 API
  const response = await register(
    formData.email, 
    formData.password, 
    formData.name
  );
  
  // 注册成功后跳转到登录页
  if (response.success) {
    router.push('/login');
  }
};
```

**理解**：
- 用户点击"注册"按钮
- 前端收集表单数据
- 调用 `register()` 函数发送请求到后端
- 后端处理完成后返回结果
- 前端根据结果显示成功或错误

---

### 第4步：API 请求发送

**文件**: `frontend/lib/api.ts`（第 50-58 行）

```typescript
export const register = async (email, password, name) => {
  const response = await apiClient.post('/api/auth/register', {
    email,
    password,
    name,
  });
  return response.data;
};
```

**理解**：
- 使用 `axios` 发送 HTTP POST 请求
- 请求地址：`http://localhost:8080/api/auth/register`
- 请求体：包含邮箱、密码、用户名

---

### 第5步：后端接收请求

**文件**: `backend/internal/handlers/auth.go`（第 30-50 行）

```go
func (h *AuthHandler) Register(c *gin.Context) {
    // 1. 解析请求数据
    var req service.RegisterRequest
    c.ShouldBindJSON(&req)
    
    // 2. 调用服务层处理
    user, err := h.authService.Register(req)
    
    // 3. 返回响应
    c.JSON(http.StatusCreated, response)
}
```

**理解**：
- `Register` 函数处理注册请求
- 从请求中提取数据
- 调用业务逻辑层处理
- 返回 JSON 响应

---

### 第6步：业务逻辑处理

**文件**: `backend/internal/service/auth_service.go`（第 50-90 行）

```go
func (s *AuthService) Register(req RegisterRequest) (*models.User, error) {
    // 1. 验证邮箱格式
    if !emailRegex.MatchString(req.Email) {
        return nil, ErrInvalidEmail
    }
    
    // 2. 检查邮箱是否已存在
    // ...
    
    // 3. 加密密码
    hashedPassword, _ := bcrypt.GenerateFromPassword(
        []byte(req.Password), 
        bcrypt.DefaultCost
    )
    
    // 4. 保存到数据库
    user := &models.User{
        Email:    req.Email,
        Password: string(hashedPassword),
    }
    database.DB.Create(user)
    
    return user, nil
}
```

**理解**：
- 验证输入数据
- 加密密码（使用 bcrypt）
- 保存到数据库
- 返回创建的用户

---

### 第7步：返回结果

数据流回前端：
```
后端 → JSON 响应 → 前端接收 → 显示结果
```

---

## 🎨 代码结构一目了然

```
用户操作
  ↓
前端页面 (React 组件)
  ↓
API 调用 (axios)
  ↓
HTTP 请求
  ↓
后端路由 (Gin)
  ↓
处理器 (Handler)
  ↓
业务逻辑 (Service)
  ↓
数据库 (GORM)
```

## 🧩 关键文件速查表

| 功能 | 前端文件 | 后端文件 |
|------|---------|---------|
| 首页 | `frontend/app/page.tsx` | - |
| 注册页 | `frontend/app/register/page.tsx` | - |
| 注册表单 | `frontend/components/RegisterForm.tsx` | - |
| API 调用 | `frontend/lib/api.ts` | - |
| 注册接口 | - | `backend/internal/handlers/auth.go` |
| 注册逻辑 | - | `backend/internal/service/auth_service.go` |
| 用户模型 | - | `backend/internal/models/user.go` |

## 💡 学习建议

### 1. 先看前端（更容易理解）
- 从 `frontend/app/page.tsx` 开始
- 看组件如何组织
- 看数据如何流动

### 2. 再看后端（理解业务逻辑）
- 从 `backend/cmd/server/main.go` 开始
- 看路由如何配置
- 看请求如何处理

### 3. 最后看完整流程
- 追踪一个完整的功能（如注册）
- 从前端到后端，理解每一步

## 🎯 现在试试这个

1. **打开** `frontend/app/page.tsx`
2. **找到** "注册新账户" 按钮
3. **修改** 按钮文字（如改为"立即注册"）
4. **保存** 文件
5. **刷新** 浏览器
6. **看到** 变化了吗？这就是前端开发！

## 📚 下一步

- 想深入学习？看 `LEARNING_GUIDE.md`
- 想运行项目？看 `START.md`
- 想了解功能？看 `README.md`

---

**记住**：代码不是用来背的，是用来理解的！🎉

