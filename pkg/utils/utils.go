package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

// MD5 计算字符串的MD5值
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// RandomString 生成指定长度的随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ParseTime 解析时间字符串
func ParseTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", str)
}

// GetCurrentTimestamp 获取当前时间戳
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// GetCurrentTime 获取当前时间字符串
func GetCurrentTime() string {
	return FormatTime(time.Now())
}