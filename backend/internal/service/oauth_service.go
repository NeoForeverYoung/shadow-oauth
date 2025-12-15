package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/database"
	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var (
	// ErrInvalidClient 无效的客户端
	ErrInvalidClient = errors.New("无效的客户端ID或密钥")
	// ErrInvalidRedirectURI 无效的重定向URI
	ErrInvalidRedirectURI = errors.New("重定向URI不匹配")
	// ErrInvalidAuthorizationCode 无效的授权码
	ErrInvalidAuthorizationCode = errors.New("无效或已过期的授权码")
	// ErrAuthorizationCodeUsed 授权码已使用
	ErrAuthorizationCodeUsed = errors.New("授权码已被使用")
)

// OAuthService OAuth 服务
type OAuthService struct {
	jwtSecret string        // JWT 签名密钥
	jwtExpire time.Duration // Token 过期时间
}

// NewOAuthService 创建 OAuth 服务实例
func NewOAuthService(jwtSecret string, expireHours int) *OAuthService {
	return &OAuthService{
		jwtSecret: jwtSecret,
		jwtExpire: time.Duration(expireHours) * time.Hour,
	}
}

// ValidateClient 验证客户端（简化版：只验证 client_id 和 client_secret）
func (s *OAuthService) ValidateClient(clientID, clientSecret string) (*models.OAuthClient, error) {
	var client models.OAuthClient
	
	// 查询客户端
	if err := database.DB.Where("client_id = ? AND client_secret = ?", clientID, clientSecret).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidClient
		}
		return nil, fmt.Errorf("查询客户端失败: %w", err)
	}
	
	return &client, nil
}

// ValidateClientID 只验证客户端ID（用于授权页面，不需要密钥）
func (s *OAuthService) ValidateClientID(clientID string) (*models.OAuthClient, error) {
	var client models.OAuthClient
	
	if err := database.DB.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidClient
		}
		return nil, fmt.Errorf("查询客户端失败: %w", err)
	}
	
	return &client, nil
}

// ValidateRedirectURI 验证重定向URI是否匹配
func (s *OAuthService) ValidateRedirectURI(client *models.OAuthClient, redirectURI string) error {
	if client.RedirectURI != redirectURI {
		return ErrInvalidRedirectURI
	}
	return nil
}

// GenerateAuthorizationCode 生成授权码
// 授权码是临时的一次性令牌，用于交换 Access Token
func (s *OAuthService) GenerateAuthorizationCode(clientID string, userID uint, redirectURI string) (string, error) {
	// 1. 生成随机授权码（32字节，64个十六进制字符）
	codeBytes := make([]byte, 32)
	if _, err := rand.Read(codeBytes); err != nil {
		return "", fmt.Errorf("生成授权码失败: %w", err)
	}
	code := hex.EncodeToString(codeBytes)
	
	// 2. 创建授权码记录（有效期10分钟）
	authCode := &models.AuthorizationCode{
		Code:        code,
		ClientID:    clientID,
		UserID:      userID,
		RedirectURI: redirectURI,
		ExpiresAt:   time.Now().Add(10 * time.Minute), // 10分钟过期
		Used:        false,
	}
	
	// 3. 保存到数据库
	if err := database.DB.Create(authCode).Error; err != nil {
		return "", fmt.Errorf("保存授权码失败: %w", err)
	}
	
	return code, nil
}

// ExchangeAuthorizationCode 用授权码交换 Access Token
// 这是 OAuth 2.0 的核心步骤：授权码 → Access Token
func (s *OAuthService) ExchangeAuthorizationCode(code, clientID, clientSecret, redirectURI string) (*models.AccessToken, error) {
	// 1. 验证客户端
	client, err := s.ValidateClient(clientID, clientSecret)
	if err != nil {
		return nil, err
	}
	
	// 2. 验证重定向URI
	if err := s.ValidateRedirectURI(client, redirectURI); err != nil {
		return nil, err
	}
	
	// 3. 查找授权码
	var authCode models.AuthorizationCode
	if err := database.DB.Where("code = ? AND client_id = ?", code, clientID).First(&authCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidAuthorizationCode
		}
		return nil, fmt.Errorf("查询授权码失败: %w", err)
	}
	
	// 4. 检查授权码是否有效
	if !authCode.IsValid() {
		if authCode.Used {
			return nil, ErrAuthorizationCodeUsed
		}
		return nil, ErrInvalidAuthorizationCode
	}
	
	// 5. 标记授权码为已使用（授权码只能使用一次）
	authCode.Used = true
	database.DB.Save(&authCode)
	
	// 6. 生成 Access Token（JWT格式）
	tokenString, err := s.GenerateAccessToken(authCode.UserID, clientID)
	if err != nil {
		return nil, fmt.Errorf("生成 Access Token 失败: %w", err)
	}
	
	// 7. 保存 Access Token 到数据库
	accessToken := &models.AccessToken{
		Token:     tokenString,
		ClientID:  clientID,
		UserID:    authCode.UserID,
		ExpiresAt: time.Now().Add(s.jwtExpire),
	}
	
	if err := database.DB.Create(accessToken).Error; err != nil {
		return nil, fmt.Errorf("保存 Access Token 失败: %w", err)
	}
	
	return accessToken, nil
}

// GenerateAccessToken 生成 Access Token（JWT格式）
func (s *OAuthService) GenerateAccessToken(userID uint, clientID string) (string, error) {
	// 创建 JWT Claims
	claims := jwt.MapClaims{
		"user_id":  userID,                                    // 用户ID
		"client_id": clientID,                                // 客户端ID
		"exp":      time.Now().Add(s.jwtExpire).Unix(),       // 过期时间
		"iat":      time.Now().Unix(),                         // 签发时间
		"type":     "oauth_access_token",                       // Token类型标识
	}
	
	// 创建并签名 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// ValidateAccessToken 验证 Access Token
func (s *OAuthService) ValidateAccessToken(tokenString string) (uint, string, error) {
	// 解析 Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	
	if err != nil {
		return 0, "", fmt.Errorf("Token 解析失败: %w", err)
	}
	
	// 提取 Claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 检查 Token 类型
		if tokenType, ok := claims["type"].(string); !ok || tokenType != "oauth_access_token" {
			return 0, "", errors.New("无效的 Token 类型")
		}
		
		// 提取用户ID和客户端ID
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, "", errors.New("Token 中缺少 user_id")
		}
		
		clientID, ok := claims["client_id"].(string)
		if !ok {
			return 0, "", errors.New("Token 中缺少 client_id")
		}
		
		return uint(userID), clientID, nil
	}
	
	return 0, "", errors.New("无效的 Token")
}

// GetUserInfo 使用 Access Token 获取用户信息
// 这是 OAuth 的典型用法：第三方应用使用 Token 访问用户资源
func (s *OAuthService) GetUserInfo(tokenString string) (*models.User, error) {
	// 1. 验证 Token
	userID, _, err := s.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}
	
	// 2. 查询用户信息
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	
	return &user, nil
}

