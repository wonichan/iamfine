package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"hupu/shared/config"
	"time"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GlobalConfig.JWT.Expire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 验证签发时间、过期时间、生效时间
		now := time.Now()
		
		// 检查token是否已过期
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(now) {
			return nil, errors.New("token已过期")
		}
		
		// 检查token是否还未生效
		if claims.NotBefore != nil && claims.NotBefore.Time.After(now) {
			return nil, errors.New("token还未生效")
		}
		
		// 检查签发时间是否合理（不能是未来时间）
		if claims.IssuedAt != nil && claims.IssuedAt.Time.After(now) {
			return nil, errors.New("token签发时间无效")
		}
		
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
