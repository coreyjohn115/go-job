package blog

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT配置
var (
	JWTSecret           = []byte("jss@13&^()") // 生产环境中应该使用环境变量
	TokenExpireDuration = time.Hour * 24       // Token过期时间：24小时
)

// 自定义Claims结构
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint, username, email string) (string, error) {
	// 创建Claims
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "jss",
			Subject:   "user-token",
		},
	}

	// 创建Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名Token
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证Token并提取Claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新Token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 检查Token是否即将过期（剩余时间少于1小时）
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return "", errors.New("token doesn't need refresh yet")
	}

	// 生成新Token
	return GenerateToken(claims.UserID, claims.Username, claims.Email)
}
