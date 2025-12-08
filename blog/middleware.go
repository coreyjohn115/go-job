package blog

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization header获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenString := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:] // 去掉"Bearer "前缀
		} else {
			tokenString = authHeader // 兼容直接传token的情况
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			c.Abort()
			return
		}

		// 解析Token
		claims, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中，供后续处理函数使用
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选的JWT认证中间件（不强制要求token）
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		} else {
			tokenString = authHeader
		}

		if tokenString == "" {
			c.Next()
			return
		}

		claims, err := ParseToken(tokenString)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("email", claims.Email)
			c.Set("claims", claims)
			c.Set("authenticated", true)
		}
	}
}

// GetCurrentUser 从上下文中获取当前用户信息
func GetCurrentUser(c *gin.Context) (*CustomClaims, bool) {
	if claims, exists := c.Get("claims"); exists {
		if userClaims, ok := claims.(*CustomClaims); ok {
			return userClaims, true
		}
	}
	return nil, false
}

// GetCurrentUserID 从上下文中获取当前用户ID
func GetCurrentUserID(c *gin.Context) (int, bool) {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(int); ok {
			return id, true
		}
	}
	return 0, false
}
