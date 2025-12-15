package handlers

import (
	"net/http"

	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/models"
	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 处理用户注册请求
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册信息"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest

	// 1. 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("请求参数无效", err))
		return
	}

	// 2. 调用服务层处理注册逻辑
	user, err := h.authService.Register(req)
	if err != nil {
		// 根据不同错误类型返回不同的 HTTP 状态码
		switch err {
		case service.ErrInvalidEmail, service.ErrWeakPassword:
			c.JSON(http.StatusBadRequest, models.ErrorResponse("注册失败", err))
		case service.ErrEmailExists:
			c.JSON(http.StatusConflict, models.ErrorResponse("注册失败", err))
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse("注册失败", err))
		}
		return
	}

	// 3. 返回成功响应
	c.JSON(http.StatusCreated, models.SuccessResponse("注册成功", user.ToResponse()))
}

// Login 处理用户登录请求
// @Summary 用户登录
// @Description 验证用户凭证并返回 JWT Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest

	// 1. 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("请求参数无效", err))
		return
	}

	// 2. 调用服务层处理登录逻辑
	loginResp, err := h.authService.Login(req)
	if err != nil {
		// 根据不同错误类型返回不同的 HTTP 状态码
		switch err {
		case service.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("登录失败", err))
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse("登录失败", err))
		}
		return
	}

	// 3. 返回成功响应（包含 Token 和用户信息）
	c.JSON(http.StatusOK, models.SuccessResponse("登录成功", loginResp))
}

// GetCurrentUser 获取当前登录用户信息
// @Summary 获取当前用户
// @Description 根据 JWT Token 获取当前登录用户的信息
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /api/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// 从上下文中获取用户ID（由 JWT 中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("未授权", nil))
		return
	}

	// 查询用户信息
	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("获取用户信息失败", err))
		return
	}

	// 返回用户信息
	c.JSON(http.StatusOK, models.SuccessResponse("获取成功", user.ToResponse()))
}
