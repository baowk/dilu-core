package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/baowk/dilu-core/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 自定义 JWT Claims
type CustomClaims struct {
	UserId   int            `json:"uid,omitempty"`
	RoleId   int            `json:"rid,omitempty"`
	Phone    string         `json:"mob,omitempty"`
	Nickname string         `json:"nick,omitempty"`
	Data     map[string]any `json:"data,omitempty"`
	jwt.RegisteredClaims
}

// JWTConfig JWT 中间件配置
type JWTConfig struct {
	SignKey    string // 签名密钥
	Expires   int    // 过期时间（分钟）
	Refresh   int    // 刷新窗口（分钟）
	Issuer    string // 签发人
	Subject   string // 签发主体
	HeaderKey string // 请求头 key（默认 Authorization）
	QueryKey  string // URL 参数 key（可选）
	CookieKey string // Cookie key（可选）
}

// NewJWTConfigFromCfg 从 dilu-core 配置创建 JWT 配置
func NewJWTConfigFromCfg(cfg config.JWT) JWTConfig {
	return JWTConfig{
		SignKey:    cfg.SignKey,
		Expires:   cfg.Expires,
		Refresh:   cfg.Refresh,
		Issuer:    cfg.Issuer,
		Subject:   cfg.Subject,
		HeaderKey: "Authorization",
	}
}

// JWTMiddleware 返回 Gin JWT 认证中间件
func JWTMiddleware(cfg JWTConfig) gin.HandlerFunc {
	if cfg.HeaderKey == "" {
		cfg.HeaderKey = "Authorization"
	}
	return func(c *gin.Context) {
		tokenStr := extractToken(c, cfg)
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未提供认证Token"})
			c.Abort()
			return
		}

		claims := &CustomClaims{}
		if err := parseToken(tokenStr, claims, cfg.SignKey, cfg.Subject); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token无效或已过期"})
			c.Abort()
			return
		}

		// 自动刷新：距离过期不足 Refresh 分钟时刷新
		if cfg.Refresh > 0 {
			if exp, err := claims.GetExpirationTime(); err == nil && exp != nil {
				remaining := time.Until(exp.Time)
				if remaining < time.Duration(cfg.Refresh)*time.Minute {
					newExp := time.Now().Add(time.Duration(cfg.Expires) * time.Minute)
					claims.ExpiresAt = jwt.NewNumericDate(newExp)
					if newToken, err := GenerateToken(claims, cfg.SignKey); err == nil {
						c.Header("X-Refresh-Token", newToken)
						c.Header("X-Token-Expires", strconv.FormatInt(newExp.Unix(), 10))
					}
				}
			}
		}

		// 将用户信息注入上下文 + 请求头（供下游/网关使用）
		c.Set("jwt_claims", claims)
		c.Set("jwt_uid", claims.UserId)
		c.Set("jwt_rid", claims.RoleId)
		c.Request.Header.Set("a_uid", fmt.Sprintf("%d", claims.UserId))
		c.Request.Header.Set("a_rid", fmt.Sprintf("%d", claims.RoleId))
		c.Request.Header.Set("a_mobile", claims.Phone)
		c.Request.Header.Set("a_nickname", claims.Nickname)

		c.Next()
	}
}

// GenerateToken 生成 JWT Token
func GenerateToken(claims jwt.Claims, signKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signKey))
}

// NewClaims 创建标准 Claims
func NewClaims(userId, roleId int, phone, nickname string, cfg JWTConfig) *CustomClaims {
	now := time.Now()
	return &CustomClaims{
		UserId:   userId,
		RoleId:   roleId,
		Phone:    phone,
		Nickname: nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.Expires) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    cfg.Issuer,
			Subject:   cfg.Subject,
		},
	}
}

// extractToken 从请求中提取 Token
func extractToken(c *gin.Context, cfg JWTConfig) string {
	// 1. 从 Header 提取
	if cfg.HeaderKey != "" {
		authHeader := c.GetHeader(cfg.HeaderKey)
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				return parts[1]
			}
			// 兼容无 Bearer 前缀
			if len(parts) == 1 {
				return parts[0]
			}
		}
	}

	// 2. 从 Query 参数提取
	if cfg.QueryKey != "" {
		if token := c.Query(cfg.QueryKey); token != "" {
			return token
		}
	}

	// 3. 从 Cookie 提取
	if cfg.CookieKey != "" {
		if cookie, err := c.Cookie(cfg.CookieKey); err == nil && cookie != "" {
			return cookie
		}
	}

	return ""
}

func parseToken(tokenString string, claims *CustomClaims, secret, subject string) error {
	opts := []jwt.ParserOption{}
	if subject != "" {
		opts = append(opts, jwt.WithSubject(subject))
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, opts...)
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}
