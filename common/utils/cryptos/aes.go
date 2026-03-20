package cryptos

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// =================== GCM（推荐）======================
// AES-GCM 同时提供保密性 + 完整性认证，无 Padding Oracle 风险。
// 密文格式：nonce(12B) + ciphertext + tag(16B)

// AesEncrypt 使用 AES-GCM 加密（推荐）
// key 长度必须为 16、24 或 32 字节（对应 AES-128/192/256）
func AesEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	// Seal 将 nonce 作为前缀附加到密文
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// AesDecrypt 使用 AES-GCM 解密，同时校验完整性（推荐）
func AesDecrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize+gcm.Overhead() {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// 统一返回，不暴露是 nonce 还是 tag 校验失败
		return nil, errors.New("decryption failed")
	}
	return plaintext, nil
}

// =================== CBC（兼容保留）======================
// 注意：CBC 不提供消息认证，存在 Padding Oracle 攻击风险。
// 新业务请使用 AesEncrypt / AesDecrypt（GCM）。

// AesEncryptCBC 使用随机 IV 进行 AES-CBC 加密
// 密文格式：IV(16B) + ciphertext
// Deprecated: 请使用 AesEncrypt（GCM）替代
func AesEncryptCBC(origData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pkcs7Padding(origData, blockSize)

	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(origData, origData)

	encrypted := make([]byte, blockSize+len(origData))
	copy(encrypted[:blockSize], iv)
	copy(encrypted[blockSize:], origData)
	return encrypted, nil
}

// AesDecryptCBC 解密 AES-CBC（从密文中提取 IV）
// Deprecated: 请使用 AesDecrypt（GCM）替代
func AesDecryptCBC(encrypted []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(encrypted) < blockSize*2 || len(encrypted)%blockSize != 0 {
		return nil, errors.New("invalid ciphertext")
	}
	iv, ciphertext := encrypted[:blockSize], encrypted[blockSize:]
	decrypted := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(decrypted, ciphertext)
	return pkcs7UnPadding(decrypted, blockSize)
}

// AesEncryptCBCWithIV 使用指定 IV 进行 AES-CBC 加密
// Deprecated: 固定 IV 不安全，请使用 AesEncrypt（GCM）替代
func AesEncryptCBCWithIV(origData []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		return nil, errors.New("IV length must equal block size")
	}
	origData = pkcs7Padding(origData, blockSize)
	encrypted := make([]byte, len(origData))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(encrypted, origData)
	return encrypted, nil
}

// AesDecryptCBCWithIV 使用指定 IV 进行 AES-CBC 解密
// Deprecated: 固定 IV 不安全，请使用 AesDecrypt（GCM）替代
func AesDecryptCBCWithIV(encrypted []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		return nil, errors.New("IV length must equal block size")
	}
	if len(encrypted)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}
	decrypted := make([]byte, len(encrypted))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(decrypted, encrypted)
	return pkcs7UnPadding(decrypted, blockSize)
}

// =================== 内部辅助 ======================

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padByte := byte(padding)
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = padByte
	}
	return padded
}

// pkcs7UnPadding 使用恒定时间比较，防止 Padding Oracle timing 攻击
func pkcs7UnPadding(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("empty data")
	}
	padding := int(data[length-1])
	// padding 合法范围：[1, blockSize]
	if padding == 0 || padding > blockSize || padding > length {
		return nil, errors.New("invalid padding")
	}
	// 恒定时间遍历所有 padding 字节，消除 timing oracle
	var invalid byte
	for i := length - padding; i < length; i++ {
		invalid |= data[i] ^ byte(padding)
	}
	if invalid != 0 {
		return nil, errors.New("invalid padding")
	}
	return data[:length-padding], nil
}
