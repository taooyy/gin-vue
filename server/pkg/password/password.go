// server/pkg/password/password.go
package password

import "golang.org/x/crypto/bcrypt"

// Hash 使用 bcrypt 加密密码
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Check 检查密码与哈希值是否匹配
func Check(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
