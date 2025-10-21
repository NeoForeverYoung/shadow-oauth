# 从零构建 IAM/SSO 认证系统实战指南

> **项目名称**：MiniAuth（简化版 Casdoor）  
> **技术栈**：Go + Gin + GORM + React + MySQL  
> **难度级别**：中高级  
> **预计时间**：6-8周

---

## 🎯 项目目标

构建一个包含核心功能的认证授权系统，实现：
- ✅ 用户注册、登录
- ✅ JWT Token 认证
- ✅ OAuth 2.0 授权码模式
- ✅ RBAC 权限控制
- ✅ 第三方登录（GitHub）
- ✅ 管理后台

---

## 📐 整体架构设计

### 技术选型理由

| 组件 | 技术选择 | 理由 |
|------|---------|------|
| 后端框架 | Gin | 比 Beego 更轻量，性能更好，更现代 |
| ORM | GORM | 活跃度高，文档完善，易上手 |
| 数据库 | MySQL | 关系型数据，成熟稳定 |
| 缓存 | Redis | Session 存储、Token 黑名单 |
| 前端 | React + Ant Design | 组件丰富，开发效率高 |
| 权限 | Casbin | 成熟的权限框架 |

### 项目结构

```
miniauth/
├── backend/                 # 后端代码
│   ├── cmd/
│   │   └── server/
│   │       └── main.go     # 程序入口
│   ├── internal/
│   │   ├── handler/        # HTTP 处理器
│   │   ├── service/        # 业务逻辑
│   │   ├── repository/     # 数据访问层
│   │   ├── model/          # 数据模型
│   │   ├── middleware/     # 中间件
│   │   └── pkg/            # 工具包
│   ├── config/             # 配置文件
│   ├── migrations/         # 数据库迁移
│   └── go.mod
├── frontend/               # 前端代码
│   ├── src/
│   │   ├── pages/         # 页面
│   │   ├── components/    # 组件
│   │   ├── services/      # API 调用
│   │   └── utils/         # 工具函数
│   └── package.json
├── docs/                   # 文档
└── README.md
```

---

## 🚀 迭代开发计划

---

## 第一阶段：基础框架搭建（Week 1）

### 目标：搭建项目骨架，实现基础用户管理

### Day 1-2：项目初始化

#### 创建项目结构
```bash
mkdir -p miniauth/{backend,frontend,docs}
cd miniauth/backend

# 初始化 Go 模块
go mod init github.com/yourusername/miniauth

# 创建目录结构
mkdir -p cmd/server internal/{handler,service,repository,model,middleware,pkg} config migrations
```

#### 安装核心依赖
```bash
# 后端依赖
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/redis/go-redis/v9
go get -u github.com/spf13/viper
go get -u golang.org/x/crypto/bcrypt
```

#### 创建配置文件

**config/config.yaml**
```yaml
server:
  port: 8080
  mode: debug  # debug, release

database:
  host: localhost
  port: 3306
  user: root
  password: "123456"
  dbname: miniauth
  charset: utf8mb4

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key-change-in-production"
  expire_hours: 24
  refresh_expire_hours: 168  # 7天

app:
  name: "MiniAuth"
  version: "1.0.0"
```

#### 配置加载代码

**internal/pkg/config/config.go**
```go
package config

import (
    "github.com/spf13/viper"
    "log"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    App      AppConfig      `mapstructure:"app"`
}

type ServerConfig struct {
    Port int    `mapstructure:"port"`
    Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"dbname"`
    Charset  string `mapstructure:"charset"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
    Secret            string `mapstructure:"secret"`
    ExpireHours       int    `mapstructure:"expire_hours"`
    RefreshExpireHours int   `mapstructure:"refresh_expire_hours"`
}

type AppConfig struct {
    Name    string `mapstructure:"name"`
    Version string `mapstructure:"version"`
}

var GlobalConfig *Config

func LoadConfig(configPath string) error {
    viper.SetConfigFile(configPath)
    viper.SetConfigType("yaml")
    
    if err := viper.ReadInConfig(); err != nil {
        return err
    }
    
    GlobalConfig = &Config{}
    if err := viper.Unmarshal(GlobalConfig); err != nil {
        return err
    }
    
    log.Println("配置加载成功")
    return nil
}
```

### Day 3-4：数据库模型设计

#### 核心数据模型

**internal/model/user.go**
```go
package model

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    
    // 基础字段
    Username     string `gorm:"uniqueIndex;size:50;not null" json:"username"`
    Password     string `gorm:"size:255;not null" json:"-"` // 不返回给前端
    Email        string `gorm:"uniqueIndex;size:100" json:"email"`
    Phone        string `gorm:"size:20" json:"phone"`
    
    // 个人信息
    DisplayName  string `gorm:"size:100" json:"display_name"`
    Avatar       string `gorm:"size:500" json:"avatar"`
    Bio          string `gorm:"type:text" json:"bio"`
    
    // 状态
    Status       string `gorm:"size:20;default:'active'" json:"status"` // active, disabled, locked
    IsAdmin      bool   `gorm:"default:false" json:"is_admin"`
    EmailVerified bool  `gorm:"default:false" json:"email_verified"`
    
    // 组织相关
    OrganizationID uint `gorm:"index" json:"organization_id"`
    Organization   *Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
    
    // 关联
    Roles        []Role        `gorm:"many2many:user_roles;" json:"roles,omitempty"`
    OAuth        []OAuthLink   `gorm:"foreignKey:UserID" json:"-"`
}

type Organization struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    
    Name        string `gorm:"uniqueIndex;size:100;not null" json:"name"`
    DisplayName string `gorm:"size:200" json:"display_name"`
    Description string `gorm:"type:text" json:"description"`
    Logo        string `gorm:"size:500" json:"logo"`
    Website     string `gorm:"size:200" json:"website"`
    
    Users []User `gorm:"foreignKey:OrganizationID" json:"-"`
}

type Role struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    
    Name        string `gorm:"uniqueIndex;size:50;not null" json:"name"`
    DisplayName string `gorm:"size:100" json:"display_name"`
    Description string `gorm:"type:text" json:"description"`
    
    OrganizationID uint `gorm:"index" json:"organization_id"`
    
    Users []User `gorm:"many2many:user_roles;" json:"-"`
}

type Permission struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    
    Subject  string `gorm:"size:100;not null" json:"subject"` // user, role
    Object   string `gorm:"size:200;not null" json:"object"`  // resource
    Action   string `gorm:"size:50;not null" json:"action"`   // read, write, delete
    Effect   string `gorm:"size:20;default:'allow'" json:"effect"` // allow, deny
    
    RoleID uint `gorm:"index" json:"role_id"`
    Role   *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

type Application struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    
    Name         string `gorm:"uniqueIndex;size:100;not null" json:"name"`
    DisplayName  string `gorm:"size:200" json:"display_name"`
    Description  string `gorm:"type:text" json:"description"`
    Logo         string `gorm:"size:500" json:"logo"`
    
    // OAuth 配置
    ClientID     string `gorm:"uniqueIndex;size:100;not null" json:"client_id"`
    ClientSecret string `gorm:"size:255;not null" json:"client_secret"`
    RedirectURIs string `gorm:"type:text" json:"redirect_uris"` // JSON 数组
    GrantTypes   string `gorm:"size:500" json:"grant_types"`    // authorization_code,refresh_token
    
    OrganizationID uint `gorm:"index" json:"organization_id"`
    Organization   *Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}

type OAuthLink struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    
    UserID   uint   `gorm:"index;not null" json:"user_id"`
    Provider string `gorm:"size:50;not null" json:"provider"` // github, google, etc
    ProviderUserID string `gorm:"size:200;not null" json:"provider_user_id"`
    AccessToken    string `gorm:"type:text" json:"-"`
    RefreshToken   string `gorm:"type:text" json:"-"`
    
    User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Token struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    
    UserID       uint      `gorm:"index;not null" json:"user_id"`
    ApplicationID uint     `gorm:"index" json:"application_id"`
    
    AccessToken  string    `gorm:"uniqueIndex;size:500;not null" json:"access_token"`
    RefreshToken string    `gorm:"uniqueIndex;size:500" json:"refresh_token"`
    TokenType    string    `gorm:"size:20;default:'Bearer'" json:"token_type"`
    ExpiresAt    time.Time `gorm:"index" json:"expires_at"`
    Scope        string    `gorm:"size:500" json:"scope"`
    
    User        *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Application *Application `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
}
```

#### 数据库初始化

**internal/pkg/database/mysql.go**
```go
package database

import (
    "fmt"
    "log"
    "time"
    
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    
    "github.com/yourusername/miniauth/internal/model"
    "github.com/yourusername/miniauth/internal/pkg/config"
)

var DB *gorm.DB

func InitMySQL() error {
    cfg := config.GlobalConfig.Database
    
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
        cfg.User,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.DBName,
        cfg.Charset,
    )
    
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().Local()
        },
    })
    
    if err != nil {
        return fmt.Errorf("连接数据库失败: %v", err)
    }
    
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    
    // 连接池配置
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    log.Println("数据库连接成功")
    return nil
}

func AutoMigrate() error {
    return DB.AutoMigrate(
        &model.User{},
        &model.Organization{},
        &model.Role{},
        &model.Permission{},
        &model.Application{},
        &model.OAuthLink{},
        &model.Token{},
    )
}
```

### Day 5-7：实现基础认证功能

#### 密码加密工具

**internal/pkg/auth/password.go**
```go
package auth

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

#### JWT Token 生成与验证

**internal/pkg/auth/jwt.go**
```go
package auth

import (
    "errors"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/yourusername/miniauth/internal/pkg/config"
)

type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    IsAdmin  bool   `json:"is_admin"`
    jwt.RegisteredClaims
}

// GenerateToken 生成访问令牌
func GenerateToken(userID uint, username string, isAdmin bool) (string, error) {
    cfg := config.GlobalConfig.JWT
    
    claims := Claims{
        UserID:   userID,
        Username: username,
        IsAdmin:  isAdmin,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(cfg.ExpireHours))),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    config.GlobalConfig.App.Name,
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.Secret))
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID uint, username string, isAdmin bool) (string, error) {
    cfg := config.GlobalConfig.JWT
    
    claims := Claims{
        UserID:   userID,
        Username: username,
        IsAdmin:  isAdmin,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(cfg.RefreshExpireHours))),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    config.GlobalConfig.App.Name,
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析令牌
func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(config.GlobalConfig.JWT.Secret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("无效的令牌")
}
```

#### 用户服务层

**internal/service/user_service.go**
```go
package service

import (
    "errors"
    "github.com/yourusername/miniauth/internal/model"
    "github.com/yourusername/miniauth/internal/repository"
    "github.com/yourusername/miniauth/internal/pkg/auth"
)

type UserService struct {
    userRepo *repository.UserRepository
}

func NewUserService() *UserService {
    return &UserService{
        userRepo: repository.NewUserRepository(),
    }
}

// Register 用户注册
func (s *UserService) Register(username, password, email string) (*model.User, error) {
    // 检查用户是否已存在
    existingUser, _ := s.userRepo.FindByUsername(username)
    if existingUser != nil {
        return nil, errors.New("用户名已存在")
    }
    
    if email != "" {
        existingEmail, _ := s.userRepo.FindByEmail(email)
        if existingEmail != nil {
            return nil, errors.New("邮箱已被使用")
        }
    }
    
    // 加密密码
    hashedPassword, err := auth.HashPassword(password)
    if err != nil {
        return nil, err
    }
    
    // 创建用户
    user := &model.User{
        Username:    username,
        Password:    hashedPassword,
        Email:       email,
        DisplayName: username,
        Status:      "active",
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    return user, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*model.User, string, string, error) {
    // 查找用户
    user, err := s.userRepo.FindByUsername(username)
    if err != nil {
        return nil, "", "", errors.New("用户名或密码错误")
    }
    
    // 验证密码
    if !auth.CheckPassword(password, user.Password) {
        return nil, "", "", errors.New("用户名或密码错误")
    }
    
    // 检查用户状态
    if user.Status != "active" {
        return nil, "", "", errors.New("账号已被禁用")
    }
    
    // 生成 Token
    accessToken, err := auth.GenerateToken(user.ID, user.Username, user.IsAdmin)
    if err != nil {
        return nil, "", "", err
    }
    
    refreshToken, err := auth.GenerateRefreshToken(user.ID, user.Username, user.IsAdmin)
    if err != nil {
        return nil, "", "", err
    }
    
    return user, accessToken, refreshToken, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
    return s.userRepo.FindByID(id)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *model.User) error {
    return s.userRepo.Update(user)
}
```

#### 数据访问层

**internal/repository/user_repository.go**
```go
package repository

import (
    "github.com/yourusername/miniauth/internal/model"
    "github.com/yourusername/miniauth/internal/pkg/database"
    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository() *UserRepository {
    return &UserRepository{
        db: database.DB,
    }
}

func (r *UserRepository) Create(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
    var user model.User
    err := r.db.Preload("Organization").Preload("Roles").First(&user, id).Error
    return &user, err
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
    var user model.User
    err := r.db.Where("username = ?", username).First(&user).Error
    return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
    var user model.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}

func (r *UserRepository) Update(user *model.User) error {
    return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
    return r.db.Delete(&model.User{}, id).Error
}

func (r *UserRepository) List(page, pageSize int) ([]*model.User, int64, error) {
    var users []*model.User
    var total int64
    
    offset := (page - 1) * pageSize
    
    err := r.db.Model(&model.User{}).Count(&total).Error
    if err != nil {
        return nil, 0, err
    }
    
    err = r.db.Offset(offset).Limit(pageSize).
        Preload("Organization").
        Find(&users).Error
    
    return users, total, err
}
```

#### HTTP 处理器

**internal/handler/auth_handler.go**
```go
package handler

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/yourusername/miniauth/internal/service"
)

type AuthHandler struct {
    userService *service.UserService
}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{
        userService: service.NewUserService(),
    }
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Password string `json:"password" binding:"required,min=6"`
    Email    string `json:"email" binding:"email"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    User         interface{} `json:"user"`
    AccessToken  string      `json:"access_token"`
    RefreshToken string      `json:"refresh_token"`
    TokenType    string      `json:"token_type"`
}

// Register 注册
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.userService.Register(req.Username, req.Password, req.Email)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "注册成功",
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        },
    })
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, accessToken, refreshToken, err := h.userService.Login(req.Username, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, LoginResponse{
        User: gin.H{
            "id":           user.ID,
            "username":     user.Username,
            "email":        user.Email,
            "display_name": user.DisplayName,
            "avatar":       user.Avatar,
            "is_admin":     user.IsAdmin,
        },
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
    })
}

// GetProfile 获取当前用户信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
    userID, _ := c.Get("user_id")
    
    user, err := h.userService.GetUserByID(userID.(uint))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "user": gin.H{
            "id":           user.ID,
            "username":     user.Username,
            "email":        user.Email,
            "display_name": user.DisplayName,
            "avatar":       user.Avatar,
            "is_admin":     user.IsAdmin,
            "status":       user.Status,
        },
    })
}
```

#### 认证中间件

**internal/middleware/auth.go**
```go
package middleware

import (
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/yourusername/miniauth/internal/pkg/auth"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证令牌"})
            c.Abort()
            return
        }
        
        // Bearer token
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
            c.Abort()
            return
        }
        
        tokenString := parts[1]
        claims, err := auth.ParseToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
            c.Abort()
            return
        }
        
        // 将用户信息存入上下文
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("is_admin", claims.IsAdmin)
        
        c.Next()
    }
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        isAdmin, exists := c.Get("is_admin")
        if !exists || !isAdmin.(bool) {
            c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

#### 主程序入口

**cmd/server/main.go**
```go
package main

import (
    "fmt"
    "log"
    
    "github.com/gin-gonic/gin"
    "github.com/yourusername/miniauth/internal/handler"
    "github.com/yourusername/miniauth/internal/middleware"
    "github.com/yourusername/miniauth/internal/pkg/config"
    "github.com/yourusername/miniauth/internal/pkg/database"
)

func main() {
    // 加载配置
    if err := config.LoadConfig("config/config.yaml"); err != nil {
        log.Fatal("配置加载失败:", err)
    }
    
    // 初始化数据库
    if err := database.InitMySQL(); err != nil {
        log.Fatal("数据库初始化失败:", err)
    }
    
    // 自动迁移数据表
    if err := database.AutoMigrate(); err != nil {
        log.Fatal("数据表迁移失败:", err)
    }
    
    // 设置 Gin 模式
    gin.SetMode(config.GlobalConfig.Server.Mode)
    
    // 创建路由
    r := gin.Default()
    
    // 跨域中间件
    r.Use(middleware.CORSMiddleware())
    
    // 初始化处理器
    authHandler := handler.NewAuthHandler()
    
    // 公开路由
    public := r.Group("/api/v1")
    {
        public.POST("/register", authHandler.Register)
        public.POST("/login", authHandler.Login)
    }
    
    // 需要认证的路由
    protected := r.Group("/api/v1")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/profile", authHandler.GetProfile)
    }
    
    // 管理员路由
    admin := r.Group("/api/v1/admin")
    admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
    {
        // 管理员接口
    }
    
    // 启动服务器
    addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
    log.Printf("服务器启动在 %s\n", addr)
    if err := r.Run(addr); err != nil {
        log.Fatal("服务器启动失败:", err)
    }
}
```

#### CORS 中间件

**internal/middleware/cors.go**
```go
package middleware

import (
    "github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

### 第一阶段验收标准 ✅

- [ ] 项目结构清晰，可以正常编译运行
- [ ] 数据库表自动创建成功
- [ ] 可以通过 API 注册新用户
- [ ] 可以使用用户名密码登录并获取 JWT Token
- [ ] 可以使用 Token 访问受保护的 API
- [ ] 使用 Postman 测试所有 API 正常工作

---

## 第二阶段：OAuth 2.0 授权服务器（Week 2）

### 目标：实现完整的 OAuth 2.0 授权码模式

### Day 1-2：OAuth 数据模型与授权流程

#### OAuth 授权码模型

**internal/model/oauth.go**
```go
package model

import (
    "time"
    "gorm.io/gorm"
)

type OAuthAuthorizationCode struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    
    Code        string    `gorm:"uniqueIndex;size:500;not null" json:"code"`
    UserID      uint      `gorm:"index;not null" json:"user_id"`
    ClientID    string    `gorm:"index;size:100;not null" json:"client_id"`
    RedirectURI string    `gorm:"type:text;not null" json:"redirect_uri"`
    Scope       string    `gorm:"size:500" json:"scope"`
    ExpiresAt   time.Time `gorm:"index" json:"expires_at"`
    Used        bool      `gorm:"default:false" json:"used"`
    
    User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (OAuthAuthorizationCode) TableName() string {
    return "oauth_authorization_codes"
}
```

#### OAuth 服务

**internal/service/oauth_service.go**
```go
package service

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "strings"
    "time"
    
    "github.com/yourusername/miniauth/internal/model"
    "github.com/yourusername/miniauth/internal/repository"
    "github.com/yourusername/miniauth/internal/pkg/auth"
)

type OAuthService struct {
    appRepo  *repository.ApplicationRepository
    codeRepo *repository.OAuthCodeRepository
    tokenRepo *repository.TokenRepository
}

func NewOAuthService() *OAuthService {
    return &OAuthService{
        appRepo:  repository.NewApplicationRepository(),
        codeRepo: repository.NewOAuthCodeRepository(),
        tokenRepo: repository.NewTokenRepository(),
    }
}

// 生成授权码
func (s *OAuthService) GenerateAuthorizationCode(
    userID uint,
    clientID, redirectURI, scope string,
) (string, error) {
    // 验证应用
    app, err := s.appRepo.FindByClientID(clientID)
    if err != nil {
        return "", errors.New("无效的客户端")
    }
    
    // 验证 redirect_uri
    if !s.validateRedirectURI(app.RedirectURIs, redirectURI) {
        return "", errors.New("无效的回调地址")
    }
    
    // 生成授权码
    code := generateRandomString(32)
    
    authCode := &model.OAuthAuthorizationCode{
        Code:        code,
        UserID:      userID,
        ClientID:    clientID,
        RedirectURI: redirectURI,
        Scope:       scope,
        ExpiresAt:   time.Now().Add(10 * time.Minute), // 10分钟有效期
        Used:        false,
    }
    
    if err := s.codeRepo.Create(authCode); err != nil {
        return "", err
    }
    
    return code, nil
}

// 用授权码换取 Token
func (s *OAuthService) ExchangeToken(
    code, clientID, clientSecret, redirectURI string,
) (accessToken, refreshToken string, expiresIn int, err error) {
    // 验证客户端
    app, err := s.appRepo.FindByClientID(clientID)
    if err != nil || app.ClientSecret != clientSecret {
        return "", "", 0, errors.New("客户端认证失败")
    }
    
    // 查找授权码
    authCode, err := s.codeRepo.FindByCode(code)
    if err != nil {
        return "", "", 0, errors.New("无效的授权码")
    }
    
    // 检查授权码是否已使用
    if authCode.Used {
        return "", "", 0, errors.New("授权码已使用")
    }
    
    // 检查是否过期
    if time.Now().After(authCode.ExpiresAt) {
        return "", "", 0, errors.New("授权码已过期")
    }
    
    // 验证 redirect_uri
    if authCode.RedirectURI != redirectURI {
        return "", "", 0, errors.New("回调地址不匹配")
    }
    
    // 标记授权码为已使用
    authCode.Used = true
    s.codeRepo.Update(authCode)
    
    // 生成访问令牌
    accessToken = generateRandomString(64)
    refreshToken = generateRandomString(64)
    expiresIn = 3600 * 24 // 24小时
    
    // 保存 Token
    token := &model.Token{
        UserID:        authCode.UserID,
        ApplicationID: app.ID,
        AccessToken:   accessToken,
        RefreshToken:  refreshToken,
        TokenType:     "Bearer",
        ExpiresAt:     time.Now().Add(time.Duration(expiresIn) * time.Second),
        Scope:         authCode.Scope,
    }
    
    if err := s.tokenRepo.Create(token); err != nil {
        return "", "", 0, err
    }
    
    return accessToken, refreshToken, expiresIn, nil
}

// 验证回调地址
func (s *OAuthService) validateRedirectURI(allowedURIs, uri string) bool {
    uris := strings.Split(allowedURIs, ",")
    for _, allowed := range uris {
        if strings.TrimSpace(allowed) == uri {
            return true
        }
    }
    return false
}

func generateRandomString(length int) string {
    b := make([]byte, length)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
```

### Day 3-4：OAuth API 实现

**internal/handler/oauth_handler.go**
```go
package handler

import (
    "net/http"
    "net/url"
    
    "github.com/gin-gonic/gin"
    "github.com/yourusername/miniauth/internal/service"
)

type OAuthHandler struct {
    oauthService *service.OAuthService
    userService  *service.UserService
}

func NewOAuthHandler() *OAuthHandler {
    return &OAuthHandler{
        oauthService: service.NewOAuthService(),
        userService:  service.NewUserService(),
    }
}

// Authorize 授权端点
// GET /oauth/authorize?response_type=code&client_id=xxx&redirect_uri=xxx&scope=xxx&state=xxx
func (h *OAuthHandler) Authorize(c *gin.Context) {
    responseType := c.Query("response_type")
    clientID := c.Query("client_id")
    redirectURI := c.Query("redirect_uri")
    scope := c.Query("scope")
    state := c.Query("state")
    
    // 参数验证
    if responseType != "code" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_response_type"})
        return
    }
    
    if clientID == "" || redirectURI == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
        return
    }
    
    // 检查用户是否已登录（从中间件获取）
    userID, exists := c.Get("user_id")
    if !exists {
        // 重定向到登录页面
        loginURL := "/login?redirect=" + url.QueryEscape(c.Request.URL.String())
        c.Redirect(http.StatusFound, loginURL)
        return
    }
    
    // 生成授权码
    code, err := h.oauthService.GenerateAuthorizationCode(
        userID.(uint),
        clientID,
        redirectURI,
        scope,
    )
    
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 重定向回客户端
    redirectURL, _ := url.Parse(redirectURI)
    query := redirectURL.Query()
    query.Set("code", code)
    if state != "" {
        query.Set("state", state)
    }
    redirectURL.RawQuery = query.Encode()
    
    c.Redirect(http.StatusFound, redirectURL.String())
}

type TokenRequest struct {
    GrantType    string `json:"grant_type" binding:"required"`
    Code         string `json:"code"`
    RedirectURI  string `json:"redirect_uri"`
    ClientID     string `json:"client_id" binding:"required"`
    ClientSecret string `json:"client_secret" binding:"required"`
    RefreshToken string `json:"refresh_token"`
}

// Token 令牌端点
// POST /oauth/token
func (h *OAuthHandler) Token(c *gin.Context) {
    var req TokenRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
        return
    }
    
    switch req.GrantType {
    case "authorization_code":
        h.handleAuthorizationCodeGrant(c, req)
    case "refresh_token":
        h.handleRefreshTokenGrant(c, req)
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_grant_type"})
    }
}

func (h *OAuthHandler) handleAuthorizationCodeGrant(c *gin.Context, req TokenRequest) {
    if req.Code == "" || req.RedirectURI == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
        return
    }
    
    accessToken, refreshToken, expiresIn, err := h.oauthService.ExchangeToken(
        req.Code,
        req.ClientID,
        req.ClientSecret,
        req.RedirectURI,
    )
    
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "access_token":  accessToken,
        "token_type":    "Bearer",
        "expires_in":    expiresIn,
        "refresh_token": refreshToken,
    })
}

func (h *OAuthHandler) handleRefreshTokenGrant(c *gin.Context, req TokenRequest) {
    // TODO: 实现刷新令牌逻辑
    c.JSON(http.StatusNotImplemented, gin.H{"error": "not_implemented"})
}
```

### Day 5-7：完善 OAuth 功能

继续实现：
- [ ] 用户授权同意页面
- [ ] Token 刷新机制
- [ ] Token 撤销接口
- [ ] Scope 权限控制
- [ ] PKCE 支持（增强安全性）

### 第二阶段验收标准 ✅

- [ ] 完整实现 OAuth 2.0 授权码模式
- [ ] 可以创建和管理应用（Application）
- [ ] 授权流程完整可用
- [ ] 能够用授权码换取 Access Token
- [ ] 实现 Token 刷新功能
- [ ] 编写完整的 API 测试用例

---

## 第三阶段：权限控制系统（Week 3）

### 目标：基于 Casbin 实现 RBAC 权限控制

### Day 1-2：集成 Casbin

#### 安装 Casbin
```bash
go get github.com/casbin/casbin/v2
go get github.com/casbin/gorm-adapter/v3
```

#### Casbin 配置

**config/rbac_model.conf**
```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

#### 初始化 Casbin

**internal/pkg/casbin/enforcer.go**
```go
package casbin

import (
    "github.com/casbin/casbin/v2"
    gormadapter "github.com/casbin/gorm-adapter/v3"
    "github.com/yourusername/miniauth/internal/pkg/database"
)

var Enforcer *casbin.Enforcer

func InitCasbin() error {
    adapter, err := gormadapter.NewAdapterByDB(database.DB)
    if err != nil {
        return err
    }
    
    enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
    if err != nil {
        return err
    }
    
    // 加载策略
    enforcer.LoadPolicy()
    
    Enforcer = enforcer
    return nil
}

// 添加角色策略
func AddRolePolicy(role, resource, action string) error {
    _, err := Enforcer.AddPolicy(role, resource, action)
    return err
}

// 为用户分配角色
func AddRoleForUser(username, role string) error {
    _, err := Enforcer.AddGroupingPolicy(username, role)
    return err
}

// 检查权限
func CheckPermission(username, resource, action string) (bool, error) {
    return Enforcer.Enforce(username, resource, action)
}
```

### Day 3-5：实现角色和权限管理

**internal/service/role_service.go**
```go
package service

import (
    "github.com/yourusername/miniauth/internal/model"
    "github.com/yourusername/miniauth/internal/repository"
    pkgcasbin "github.com/yourusername/miniauth/internal/pkg/casbin"
)

type RoleService struct {
    roleRepo *repository.RoleRepository
}

func NewRoleService() *RoleService {
    return &RoleService{
        roleRepo: repository.NewRoleRepository(),
    }
}

func (s *RoleService) CreateRole(role *model.Role) error {
    return s.roleRepo.Create(role)
}

func (s *RoleService) AssignRoleToUser(username, roleName string) error {
    return pkgcasbin.AddRoleForUser(username, roleName)
}

func (s *RoleService) AddPermissionToRole(roleName, resource, action string) error {
    return pkgcasbin.AddRolePolicy(roleName, resource, action)
}

func (s *RoleService) CheckUserPermission(username, resource, action string) (bool, error) {
    return pkgcasbin.CheckPermission(username, resource, action)
}
```

### Day 6-7：权限中间件

**internal/middleware/permission.go**
```go
package middleware

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    pkgcasbin "github.com/yourusername/miniauth/internal/pkg/casbin"
)

// PermissionMiddleware 权限检查中间件
func PermissionMiddleware(resource, action string) gin.HandlerFunc {
    return func(c *gin.Context) {
        username, exists := c.Get("username")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
            c.Abort()
            return
        }
        
        ok, err := pkgcasbin.CheckPermission(username.(string), resource, action)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
            c.Abort()
            return
        }
        
        if !ok {
            c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### 第三阶段验收标准 ✅

- [ ] 成功集成 Casbin
- [ ] 可以创建和管理角色
- [ ] 可以为角色添加权限
- [ ] 可以为用户分配角色
- [ ] API 端点受权限保护
- [ ] 测试不同角色的访问控制

---

## 第四阶段：第三方登录（Week 4）

### 目标：实现 GitHub OAuth 登录

### Day 1-3：GitHub OAuth 集成

#### 配置 GitHub OAuth Provider

**internal/model/provider.go**
```go
package model

import (
    "time"
    "gorm.io/gorm"
)

type Provider struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    
    Name         string `gorm:"uniqueIndex;size:100;not null" json:"name"`
    DisplayName  string `gorm:"size:200" json:"display_name"`
    Type         string `gorm:"size:50;not null" json:"type"` // github, google, etc
    ClientID     string `gorm:"size:200;not null" json:"client_id"`
    ClientSecret string `gorm:"size:500;not null" json:"client_secret"`
    AuthURL      string `gorm:"size:500" json:"auth_url"`
    TokenURL     string `gorm:"size:500" json:"token_url"`
    UserInfoURL  string `gorm:"size:500" json:"user_info_url"`
    Scopes       string `gorm:"size:500" json:"scopes"`
    
    OrganizationID uint `gorm:"index" json:"organization_id"`
}
```

#### GitHub OAuth 服务

**internal/service/github_oauth.go**
```go
package service

import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/github"
)

type GitHubOAuthService struct {
    config *oauth2.Config
}

func NewGitHubOAuthService(clientID, clientSecret, redirectURL string) *GitHubOAuthService {
    config := &oauth2.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        RedirectURL:  redirectURL,
        Scopes:       []string{"user:email"},
        Endpoint:     github.Endpoint,
    }
    
    return &GitHubOAuthService{
        config: config,
    }
}

func (s *GitHubOAuthService) GetAuthURL(state string) string {
    return s.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (s *GitHubOAuthService) ExchangeToken(code string) (*oauth2.Token, error) {
    return s.config.Exchange(context.Background(), code)
}

type GitHubUser struct {
    ID        int    `json:"id"`
    Login     string `json:"login"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    AvatarURL string `json:"avatar_url"`
}

func (s *GitHubOAuthService) GetUserInfo(token *oauth2.Token) (*GitHubUser, error) {
    client := s.config.Client(context.Background(), token)
    
    resp, err := client.Get("https://api.github.com/user")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var user GitHubUser
    if err := json.Unmarshal(body, &user); err != nil {
        return nil, err
    }
    
    return &user, nil
}
```

### Day 4-7：完整实现第三方登录流程

- [ ] 实现 Provider 管理接口
- [ ] 实现第三方登录回调处理
- [ ] 账号关联和解绑功能
- [ ] 支持多个第三方账号绑定
- [ ] 新用户自动注册

### 第四阶段验收标准 ✅

- [ ] 成功配置 GitHub OAuth App
- [ ] 可以通过 GitHub 登录
- [ ] 首次登录自动创建账号
- [ ] 可以绑定和解绑 GitHub 账号
- [ ] 测试完整的第三方登录流程

---

## 第五阶段：前端开发（Week 5）

### 目标：实现管理后台前端

### Day 1-2：前端项目初始化

```bash
cd miniauth/frontend
npx create-react-app . --template typescript
npm install antd axios react-router-dom @types/react-router-dom
```

### Day 3-7：核心页面开发

#### 页面列表
1. 登录页面
2. 注册页面
3. 用户管理页面
4. 应用管理页面
5. 角色管理页面
6. 个人设置页面

#### 示例：登录页面

**src/pages/Login.tsx**
```typescript
import React, { useState } from 'react';
import { Form, Input, Button, Card, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { login } from '../services/auth';

const Login: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      const response = await login(values.username, values.password);
      localStorage.setItem('access_token', response.data.access_token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
      message.success('登录成功');
      navigate('/dashboard');
    } catch (error: any) {
      message.error(error.response?.data?.error || '登录失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ 
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center', 
      height: '100vh' 
    }}>
      <Card title="登录" style={{ width: 400 }}>
        <Form
          name="login"
          onFinish={onFinish}
          autoComplete="off"
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="用户名" />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password placeholder="密码" />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              登录
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default Login;
```

#### API 服务封装

**src/services/auth.ts**
```typescript
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
});

// 请求拦截器
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const login = (username: string, password: string) => {
  return api.post('/login', { username, password });
};

export const register = (username: string, password: string, email: string) => {
  return api.post('/register', { username, password, email });
};

export const getProfile = () => {
  return api.get('/profile');
};
```

### 第五阶段验收标准 ✅

- [ ] 前端项目可以正常运行
- [ ] 完成所有核心页面
- [ ] API 调用正常
- [ ] 路由配置完整
- [ ] UI 美观，用户体验良好

---

## 第六阶段：高级功能（Week 6-7）

### 可选功能模块

1. **多因素认证 (MFA)**
   - TOTP（Time-based One-Time Password）
   - 短信验证码
   - 邮箱验证码

2. **SAML 2.0 支持**
   - 作为 IdP（身份提供商）
   - SAML 断言生成

3. **LDAP 服务器**
   - 基础 LDAP 协议支持
   - 用户同步

4. **审计日志**
   - 操作日志记录
   - 日志查询和导出

5. **Webhook**
   - 事件通知机制
   - 自定义 Webhook

---

## 📋 开发最佳实践

### 代码规范
- 使用 `gofmt` 格式化代码
- 遵循 Go 命名规范
- 编写单元测试（目标覆盖率 >60%）
- 添加代码注释

### Git 提交规范
```
feat: 添加用户注册功能
fix: 修复登录 token 过期问题
docs: 更新 API 文档
test: 添加用户服务测试用例
refactor: 重构权限检查逻辑
```

### 测试策略
- 单元测试：每个 Service 层函数
- 集成测试：API 端点
- E2E 测试：核心业务流程

---

## 🎯 每周里程碑检查清单

### Week 1 ✅
- [ ] 项目结构搭建完成
- [ ] 数据库模型设计完成
- [ ] 基础认证功能实现
- [ ] JWT Token 机制正常工作

### Week 2 ✅
- [ ] OAuth 2.0 授权码模式实现
- [ ] 应用管理功能完成
- [ ] Token 管理完善

### Week 3 ✅
- [ ] Casbin 权限系统集成
- [ ] 角色和权限管理实现
- [ ] API 权限控制生效

### Week 4 ✅
- [ ] GitHub OAuth 登录实现
- [ ] 账号绑定功能完成
- [ ] 第三方登录流程测试通过

### Week 5 ✅
- [ ] 前端项目搭建完成
- [ ] 核心页面开发完成
- [ ] 前后端联调成功

### Week 6-7 ✅
- [ ] 至少实现一个高级功能
- [ ] 编写完整项目文档
- [ ] 性能优化和安全加固

---

## 📚 学习资源推荐

### 技术文档
- Gin 框架：https://gin-gonic.com/docs/
- GORM：https://gorm.io/docs/
- Casbin：https://casbin.org/docs/
- OAuth 2.0：https://oauth.net/2/
- JWT：https://jwt.io/introduction

### 开源项目参考
- Casdoor：https://github.com/casdoor/casdoor
- Dex：https://github.com/dexidp/dex
- Keycloak：https://www.keycloak.org/
- ORY Hydra：https://www.ory.sh/hydra/

### 在线工具
- JWT 调试：https://jwt.io
- OAuth Playground：https://www.oauth.com/playground/
- Postman：API 测试工具

---

## 🔧 故障排查指南

### 常见问题

**1. 数据库连接失败**
```go
// 检查配置文件
// 确认 MySQL 服务运行
// 检查防火墙设置
```

**2. JWT Token 验证失败**
```go
// 检查密钥配置
// 确认 token 未过期
// 验证 token 格式
```

**3. OAuth 回调失败**
```go
// 检查 redirect_uri 配置
// 确认应用配置正确
// 查看错误日志
```

---

## 🎓 项目总结与提升

### 完成项目后的学习目标

1. **技术深度**
   - 深入理解 OAuth 2.0 和 OIDC
   - 掌握 JWT 的安全实践
   - 理解 RBAC 权限模型

2. **工程能力**
   - 项目架构设计能力
   - 代码组织和模块化
   - 测试驱动开发

3. **安全意识**
   - 密码安全存储
   - Token 安全管理
   - API 安全防护

### 后续优化方向

1. **性能优化**
   - Redis 缓存
   - 数据库查询优化
   - 接口性能测试

2. **功能扩展**
   - 多租户支持
   - SSO 单点登录
   - 移动端 SDK

3. **运维部署**
   - Docker 容器化
   - Kubernetes 部署
   - 监控和日志

---

## 🤝 获取帮助

如果在开发过程中遇到问题：

1. 查阅项目文档和代码注释
2. 搜索相关技术问题
3. 参考 Casdoor 源码实现
4. 在技术社区提问
5. 随时向我寻求帮助！

---

## ✨ 结语

从零构建一个 IAM 系统是一个挑战，但也是最好的学习方式。通过这个项目，你将：

- 💪 深入理解认证授权机制
- 🚀 提升全栈开发能力
- 📚 积累实战项目经验
- 🎯 掌握核心安全知识

**祝你开发顺利！** 🎉

记住：代码不是一次写成的，而是不断迭代改进的。保持学习，持续优化！

