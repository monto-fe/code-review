package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

var (
	cache     *lru.Cache[string, int64]
	cacheOnce sync.Once
	ttl       = 5 * time.Second
)

func initCache() {
	var err error
	cache, err = lru.New[string, int64](500) // 缓存容量500条
	if err != nil {
		panic(err)
	}
}

// IsDuplicateWebhook 检查是否为重复请求（5秒内重复视为重复）
func IsDuplicateWebhook(body any) bool {
	cacheOnce.Do(initCache)

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("请求体序列化失败: %v\n", err)
		return false
	}
	hash := hashPayload(jsonBytes)

	now := time.Now().Unix()
	if ts, ok := cache.Get(hash); ok {
		if now-ts <= int64(ttl.Seconds()) {
			return true
		}
	}
	cache.Add(hash, now)
	return false
}

func hashPayload(b []byte) string {
	h := sha256.Sum256(b)
	return fmt.Sprintf("%x", h[:])
}
