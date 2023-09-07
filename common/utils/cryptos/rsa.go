package cryptos

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

// // 将RSA私钥转换为字符串
// func PrivateKeyToPemStr(privateKey *rsa.PrivateKey) string {
// 	return string(PrivateKeyToPem(privateKey))
// }

// 将RSA私钥转换为字符串
func PrivateKeyToPem(privateKey *rsa.PrivateKey) []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return pem.EncodeToMemory(privateKeyBlock)
}

// 将RSA公钥转换为字符串
func PublicKeyToPemStr(publicKey *rsa.PublicKey) (string, error) {
	b, err := PublicKeyToPem(publicKey)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// 将RSA公钥转换为字符串
func PublicKeyToPem(publicKey *rsa.PublicKey) ([]byte, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	return pem.EncodeToMemory(publicKeyBlock), nil
}

func ParsePriKey(privateKey []byte) (*rsa.PrivateKey, error) {
	privateKeyBlock, _ := pem.Decode(privateKey)
	if privateKeyBlock == nil || privateKeyBlock.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("无效的私钥")
	}
	priKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return priKey, nil
}

// 将字符串转换为RSA公钥
func ParsePubKey(publicKey []byte) (*rsa.PublicKey, error) {
	publicKeyBlock, _ := pem.Decode(publicKey)
	if publicKeyBlock == nil || publicKeyBlock.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("无效的公钥")
	}
	pubKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublicKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("无效的RSA公钥")
	}
	return rsaPublicKey, nil
}

// 生成密钥对
func GenerateRsaKeyStr(len int) (publicKey string, privateKey string, err error) {
	pub, pri, err := GenerateRsaKey(len)
	if err != nil {
		return
	}
	publicKey = string(pub)
	privateKey = string(pri)
	return
}

// 生成密钥对
func GenerateRsaKey(len int) (publicKey []byte, privateKey []byte, err error) {
	if len != 1024 && len != 4096 {
		len = 2048
	}
	// 生成RSA密钥对
	key, err := rsa.GenerateKey(rand.Reader, len)
	if err != nil {
		return
	}

	publicKey, err = PublicKeyToPem(&key.PublicKey)
	if err != nil {
		return
	}
	privateKey = PrivateKeyToPem(key)
	return
}

// 公钥加密
func RSA_Encrypt(message []byte, pubKey string) (string, error) {
	return RsaEncrypt(message, []byte(pubKey))
}

// 公钥加密
func RsaEncrypt(message []byte, pubKey []byte) (string, error) {
	publicKey, err := ParsePubKey(pubKey)
	if err != nil {
		return "", err
	}
	// 使用公钥加密消息
	encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(encryptedMessage), nil
}

// 私钥解密
func RsaDecrypt(encryptedMsg string, priKey []byte) ([]byte, error) {
	privateKey, err := ParsePriKey(priKey)
	if err != nil {
		return nil, err
	}
	// 使用私钥解密消息
	cipherText, err := hex.DecodeString(encryptedMsg)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
}

// 私钥解密
func RSA_Decrypt(encryptedMsg, priKey string) ([]byte, error) {
	return RsaDecrypt(encryptedMsg, []byte(priKey))
}

// RsaSign 私钥加签
func RSA_Sign(publicKeyPEM string, message []byte) ([]byte, error) {
	return RsaSign([]byte(publicKeyPEM), message)
}

// RsaSign 私钥加签
func RsaSign(publicKeyPEM, message []byte) ([]byte, error) {
	privateKey, _ := ParsePriKey(publicKeyPEM)
	hash := sha256.Sum256(message)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
}

// RsaVerify 公钥验签
func RSA_Verify(publicKeyPEM string, message []byte, signature []byte) error {
	return RsaVerify([]byte(publicKeyPEM), message, signature)
}

// RsaVerify 公钥验签
func RsaVerify(publicKeyPEM, message []byte, signature []byte) error {
	publicKey, _ := ParsePubKey(publicKeyPEM)
	hash := sha256.Sum256(message)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
}
