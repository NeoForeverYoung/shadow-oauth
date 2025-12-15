# Shadow IAM 前端

基于 Next.js 15 + TypeScript + Tailwind CSS 的现代化身份认证系统前端。

## 运行方式

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 启动生产服务器
npm start
```

访问 http://localhost:3000

## 环境变量

创建 `.env.local` 文件：

```
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

## 功能特性

- ✅ 用户注册和登录
- ✅ JWT Token 认证
- ✅ 受保护的路由
- ✅ 响应式设计
- ✅ 表单验证
- ✅ 错误处理

## 技术栈

- **框架**: Next.js 15 (App Router)
- **语言**: TypeScript
- **样式**: Tailwind CSS
- **HTTP 客户端**: Axios
- **状态管理**: React Hooks + LocalStorage

