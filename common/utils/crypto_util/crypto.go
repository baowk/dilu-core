package crypto_util

import (
	"crypto/rand"
	"crypto/rsa"
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
	privateKeyBlock, _ := pem.Decode([]byte(privateKeyStr))
	if privateKeyBlock == nil || privateKeyBlock.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("无效的私钥")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// 将字符串转换为RSA公钥
func ParsePublicKey(publicKeyStr string) (*rsa.PublicKey, error) {
	publicKeyBlock, _ := pem.Decode([]byte(publicKeyStr))
	if publicKeyBlock == nil || publicKeyBlock.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("无效的公钥")
	}
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("无效的RSA公钥")
	}
	return rsaPublicKey, nil
}

func GenerateRsaKey(len int) (string, string) {
	if len == 0 {
		len = 2048
	}
	// 生成RSA密钥对
	key, err := rsa.GenerateKey(rand.Reader, len)
	if err != nil {
		fmt.Println("无法生成RSA密钥对：", err)
		panic(err)
	}

	publicKeyStr, err := PublicKeyToString(&key.PublicKey)
	if err != nil {
		panic(err)
	}
	privateKeyStr, err := PrivateKeyToString(key)
	if err != nil {
		panic(err)
	}
	return publicKeyStr, privateKeyStr

}

func RSA_Encrypt(message, pubKey string) string {

	publicKey, err := ParsePublicKey(pubKey)
	if err != nil {
		panic(err)
	}
	// 使用公钥加密消息
	encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(message))
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(encryptedMessage)
}

func RSA_Decrypt(encryptedMsg, priKey string) string {
	privateKey, err := ParsePrivateKey(priKey)
	if err != nil {
		panic(err)
	}
	// 使用私钥解密消息
	cipherText, err := hex.DecodeString(encryptedMsg)
	if err != nil {
		panic(err)
	}
	decryptedMessage, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		panic(err)
	}

	return string(decryptedMessage[:])

}
