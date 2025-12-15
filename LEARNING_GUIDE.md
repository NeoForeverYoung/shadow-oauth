# 📚 Shadow IAM 学习指南

## 🎯 学习目标

通过本指南，你将学会：
1. 理解 IAM 系统的基本架构
2. 掌握前后端分离的开发模式
3. 学会 JWT 认证的实现原理
4. 理解 RESTful API 的设计
5. 掌握现代 Web 开发的最佳实践

## 📖 学习路径（建议按顺序学习）

### 第一阶段：理解整体架构（30分钟）

#### 1. 先看项目结构
```
shadow-oauth/
├── backend/          # 后端服务（处理业务逻辑）
└── frontend/        # 前端应用（用户界面）
```

**理解要点**：
- 前后端分离：前端负责展示，后端负责数据处理
- 通过 HTTP API 通信

#### 2. 理解数据流向
```
用户操作 → 前端页面 → API 请求 → 后端处理 → 数据库 → 返回结果 → 前端显示
```

**推荐阅读顺序**：
1. 📄 `README.md` - 了解项目整体
2. 📄 `PROJECT_SUMMARY.md` - 查看完成的功能
3. 📄 `backend/README.md` - 后端说明
4. 📄 `frontend/README.md` - 前端说明

---

### 第二阶段：从简单开始 - 前端首页（20分钟）

#### 1. 查看首页代码
**文件**: `frontend/app/page.tsx`

**学习重点**：
```tsx
// 这是一个简单的 React 组件
export default function Home() {
  return (
    <div>...</div>  // JSX 语法，类似 HTML
  )
}
```

**理解要点**：
- React 组件就是返回 JSX 的函数
- Tailwind CSS 类名（如 `bg-blue-50`）用于样式
- Next.js 的 `Link` 组件用于页面跳转

**动手实践**：
1. 打开 `frontend/app/page.tsx`
2. 尝试修改文字内容
3. 尝试修改颜色（如 `bg-blue-50` 改为 `bg-red-50`）
4. 刷新浏览器看效果

---

### 第三阶段：理解用户注册流程（40分钟）

#### 1. 前端注册表单
**文件**: `frontend/components/RegisterForm.tsx`

**关键代码解析**：

```tsx
// 1. 状态管理 - 存储表单数据
const [formData, setFormData] = useState({
  email: '',
  password: '',
  // ...
});

// 2. 处理输入变化
const handleChange = (e) => {
  setFormData({
    ...formData,
    [e.target.name]: e.target.value,  // 更新对应字段
  });
};

// 3. 表单提交
const handleSubmit = async (e) => {
  e.preventDefault();  // 阻止默认提交行为
  
  // 调用 API
  const response = await register(
    formData.email, 
    formData.password
  );
  
  // 处理响应
  if (response.success) {
    router.push('/login');  // 跳转到登录页
  }
};
```

**学习重点**：
- `useState` - React 的状态管理 Hook
- `async/await` - 异步操作处理
- 表单验证逻辑
- API 调用方式

#### 2. API 调用层
**文件**: `frontend/lib/api.ts`

**关键代码**：
```typescript
// 创建 axios 实例（HTTP 客户端）
const apiClient = axios.create({
  baseURL: 'http://localhost:8080',  // 后端地址
});

// 注册函数
export const register = async (email, password, name) => {
  const response = await apiClient.post('/api/auth/register', {
    email,
    password,
    name,
  });
  return response.data;
};
```

**理解要点**：
- Axios 用于发送 HTTP 请求
- `POST` 方法用于创建数据
- 请求体包含用户信息

#### 3. 后端注册接口
**文件**: `backend/internal/handlers/auth.go`

**关键代码**：
```go
// 1. 接收请求
func (h *AuthHandler) Register(c *gin.Context) {
    var req service.RegisterRequest
    
    // 2. 解析 JSON 数据
    c.ShouldBindJSON(&req)
    
    // 3. 调用服务层处理业务逻辑
    user, err := h.authService.Register(req)
    
    // 4. 返回响应
    c.JSON(http.StatusCreated, response)
}
```

**学习重点**：
- Handler 层：处理 HTTP 请求和响应
- Gin 框架的上下文 `c *gin.Context`
- JSON 数据绑定

#### 4. 业务逻辑层
**文件**: `backend/internal/service/auth_service.go`

**关键代码**：
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
    
    // 4. 创建用户
    user := &models.User{
        Email:    req.Email,
        Password: string(hashedPassword),
    }
    database.DB.Create(user)
    
    return user, nil
}
```

**学习重点**：
- Service 层：核心业务逻辑
- 密码加密：使用 bcrypt
- 数据库操作：通过 GORM

#### 5. 数据模型
**文件**: `backend/internal/models/user.go`

**关键代码**：
```go
type User struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    Email     string    `gorm:"uniqueIndex" json:"email"`
    Password  string    `gorm:"not null" json:"-"`  // json:"-" 表示不返回
    CreatedAt time.Time `json:"created_at"`
}
```

**理解要点**：
- GORM 标签：定义数据库字段
- JSON 标签：定义 API 响应格式
- 结构体：Go 语言的数据类型

---

### 第四阶段：理解登录和 JWT 认证（50分钟）

#### 1. 登录流程
```
用户输入邮箱密码 
  → 前端发送 POST /api/auth/login
  → 后端验证密码
  → 生成 JWT Token
  → 返回 Token 给前端
  → 前端保存 Token
```

#### 2. JWT Token 生成
**文件**: `backend/internal/service/auth_service.go`

**关键代码**：
```go
func (s *AuthService) GenerateToken(userID uint) (string, error) {
    // 1. 创建 Claims（Token 内容）
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp": time.Now().Add(24 * time.Hour).Unix(),  // 过期时间
    }
    
    // 2. 创建 Token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // 3. 签名（使用密钥）
    tokenString, _ := token.SignedString([]byte(s.jwtSecret))
    
    return tokenString, nil
}
```

**理解要点**：
- JWT 包含三部分：Header、Payload、Signature
- Payload 包含用户信息（如 user_id）
- 使用密钥签名，防止篡改

#### 3. JWT 中间件
**文件**: `backend/internal/middleware/jwt.go`

**关键代码**：
```go
func JWTAuth(authService *service.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从 Header 获取 Token
        authHeader := c.GetHeader("Authorization")
        // 格式: "Bearer <token>"
        
        // 2. 验证 Token
        userID, err := authService.ValidateToken(tokenString)
        
        // 3. 将用户ID存入上下文
        c.Set("userID", userID)
        
        // 4. 继续处理请求
        c.Next()
    }
}
```

**学习重点**：
- 中间件：在请求处理前执行
- Token 验证：确保请求合法
- 上下文传递：将用户信息传给后续处理器

#### 4. 前端 Token 管理
**文件**: `frontend/lib/auth.ts`

**关键代码**：
```typescript
// 保存 Token
export const setToken = (token: string) => {
    localStorage.setItem('token', token);
};

// 获取 Token
export const getToken = (): string | null => {
    return localStorage.getItem('token');
};
```

**文件**: `frontend/lib/api.ts`

**关键代码**：
```typescript
// 请求拦截器：自动添加 Token
apiClient.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});
```

**理解要点**：
- localStorage：浏览器本地存储
- 拦截器：自动在每个请求中添加 Token
- Bearer Token：标准的认证方式

---

### 第五阶段：理解路由保护（30分钟）

#### 1. 前端路由保护
**文件**: `frontend/app/dashboard/page.tsx`

**关键代码**：
```tsx
useEffect(() => {
    // 检查是否已登录
    if (!isAuthenticated()) {
        router.push('/login');  // 未登录则跳转
        return;
    }
    
    // 获取用户信息
    const response = await getCurrentUser();
    setUser(response.data);
}, []);
```

**理解要点**：
- `useEffect`：组件加载时执行
- 检查认证状态
- 未登录自动跳转

#### 2. 后端路由保护
**文件**: `backend/cmd/server/main.go`

**关键代码**：
```go
// 公开接口（无需认证）
auth.POST("/register", authHandler.Register)
auth.POST("/login", authHandler.Login)

// 受保护接口（需要认证）
auth.GET("/me", 
    middleware.JWTAuth(authService),  // 中间件验证
    authHandler.GetCurrentUser
)
```

**理解要点**：
- 中间件链：按顺序执行
- 公开接口 vs 受保护接口
- 认证失败返回 401

---

## 🛠️ 实践建议

### 1. 按模块学习
不要一次性看完所有代码，按以下顺序：
1. ✅ 前端首页（最简单）
2. ✅ 注册功能（理解完整流程）
3. ✅ 登录功能（理解 JWT）
4. ✅ 路由保护（理解认证）

### 2. 动手实践
每学一个模块，尝试：
- 修改代码看效果
- 添加日志看执行流程
- 打断点调试
- 修改样式看 UI 变化

### 3. 画流程图
理解代码时，画出：
- 数据流向图
- 函数调用关系
- 组件层次结构

### 4. 阅读顺序建议

**第一天**（2小时）：
1. 阅读 `README.md`
2. 运行项目，体验功能
3. 查看前端首页代码
4. 理解注册表单组件

**第二天**（2小时）：
1. 理解 API 调用层
2. 理解后端 Handler
3. 理解 Service 层
4. 理解数据模型

**第三天**（2小时）：
1. 理解 JWT 认证原理
2. 理解中间件机制
3. 理解路由保护
4. 尝试修改和扩展功能

---

## 📝 关键概念解释

### 1. 前后端分离
- **前端**：负责用户界面，运行在浏览器
- **后端**：负责业务逻辑，运行在服务器
- **通信**：通过 HTTP API（JSON 格式）

### 2. RESTful API
- `GET`：获取数据
- `POST`：创建数据
- `PUT`：更新数据
- `DELETE`：删除数据

### 3. JWT 认证
- **Token**：类似"通行证"
- **生成**：登录成功后生成
- **验证**：每次请求时验证
- **过期**：设置过期时间，提高安全性

### 4. 中间件
- **作用**：在请求处理前/后执行
- **用途**：认证、日志、错误处理等
- **执行顺序**：按注册顺序执行

### 5. 状态管理
- **前端**：React Hooks (`useState`)
- **存储**：localStorage（持久化）
- **同步**：组件状态和存储同步

---

## 🎯 学习检查清单

完成以下任务，说明你已经理解了：

- [ ] 能解释前后端如何通信
- [ ] 能说出注册流程的每一步
- [ ] 能解释 JWT Token 的作用
- [ ] 能说出密码如何加密存储
- [ ] 能解释路由保护如何工作
- [ ] 能修改代码并看到效果
- [ ] 能添加新的 API 接口
- [ ] 能添加新的前端页面

---

## 💡 常见问题

### Q1: 为什么密码要加密？
**A**: 防止数据库泄露时密码被直接看到。bcrypt 是单向加密，无法解密。

### Q2: JWT Token 安全吗？
**A**: 相对安全，但要注意：
- 使用 HTTPS 传输
- 设置合理的过期时间
- 密钥要保密

### Q3: 前端如何知道用户已登录？
**A**: 检查 localStorage 中是否有 Token，有则已登录。

### Q4: 后端如何验证 Token？
**A**: 使用密钥验证 Token 签名，确保未被篡改。

### Q5: 为什么用中间件？
**A**: 避免在每个接口中重复写认证代码，提高代码复用性。

---

## 🚀 下一步学习

### 1. 扩展功能
- 添加"忘记密码"功能
- 添加"修改密码"功能
- 添加"用户资料编辑"功能

### 2. 深入学习
- 学习 OAuth 2.0 协议
- 学习 RBAC 权限模型
- 学习数据库设计
- 学习系统安全

### 3. 参考资源
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [Next.js 文档](https://nextjs.org/docs)
- [JWT 官方文档](https://jwt.io/)
- [GORM 文档](https://gorm.io/docs/)

---

## 📚 推荐学习路径

```
基础概念（1天）
  ↓
前端代码（2天）
  ↓
后端代码（2天）
  ↓
完整流程（1天）
  ↓
扩展功能（持续）
```

**总时间**：约 1 周可以完全理解

---

## 🎓 学习技巧

1. **不要死记硬背**：理解原理比记住代码更重要
2. **多动手实践**：修改代码看效果
3. **画图理解**：用流程图帮助理解
4. **循序渐进**：从简单到复杂
5. **多问为什么**：理解设计的原因

---

**祝你学习愉快！有问题随时查看代码注释或文档。** 🎉

