package util

import (
	"cyblog/conf"
	"cyblog/pkg/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId uint   `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// 生成JWT Token
func GenerateToken(user *model.User) (string, error) {
	jwtSecret := []byte(conf.GetConfig().GetString("app.jwt_secret"))
	expireHours := conf.GetConfig().GetInt("app.jwt_expire_hours")
	if expireHours == 0 {
		expireHours = 24 * 7 // 默认7天过期
	}

	claims := Claims{
		UserId: user.ID,
		Name:   user.Name,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cyblog",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	jwtSecret := []byte(conf.GetConfig().GetString("app.jwt_secret"))

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}