package cryptos

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// =================== CBC ======================

// AesEncryptCBC 使用随机 IV 进行 AES-CBC 加密
// 返回的密文格式：IV + ciphertext
// key的长度必须为16, 24或者32
func AesEncryptCBC(origData []byte, key []byte) (encrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pkcs5Padding(origData, blockSize)

	// 生成随机 IV
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(origData))
	blockMode.CryptBlocks(ciphertext, origData)

	// 将 IV 放在密文前面
	encrypted = make([]byte, blockSize+len(ciphertext))
	copy(encrypted[:blockSize], iv)
	copy(encrypted[blockSize:], ciphertext)

	return encrypted, nil
}

// AesDecryptCBC 解密 AES-CBC（从密文中提取 IV）
// 密文格式：IV + ciphertext
// key的长度必须为16, 24或者32
func AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()

	if len(encrypted) < blockSize*2 {
		return nil, errors.New("ciphertext too short")
	}
	if len(encrypted)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	// 提取 IV 和实际密文
	iv := encrypted[:blockSize]
	ciphertext := encrypted[blockSize:]

	blockMode := cipher.NewCBCDecrypter(block, iv)
	decrypted = make([]byte, len(ciphertext))
	blockMode.CryptBlocks(decrypted, ciphertext)

	decrypted, err = pkcs5UnPadding(decrypted)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

// AesEncryptCBCWithIV 使用指定 IV 进行 AES-CBC 加密（兼容旧接口）
// 注意：固定 IV 不安全，建议使用 AesEncryptCBC
func AesEncryptCBCWithIV(origData []byte, key []byte, iv []byte) (encrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		return nil, errors.New("IV length must equal block size")
	}
	origData = pkcs5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted = make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

// AesDecryptCBCWithIV 使用指定 IV 进行 AES-CBC 解密（兼容旧接口）
func AesDecryptCBCWithIV(encrypted []byte, key []byte, iv []byte) (decrypted []byte, err error) {
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
	blockMode := cipher.NewCBCDecrypter(block, iv)
	decrypted = make([]byte, len(encrypted))
	blockMode.CryptBlocks(decrypted, encrypted)
	decrypted, err = pkcs5UnPadding(decrypted)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("empty data")
	}
	unpadding := int(origData[length-1])
	if unpadding > length || unpadding == 0 {
		return nil, errors.New("invalid padding")
	}
	// 验证所有 padding 字节是否一致
	for i := length - unpadding; i < length; i++ {
		if origData[i] != byte(unpadding) {
			return nil, errors.New("invalid padding")
		}
	}
	return origData[:(length - unpadding)], nil
}
