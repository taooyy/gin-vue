package jwt

import (
	"server/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 自定义 JWT Claims
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 Token
func GenerateToken(userID uint, username, role string) (string, error) {
	// 创建一个我们自己的声明
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Cfg.Jwt.Expire) * time.Hour)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                       // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                                       // 生效时间
			Issuer:    config.Cfg.Jwt.Issuer,                                                                // 签发者
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用指定的 secret 签名并获得完整的编码后的字符串 token
	return token.SignedString([]byte(config.Cfg.Jwt.Secret))
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(config.Cfg.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}
