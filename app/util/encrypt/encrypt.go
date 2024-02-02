package encrypt

import "golang.org/x/crypto/bcrypt"

// 生成哈希加密字符串
func HashString(pass string) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hashPass), nil
}

// CheckPassword 校验密码是否匹配
func CheckPassword(check, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(check))
	return err
}
