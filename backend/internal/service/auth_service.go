package service

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/database"
	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	// ErrInvalidEmail 无效的邮箱格式
	ErrInvalidEmail = errors.New("邮箱格式无效")
	// ErrWeakPassword 密码强度不足
	ErrWeakPassword = errors.New("密码必须至少包含 6 个字符")
	// ErrEmailExists 邮箱已存在
	ErrEmailExists = errors.New("该邮箱已被注册")
	// ErrInvalidCredentials 无效的登录凭证
	ErrInvalidCredentials = errors.New("邮箱或密码错误")
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("用户不存在")
)

// emailRegex 邮箱格式验证正则表达式
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// AuthService 认证服务
type AuthService struct {
	jwtSecret string        // JWT 签名密钥
	jwtExpire time.Duration // JWT 过期时间
}

// NewAuthService 创建认证服务实例
func NewAuthService(jwtSecret string, expireHours int) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
		jwtExpire: time.Duration(expireHours) * time.Hour,
	}
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`    // 邮箱（必填）
	Password string `json:"password" binding:"required"` // 密码（必填）
	Name     string `json:"name"`                        // 用户名（可选）
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`    // 邮箱（必填）
	Password string `json:"password" binding:"required"` // 密码（必填）
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token string              `json:"token"` // JWT Token
	User  models.UserResponse `json:"user"`  // 用户信息
}

// Register 用户注册
func (s *AuthService) Register(req RegisterRequest) (*models.User, error) {
	// 1. 验证邮箱格式
	if !emailRegex.MatchString(req.Email) {
		return nil, ErrInvalidEmail
	}

	// 2. 验证密码强度
	if len(req.Password) < 6 {
		return nil, ErrWeakPassword
	}

	// 3. 检查邮箱是否已存在
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, ErrEmailExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 4. 加密密码（使用 bcrypt）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 5. 创建用户
	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	// 1. 查询用户
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 2. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 3. 生成 JWT Token
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("生成 Token 失败: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// GenerateToken 生成 JWT Token
func (s *AuthService) GenerateToken(userID uint) (string, error) {
	// 创建 JWT Claims
	claims := jwt.MapClaims{
		"user_id": userID,                             // 用户ID
		"exp":     time.Now().Add(s.jwtExpire).Unix(), // 过期时间
		"iat":     time.Now().Unix(),                  // 签发时间
	}

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并返回
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证 JWT Token
func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	// 解析 Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("Token 解析失败: %w", err)
	}

	// 提取 Claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("Token 中缺少 user_id")
		}
		return uint(userID), nil
	}

	return 0, errors.New("无效的 Token")
}

// GetUserByID 根据 ID 获取用户
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}
