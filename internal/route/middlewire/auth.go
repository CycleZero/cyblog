package middlewire

import (
	"crypto/rand"
	"cyblog/conf"
	"cyblog/internal/common"
	"cyblog/internal/domain/auth"
	"cyblog/pkg/log"
	"cyblog/pkg/repo"
	"encoding/json"
	"errors"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ... existing code ...

var ExcludePath = []string{
	"/api/auth/register",
	"/api/auth/login",
}

func JwtAuthMiddleWire(jwtBiz *auth.JwtBiz, userRepo *repo.UserRepo) func(bool) gin.HandlerFunc {

	return func(required bool) gin.HandlerFunc {
		return func(c *gin.Context) {

			if InArray(c.Request.URL.Path, ExcludePath) {
				c.Next()
				return
			}
			token := c.GetHeader("Authorization")
			token, err := ExtractBearerToken(token)
			if err != nil {
				if !required {
					log.SugaredLogger().Warnf("认证失败：Token 验证失败 - %v", err)
					c.Next()
					return
				}
				log.SugaredLogger().Warnf("认证失败：Token 格式错误 - %v", err)
				c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
				return
			}

			//if token == jwtBiz.BotToken {
			//	u, err := userRepo.GetUserById(10)
			//	if err != nil {
			//		if optional {
			//			log.SugaredLogger().Warnf("认证失败：Token 验证失败 - %v", err)
			//			c.Next()
			//			return
			//		}
			//		log.SugaredLogger().Errorf("认证失败：查询用户信息失败 - UID: %d, 错误：%v", 10, err)
			//		c.AbortWithStatusJSON(500, gin.H{"message": "Invalid token"})
			//		return
			//	}
			//	meta := common.GetRequestMetadata(c)
			//	meta.UserID = u.ID
			//	meta.User = *u
			//	if conf.IsDevMode() {
			//		log.SugaredLogger().Debug(meta)
			//		s, _ := json.Marshal(meta)
			//		log.SugaredLogger().Debug("meta:", string(s))
			//	}
			//
			//	log.SugaredLogger().Debugf("认证成功 - UID: %d, Username: %s", u.ID, u.Name)
			//
			//	common.SetRequestMetadata(c, meta)
			//	c.Next()
			//	return
			//}

			payload, err := jwtBiz.DecodeJWT(token)
			if err != nil {
				if !required {
					log.SugaredLogger().Warnf("认证失败：Token 验证失败 - %v", err)
					c.Next()

					return
				}
				log.SugaredLogger().Warnf("认证失败：Token 解码失败 - %v", err)
				c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
				return
			}

			if jwtBiz.IsBlacklist(payload.ID) {
				if !required {
					log.SugaredLogger().Warnf("认证失败：Token 验证失败 - %v", err)

					c.Next()
					return
				}
				log.SugaredLogger().Warnf("认证失败：Token 已在黑名单中 - TokenID: %s, UID: %d", payload.ID, payload.UserID)
				c.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
				return
			}

			u, err := userRepo.GetUserById(c, payload.UserID)
			if err != nil {
				if !required {
					log.SugaredLogger().Warnf("认证失败：Token 验证失败 - %v", err)
					c.Next()
					return
				}
				log.SugaredLogger().Errorf("认证失败：查询用户信息失败 - UID: %d, 错误：%v", payload.UserID, err)
				c.AbortWithStatusJSON(500, gin.H{"message": "Invalid token"})
				return
			}

			meta := common.GetRequestMetadata(c)
			meta.UserID = payload.UserID
			meta.User = *u
			if conf.IsDevMode() {
				log.SugaredLogger().Debug(meta)
				s, _ := json.Marshal(meta)
				log.SugaredLogger().Debug("meta:", string(s))
			}

			log.SugaredLogger().Debugf("认证成功 - UID: %d, Username: %s", payload.UserID, u.Name)

			common.SetRequestMetadata(c, meta)
			c.Next()

		}
	}

}

// ... existing code ...

func ExtractBearerToken(token string) (string, error) {
	if token == "" {
		return "", errors.New("token is empty")
	}
	if !strings.HasPrefix(token, "Bearer ") {
		return token, nil
	}
	return strings.TrimPrefix(token, "Bearer "), nil
}
func GenerateRandomStringSecure(length int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[n.Int64()]
	}
	return string(result), nil
}

func GenerateRequestId() string {
	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	randomSuffix, err := GenerateRandomStringSecure(8)
	if err != nil {
		// fallback to timestamp only if crypto fails
		return now
	}
	return now + randomSuffix
}

// InArray 检查值是否在数组中
func InArray(str string, arr []string) bool {
	for _, v := range arr {
		if str == v {
			return true
		}
	}
	return false
}
