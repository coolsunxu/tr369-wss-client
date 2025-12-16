// Package string 提供字符串处理工具
package string

import (
	"math/rand"
	"strings"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdBits = 6
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

// RandStr 生成指定长度的随机字符串
func RandStr(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return string(b)
}

// FormatMacWithColon 格式化 MAC 地址，添加冒号分隔
// 例如: 54AF971CC687 -> 54:AF:97:1C:C6:87
func FormatMacWithColon(mac string) string {
	var stringBuilder strings.Builder
	index := 2
	for index <= len(mac) {
		stringBuilder.WriteString(mac[index-2 : index])
		if index < len(mac) {
			stringBuilder.WriteString(":")
		}
		index += 2
	}
	return stringBuilder.String()
}
