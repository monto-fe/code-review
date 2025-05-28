package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"code-review-go/internal/pkg/response"
)

var (
	// WhitePathList 白名单路径
	WhitePathList = []string{
		"/v1/login",
		"/v1/register",
		"/v1/webhook/merge",
	}

	// TokenSecretKey JWT密钥
	TokenSecretKey = []byte("your-secret-key") // TODO: 从配置文件读取
)

// AuthenticateJWT JWT认证中间件
func AuthenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否在白名单中
		path := c.Request.URL.Path
		for _, whitePath := range WhitePathList {
			if path == whitePath {
				c.Next()
				return
			}
		}

		// 获取JWT token
		token := c.GetHeader("jwt_token")
		if token == "" {
			response.Error(c, nil, "Token is required", 10012)
			c.Abort()
			return
		}

		// 验证token
		claims, err := validateToken(token)
		if err != nil {
			response.Error(c, err, "Invalid token", 10014)
			c.Abort()
			return
		}

		// 检查token是否过期
		exp, ok := claims["exp"].(float64)
		if !ok {
			response.Error(c, nil, "Invalid token expiration", 10016)
			c.Abort()
			return
		}

		if time.Now().Unix() > int64(exp) {
			response.Error(c, nil, "Token expired", 10016)
			c.Abort()
			return
		}

		// 如果有效期少于5小时，更新有效期+5h
		// if int64(exp)-time.Now().Unix() < 5*60*60 {
		// 	// TODO: 实现token刷新逻辑
		// }

		// 设置用户信息到请求头
		c.Request.Header.Set("remoteUser", claims["user"].(string))
		c.Request.Header.Set("userId", fmt.Sprintf("%v", claims["id"]))
		c.Request.Header.Set("namespace", claims["namespace"].(string))

		c.Next()
	}
}

// validateToken 验证JWT token
func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return TokenSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
