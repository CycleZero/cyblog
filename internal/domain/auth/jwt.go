package auth

import (
	"context"
	"cyblog/pkg/log"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	ErrInvalidToken = errors.New("无效的 token")
	ErrExpiredToken = errors.New("token 已过期")
)

// Payload JWT 负载结构
type Payload struct {
	jwt.RegisteredClaims
	UserID   uint
	UserName string
}

type JwtBiz struct {
	logger        *log.Logger
	secretKey     string
	tokenExpire   time.Duration
	refreshExpire time.Duration
	redisClient   *redis.Client
	Issuer        string
	BotToken      string
}

func NewJwtBiz(logger *log.Logger, vc *viper.Viper, redisClient *redis.Client) *JwtBiz {
	return &JwtBiz{
		logger:        logger,
		secretKey:     vc.GetString("jwt.secret"),
		tokenExpire:   time.Hour * 24 * 7,
		refreshExpire: time.Hour * 24 * 30 * 3,
		redisClient:   redisClient,
		Issuer:        "cyblog",
		BotToken:      vc.GetString("app.bot_token"),
	}
}

// GenerateToken 生成访问令牌
func (b *JwtBiz) GenerateToken(ctx context.Context, userID uint, userName string) (string, string, error) {
	tokenID := uuid.New().String()
	now := time.Now()
	expiresAt := now.Add(b.tokenExpire)

	claims := &Payload{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    b.Issuer,
			Subject:   userName,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        tokenID,
		},
		UserID:   userID,
		UserName: userName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(b.secretKey))
	if err != nil {
		b.logger.Sugar().Errorw("生成 token 失败", "error", err)
		return "", "", err
	}

	// 生成刷新令牌
	refreshTokenID := uuid.New().String()
	refreshExpiresAt := now.Add(b.refreshExpire)

	refreshClaims := &jwt.RegisteredClaims{
		Issuer:    b.Issuer,
		Subject:   userName,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
		ID:        refreshTokenID,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(b.secretKey))
	if err != nil {
		b.logger.Sugar().Errorw("生成 refresh token 失败", "error", err)
		return "", "", err
	}

	// 将 token ID 存入 Redis，用于黑名单管理
	b.storeTokenInfo(ctx, tokenID, userID, userName, expiresAt)

	b.logger.Sugar().Infow("生成 token 成功", "userID", userID, "userName", userName, "tokenID", tokenID)
	return tokenString, refreshTokenString, nil
}

// DecodeJWT 解析并验证 token
func (b *JwtBiz) DecodeJWT(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(b.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Payload); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// RefreshToken 刷新 token
func (b *JwtBiz) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// 解析刷新令牌
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(b.secretKey), nil
	})

	if err != nil {
		return "", "", ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", "", ErrInvalidToken
	}

	// 检查刷新令牌是否在黑名单中
	if b.IsBlacklist(claims.ID) {
		return "", "", ErrInvalidToken
	}

	// 从刷新令牌中获取用户信息（需要额外存储或查询）
	// 这里简化处理，实际项目中可能需要从 Redis 中获取用户信息
	userID, userName, err := b.getTokenUserInfo(ctx, claims.ID)
	if err != nil {
		return "", "", err
	}

	// 将旧令牌加入黑名单
	b.AddToBlacklist(claims.ID)

	// 生成新的访问令牌和刷新令牌
	return b.GenerateToken(ctx, userID, userName)
}

// AddToBlacklist 将 token 加入黑名单
func (b *JwtBiz) AddToBlacklist(tokenID string) {
	ctx := context.Background()
	key := b.getBlacklistKey(tokenID)
	// 设置过期时间与 token 剩余有效期一致
	_, err := b.redisClient.Set(ctx, key, "1", b.tokenExpire).Result()
	if err != nil {
		b.logger.Sugar().Errorw("将 token 加入黑名单失败", "tokenID", tokenID, "error", err)
	}
}

// IsBlacklist 检查 token 是否在黑名单中
func (b *JwtBiz) IsBlacklist(tokenID string) bool {
	ctx := context.Background()
	key := b.getBlacklistKey(tokenID)
	exists, err := b.redisClient.Exists(ctx, key).Result()
	if err != nil {
		b.logger.Sugar().Errorw("检查 token 黑名单失败", "tokenID", tokenID, "error", err)
		return false
	}
	return exists > 0
}

// AddUserTokenToBlacklist 将用户的所有 token 加入黑名单（用于登出或修改密码）
func (b *JwtBiz) AddUserTokenToBlacklist(ctx context.Context, userID uint, issuedAt time.Time) {
	key := b.getUserTokenBlacklistKey(userID)
	// 存储用户 token 黑名单的 issuedAt 时间戳
	_, err := b.redisClient.Set(ctx, key, issuedAt.UnixNano(), 0).Result()
	if err != nil {
		b.logger.Sugar().Errorw("将用户 token 加入黑名单失败", "userID", userID, "error", err)
	}
}

// IsBlacklistUserTokens 检查用户的所有 token 是否被禁用
func (b *JwtBiz) IsBlacklistUserTokens(userID uint, issuedAt time.Time) bool {
	ctx := context.Background()
	key := b.getUserTokenBlacklistKey(userID)
	val, err := b.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return false
	}
	if err != nil {
		b.logger.Sugar().Errorw("检查用户 token 黑名单失败", "userID", userID, "error", err)
		return false
	}

	// 如果存在黑名单记录，检查 issuedAt 是否早于黑名单时间
	var blacklistTime int64
	_, err = fmt.Sscanf(val, "%d", &blacklistTime)
	if err != nil {
		return false
	}

	return issuedAt.UnixNano() < blacklistTime
}

// ClearUserTokenBlacklist 清除用户的 token 黑名单
func (b *JwtBiz) ClearUserTokenBlacklist(ctx context.Context, userID uint) {
	key := b.getUserTokenBlacklistKey(userID)
	_, err := b.redisClient.Del(ctx, key).Result()
	if err != nil {
		b.logger.Sugar().Errorw("清除用户 token 黑名单失败", "userID", userID, "error", err)
	}
}

// storeTokenInfo 存储 token 信息到 Redis
func (b *JwtBiz) storeTokenInfo(ctx context.Context, tokenID string, userID uint, userName string, expiresAt time.Time) {
	key := b.getTokenInfoKey(tokenID)
	_, err := b.redisClient.HSet(ctx, key, "user_id", userID, "user_name", userName).Result()
	if err != nil {
		b.logger.Sugar().Errorw("存储 token 信息失败", "tokenID", tokenID, "error", err)
		return
	}
	// 设置过期时间
	_, err = b.redisClient.ExpireAt(ctx, key, expiresAt).Result()
	if err != nil {
		b.logger.Sugar().Errorw("设置 token 信息过期时间失败", "tokenID", tokenID, "error", err)
	}
}

// getTokenUserInfo 从 Redis 获取 token 关联的用户信息
func (b *JwtBiz) getTokenUserInfo(ctx context.Context, tokenID string) (uint, string, error) {
	key := b.getTokenInfoKey(tokenID)
	userIDStr, err := b.redisClient.HGet(ctx, key, "user_id").Result()
	if err != nil {
		return 0, "", err
	}

	userName, err := b.redisClient.HGet(ctx, key, "user_name").Result()
	if err != nil {
		return 0, "", err
	}

	var userID uint
	_, err = fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		return 0, "", err
	}

	return userID, userName, nil
}

// getBlacklistKey 获取黑名单 key
func (b *JwtBiz) getBlacklistKey(tokenID string) string {
	return "jwt:blacklist:" + tokenID
}

// getUserTokenBlacklistKey 获取用户 token 黑名单 key
func (b *JwtBiz) getUserTokenBlacklistKey(userID uint) string {
	return fmt.Sprintf("jwt:user_blacklist:%d", userID)
}

// getTokenInfoKey 获取 token 信息 key
func (b *JwtBiz) getTokenInfoKey(tokenID string) string {
	return "jwt:token_info:" + tokenID
}
