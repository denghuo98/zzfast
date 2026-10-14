package zencrypt

import "golang.org/x/crypto/bcrypt"

// HashPassword 使用 bcrypt 生成密码哈希（成本为 bcrypt.DefaultCost）
// bcrypt 自带随机盐，无需额外盐值
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// VerifyPassword 校验明文密码与 bcrypt 哈希是否匹配
func VerifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// IsBcryptHash 判断哈希是否为 bcrypt 格式（用于兼容历史 MD5 数据）
func IsBcryptHash(hash string) bool {
	return len(hash) >= 4 && (hash[:4] == "$2a$" || hash[:4] == "$2b$" || hash[:4] == "$2y$")
}
