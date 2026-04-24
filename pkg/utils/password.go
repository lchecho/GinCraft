package utils

import "golang.org/x/crypto/bcrypt"

const bcryptCost = bcrypt.DefaultCost

// HashPassword 使用 bcrypt 对密码做散列
func HashPassword(raw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(raw), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPassword 验证密码与 bcrypt 哈希是否匹配
func CheckPassword(raw, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw)) == nil
}
