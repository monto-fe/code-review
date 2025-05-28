package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// ComparePasswords 比较密码
func ComparePasswords(hashedPassword, plainPassword string) bool {
	hashedInput := hex.EncodeToString(md5.New().Sum([]byte(plainPassword)))
	return hashedPassword == hashedInput
}

// HashPassword 加密密码
func HashPassword(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
