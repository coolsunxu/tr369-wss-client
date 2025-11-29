package common

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"unsafe"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func LoadJsonFile(filePath string) map[string]interface{} {
	var result map[string]interface{}
	for i := 0; i < 3; i++ {
		file, err := os.Open(filePath)
		if err != nil {
			slog.Warn("failed to open json file", "path", filePath, "index", i, "error", err)
			time.Sleep(1 * time.Second)
			continue
		}

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&result)
		if err != nil {
			slog.Warn("failed to decode json data into map", "path", filePath, "index", i, "error", err)
			time.Sleep(1 * time.Second)
			continue
		}

		err = file.Close()
		if err != nil {
			slog.Warn("failed to close file", "path", filePath, "index", i, "error", err)
		}
		break
	}

	return result
}

func SaveJsonFile(trtree map[string]interface{}, filePath string) {
	file, e := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			slog.Warn("close file failed", "path", filePath)
		}
	}(file)
	if e != nil {
		slog.Warn("failed to open file", "path", filePath)
	} else {
		slog.Debug("open file succeed", "path", filePath)
	}
	// 创建编码器
	encoder := json.NewEncoder(file)
	err := encoder.Encode(trtree)
	if err != nil {
		slog.Warn("encode failed", "path", filePath)
	} else {
		slog.Debug("encode succeed", "path", filePath)
	}
}

// FormatMacWithColon 54AF971CC687 -> 54:AF:97:1C:C6:87
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

func HmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func ReportQoeData(url string, authorization string, payload string) error {

	body := bytes.NewReader([]byte(payload))
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("BBF-Report-Format", "ObjectHierarchy")
	request.Header.Set("Authorization", authorization)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return errors.New(resp.Status)
	}

	return nil
}

func RandStr(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
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
	return *(*string)(unsafe.Pointer(&b))
}
