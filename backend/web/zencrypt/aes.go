package zencrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var errCiphertextTooShort = errors.New("ciphertext too short")

// AesGCMEncrypt 使用 AES-GCM 加密，返回 nonce+密文（密文尾部携带 GCM 认证标签）
// key 长度必须为 16/24/32 字节
func AesGCMEncrypt(src, key []byte) (dst []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, src, nil), nil
}

// AesGCMDecrypt 使用 AES-GCM 解密
func AesGCMDecrypt(ciphertext, key []byte) (dst []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errCiphertextTooShort
	}
	nonce, ct := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ct, nil)
}
