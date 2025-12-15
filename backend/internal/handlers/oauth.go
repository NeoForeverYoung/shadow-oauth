package handlers

import (
	"net/http"
	"net/url"

	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/models"
	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// OAuthHandler OAuth 处理器
type OAuthHandler struct {
	oauthService *service.OAuthService
	authService  *service.AuthService
}

// NewOAuthHandler 创建 OAuth 处理器实例
func NewOAuthHandler(oauthService *service.OAuthService, authService *service.AuthService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		authService:  authService,
	}
}

// AuthorizeRequest 授权请求参数
type AuthorizeRequest struct {
	ClientID     string `form:"client_id" binding:"required"`     // 客户端ID
	RedirectURI  string `form:"redirect_uri" binding:"required"`  // 重定向URI
	ResponseType string `form:"response_type" binding:"required"` // 响应类型（固定为 "code"）
	State        string `form:"state"`                            // 状态参数（用于防止CSRF攻击）
}

// Authorize 授权端点
// GET /oauth/authorize?client_id=xxx&redirect_uri=xxx&response_type=code&state=xxx
// 这是 OAuth 2.0 的第一步：第三方应用引导用户到这里进行授权
func (h *OAuthHandler) Authorize(c *gin.Context) {
	var req AuthorizeRequest

	// 1. 解析查询参数
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("请求参数无效", err))
		return
	}

	// 2. 验证响应类型（简化版只支持授权码模式）
	if req.ResponseType != "code" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("不支持的响应类型", nil))
		return
	}

	// 3. 验证客户端
	client, err := h.oauthService.ValidateClientID(req.ClientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("无效的客户端", err))
		return
	}

	// 4. 验证重定向URI
	if err := h.oauthService.ValidateRedirectURI(client, req.RedirectURI); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("重定向URI不匹配", err))
		return
	}

	// 5. 检查用户是否已登录（通过 JWT Token）
	// 如果未登录，需要先登录
	userID, exists := c.Get("userID")
	if !exists {
		// 未登录，重定向到登录页（登录后返回这里）
		// 简化版：直接返回需要登录的提示
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("请先登录", nil))
		return
	}

	// 6. 生成授权码
	code, err := h.oauthService.GenerateAuthorizationCode(
		req.ClientID,
		userID.(uint),
		req.RedirectURI,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("生成授权码失败", err))
		return
	}

	// 7. 重定向到客户端，带上授权码
	// 格式：redirect_uri?code=xxx&state=xxx
	redirectURL, _ := url.Parse(req.RedirectURI)
	query := redirectURL.Query()
	query.Set("code", code)
	if req.State != "" {
		query.Set("state", req.State) // 原样返回 state 参数
	}
	redirectURL.RawQuery = query.Encode()

	c.Redirect(http.StatusFound, redirectURL.String())
}

// TokenRequest Token 请求参数
type TokenRequest struct {
	GrantType    string `form:"grant_type" binding:"required"`    // 授权类型（固定为 "authorization_code"）
	Code         string `form:"code" binding:"required"`          // 授权码
	RedirectURI  string `form:"redirect_uri" binding:"required"`  // 重定向URI（必须与授权时一致）
	ClientID     string `form:"client_id" binding:"required"`     // 客户端ID
	ClientSecret string `form:"client_secret" binding:"required"` // 客户端密钥
}

// TokenResponse Token 响应
type TokenResponse struct {
	AccessToken string `json:"access_token"` // Access Token
	TokenType   string `json:"token_type"`   // Token 类型（固定为 "Bearer"）
	ExpiresIn   int64  `json:"expires_in"`   // 过期时间（秒）
}

// Token Token 端点
// POST /oauth/token
// 这是 OAuth 2.0 的第二步：第三方应用用授权码交换 Access Token
func (h *OAuthHandler) Token(c *gin.Context) {
	var req TokenRequest

	// 1. 解析请求参数（支持 form-urlencoded 格式）
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("请求参数无效", err))
		return
	}

	// 2. 验证授权类型（简化版只支持授权码模式）
	if req.GrantType != "authorization_code" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("不支持的授权类型", nil))
		return
	}

	// 3. 用授权码交换 Access Token
	accessToken, err := h.oauthService.ExchangeAuthorizationCode(
		req.Code,
		req.ClientID,
		req.ClientSecret,
		req.RedirectURI,
	)
	if err != nil {
		switch err {
		case service.ErrInvalidClient:
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("无效的客户端", err))
		case service.ErrInvalidRedirectURI:
			c.JSON(http.StatusBadRequest, models.ErrorResponse("重定向URI不匹配", err))
		case service.ErrInvalidAuthorizationCode, service.ErrAuthorizationCodeUsed:
			c.JSON(http.StatusBadRequest, models.ErrorResponse("无效的授权码", err))
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse("交换Token失败", err))
		}
		return
	}

	// 4. 返回 Access Token（按照 OAuth 2.0 标准格式）
	expiresIn := int64(accessToken.ExpiresAt.Sub(accessToken.CreatedAt).Seconds())
	c.JSON(http.StatusOK, TokenResponse{
		AccessToken: accessToken.Token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	})
}

// UserInfoRequest 用户信息请求
type UserInfoRequest struct {
	AccessToken string `form:"access_token" binding:"required"` // Access Token
}

// UserInfo 用户信息端点
// GET /oauth/userinfo?access_token=xxx
// 第三方应用使用 Access Token 获取用户信息
func (h *OAuthHandler) UserInfo(c *gin.Context) {
	// 1. 从查询参数或 Header 获取 Token
	token := c.Query("access_token")
	if token == "" {
		// 尝试从 Authorization Header 获取
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}
	}

	if token == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("缺少 access_token", nil))
		return
	}

	// 2. 使用 Token 获取用户信息
	user, err := h.oauthService.GetUserInfo(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("无效的 Token", err))
		return
	}

	// 3. 返回用户信息（不包含敏感信息）
	c.JSON(http.StatusOK, models.SuccessResponse("获取成功", user.ToResponse()))
}
