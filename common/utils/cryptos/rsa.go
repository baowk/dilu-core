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

// 将RSA私钥转换为字符串
func PrivateKeyToString(privateKey *rsa.PrivateKey) (string, error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyStr := string(pem.EncodeToMemory(privateKeyBlock))
	return privateKeyStr, nil
}

// 将RSA公钥转换为字符串
func PublicKeyToString(publicKey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyStr := string(pem.EncodeToMemory(publicKeyBlock))
	return publicKeyStr, nil
}

// 将字符串转换为RSA私钥
func ParsePrivateKey(privateKeyStr string) (*rsa.PrivateKey, error) {
	return ParsePriKey([]byte(privateKeyStr))
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

func ParsePublicKey(publicKeyStr string) (*rsa.PublicKey, error) {
	return ParsePubKey([]byte(publicKeyStr))
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
func GenerateRsaKey(len int) (string, string, error) {
	if len != 1024 && len != 4096 {
		len = 2048
	}
	// 生成RSA密钥对
	key, err := rsa.GenerateKey(rand.Reader, len)
	if err != nil {
		return "", "", err
	}

	publicKeyStr, err := PublicKeyToString(&key.PublicKey)
	if err != nil {
		return "", "", err
	}
	privateKeyStr, err := PrivateKeyToString(key)
	if err != nil {
		return "", "", err
	}
	return publicKeyStr, privateKeyStr, nil
}

// 公钥加密
func RSA_Encrypt(message []byte, pubKey string) (string, error) {
	publicKey, err := ParsePublicKey(pubKey)
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
func RSA_Decrypt(encryptedMsg, priKey string) ([]byte, error) {
	privateKey, err := ParsePrivateKey(priKey)
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

// RsaSign 私钥加签
func RsaSign(publicKeyPEM string, message []byte) ([]byte, error) {
	privateKey, _ := ParsePrivateKey(publicKeyPEM)
	hash := sha256.Sum256(message)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
}

// RsaVerify 公钥验签
func RsaVerify(publicKeyPEM string, message []byte, signature []byte) error {
	publicKey, _ := ParsePublicKey(publicKeyPEM)
	hash := sha256.Sum256(message)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
}
