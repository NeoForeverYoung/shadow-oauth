# Shadow OAuth 后端服务

基于 Golang + Gin + GORM 的 IAM 认证系统后端服务。

## 运行方式

```bash
# 启动服务器
go run cmd/server/main.go

# 或编译后运行
go build -o shadow-oauth cmd/server/main.go
./shadow-oauth
```

## 环境变量

- `PORT` - 服务器端口（默认：8080）
- `JWT_SECRET` - JWT签名密钥（默认：your-secret-key-change-in-production）
- `DATABASE_PATH` - 数据库文件路径（默认：./data/shadow.db）
- `JWT_EXPIRE_HOURS` - Token过期时间/小时（默认：24）

## API 接口

### 健康检查
```
GET /health
```

### 认证相关
```
POST /api/auth/register  - 用户注册
POST /api/auth/login     - 用户登录
GET  /api/auth/me        - 获取当前用户信息（需要认证）
```

