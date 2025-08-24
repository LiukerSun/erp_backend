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
	TokenType       string `json:"token_type"`       // token类型：access 或 refresh
	jwt.RegisteredClaims
}

// TokenPair 包含访问token和刷新token
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // 访问token过期时间（秒）
}

// GenerateTokenPair 生成访问token和刷新token
func GenerateTokenPair(userID uint, username, role string, passwordVersion uint) (*TokenPair, error) {
	// 生成访问token（短期）
	accessClaims := Claims{
		UserID:          userID,
		Username:        username,
		Role:            role,
		PasswordVersion: passwordVersion,
		TokenType:       "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.AppConfig.JWTExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return nil, err
	}

	// 生成刷新token（长期）
	refreshClaims := Claims{
		UserID:          userID,
		Username:        username,
		Role:            role,
		PasswordVersion: passwordVersion,
		TokenType:       "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.AppConfig.JWTRefreshExpireDays) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(config.AppConfig.JWTExpireHours * 3600), // 转换为秒
	}, nil
}

// GenerateAccessToken 仅生成访问token
func GenerateAccessToken(userID uint, username, role string, passwordVersion uint) (string, error) {
	claims := Claims{
		UserID:          userID,
		Username:        username,
		Role:            role,
		PasswordVersion: passwordVersion,
		TokenType:       "access",
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

// IsTokenExpired 检查token是否即将过期（在指定时间内）
func IsTokenExpired(claims *Claims, thresholdMinutes int) bool {
	exp := claims.ExpiresAt.Time
	threshold := time.Now().Add(time.Duration(thresholdMinutes) * time.Minute)
	return exp.Before(threshold)
}

// IsRefreshToken 检查是否为刷新token
func IsRefreshToken(claims *Claims) bool {
	return claims.TokenType == "refresh"
}

// IsAccessToken 检查是否为访问token
func IsAccessToken(claims *Claims) bool {
	return claims.TokenType == "access"
}
