package main

import (
	"log"

	"github.com/NeoForeverYoung/shadow-oauth/backend/config"
	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/database"
	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/models"
)

func main() {
	// 1. 加载配置
	cfg := config.Load()

	// 2. 初始化数据库
	if err := database.Initialize(cfg.Database.Path); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.Close()

	// 3. 自动迁移数据库表结构
	if err := database.AutoMigrate(
		&models.OAuthClient{},
		&models.AuthorizationCode{},
		&models.AccessToken{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 4. 创建测试客户端
	testClient := &models.OAuthClient{
		ClientID:     "test_client_123",
		ClientSecret: "test_secret_456",
		Name:         "OAuth 测试客户端",
		RedirectURI:  "http://localhost:3000/oauth/test-client/callback",
	}

	// 检查是否已存在
	var existing models.OAuthClient
	if err := database.DB.Where("client_id = ?", testClient.ClientID).First(&existing).Error; err == nil {
		log.Printf("测试客户端已存在，跳过创建")
		log.Printf("Client ID: %s", existing.ClientID)
		log.Printf("Client Secret: %s", existing.ClientSecret)
		log.Printf("Redirect URI: %s", existing.RedirectURI)
		return
	}

	// 创建客户端
	if err := database.DB.Create(testClient).Error; err != nil {
		log.Fatalf("创建测试客户端失败: %v", err)
	}

	log.Println("✅ 测试客户端创建成功！")
	log.Printf("Client ID: %s", testClient.ClientID)
	log.Printf("Client Secret: %s", testClient.ClientSecret)
	log.Printf("Redirect URI: %s", testClient.RedirectURI)
	log.Println("\n💡 提示：这些信息用于 OAuth 测试，请妥善保管 Client Secret")
}

