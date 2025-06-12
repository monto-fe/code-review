package middleware

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

var testSecret = []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlciI6ImFkbWluIiwibmFtZXNwYWNlIjoiYWNsIiwiY3JlYXRlX3RpbWUiOjE3NDk3MTU2MjksImlhdCI6MTc0OTcxNTYyOSwiZXhwIjoxNzQ5NzQ0NDI5fQ.OjeUzlCnwmqn8xJZECLLVBePzsTA2i-_88sn0o9rLcI") // 替换为你的 TokenSecretKey

func TestValidateToken_Success(t *testing.T) {
	// 生成一个有效的 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   1,
		"user": "admin",
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	})
	tokenString, err := token.SignedString(testSecret)
	assert.NoError(t, err)

	// 设置 TokenSecretKey
	TokenSecretKey = testSecret

	claims, err := ValidateToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), claims["id"])
	assert.Equal(t, "admin", claims["user"])
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	// 用错误的密钥签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   1,
		"user": "admin",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte("wrong_secret"))
	assert.NoError(t, err)

	TokenSecretKey = testSecret

	_, err = ValidateToken(tokenString)
	assert.Error(t, err)
}

func TestValidateToken_Expired(t *testing.T) {
	// 生成一个已过期的 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   1,
		"user": "admin",
		"exp":  time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(testSecret)
	assert.NoError(t, err)

	TokenSecretKey = testSecret

	_, err = ValidateToken(tokenString)
	assert.Error(t, err)
}
