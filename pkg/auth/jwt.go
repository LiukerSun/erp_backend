package auth

import (
	"errors"
	"time"

	"erp/config"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明
type Claims struct {
	UserID          uint   `json:"user_id"`
	Username        string `json:"username"`
	Role            string `json:"role"`
	PasswordVersion uint   `json:"password_version"` // 密码版本，用于验证token是否有效
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username, role string, passwordVersion uint) (string, error) {
	claims := Claims{
		UserID:          userID,
		Username:        username,
		Role:            role,
		PasswordVersion: passwordVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.AppConfig.JWTExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// ValidateTokenPasswordVersion 验证token的密码版本是否有效
// 这个函数需要从数据库获取当前用户的密码版本进行比较
func ValidateTokenPasswordVersion(claims *Claims, currentPasswordVersion uint) bool {
	return claims.PasswordVersion == currentPasswordVersion
}
