package middleware

import (
	"net/http"
	"strings"

	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/models"
	"github.com/NeoForeverYoung/shadow-oauth/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Header 中获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("缺少认证令牌", nil))
			c.Abort()
			return
		}

		// 2. 检查 Bearer Token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("认证令牌格式错误", nil))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. 验证 Token
		userID, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse("无效的认证令牌", err))
			c.Abort()
			return
		}

		// 4. 将用户ID存入上下文，供后续处理器使用
		c.Set("userID", userID)

		// 5. 继续处理请求
		c.Next()
	}
}

