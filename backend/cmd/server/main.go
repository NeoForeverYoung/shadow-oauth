package main

import (
	"log"
	"net/http"

	"github.com/NeoForeverYoung/shadow-oauth/backend/config"
	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/database"
	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/handlers"
	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/middleware"
	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/models"
	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	cfg := config.Load()
	log.Printf("配置加载成功，服务器端口: %s", cfg.Server.Port)

	// 2. 初始化数据库
	if err := database.Initialize(cfg.Database.Path); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.Close()

	// 3. 自动迁移数据库表结构
	if err := database.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 4. 初始化 Gin 路由
	router := setupRouter(cfg)

	// 5. 启动服务器
	addr := ":" + cfg.Server.Port
	log.Printf("🚀 服务器启动在 http://localhost%s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// setupRouter 配置路由和中间件
func setupRouter(cfg *config.Config) *gin.Engine {
	// 设置 Gin 模式（可通过环境变量 GIN_MODE=release 切换为生产模式）
	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// 配置 CORS（允许前端跨域访问）
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // 允许的前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 健康检查接口
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.SuccessResponse("服务运行正常", gin.H{
			"status": "healthy",
		}))
	})

	// 初始化服务层
	authService := service.NewAuthService(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(authService)

	// API 路由组
	api := router.Group("/api")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			// 公开接口（无需认证）
			auth.POST("/register", authHandler.Register) // 用户注册
			auth.POST("/login", authHandler.Login)       // 用户登录

			// 受保护接口（需要认证）
			auth.GET("/me", middleware.JWTAuth(authService), authHandler.GetCurrentUser) // 获取当前用户信息
		}
	}

	return router
}
