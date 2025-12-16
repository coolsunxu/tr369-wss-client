// Package crypto 提供加密相关工具
package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HmacSha256 使用 HMAC-SHA256 算法计算签名
func HmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
