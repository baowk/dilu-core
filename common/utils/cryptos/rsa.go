package cryptos

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

// 将RSA私钥转换为byte
func PrivateKeyToPem(privateKey *rsa.PrivateKey) []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return pem.EncodeToMemory(privateKeyBlock)
}

// 将RSA公钥转换为byte
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

// 将byte转为私钥
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

// 将byte转换为RSA公钥
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

// 将RSAP kcs8私钥转换为byte
func PrivateKeyPkcs8ToPem(privateKey *rsa.PrivateKey) ([]byte, error) {
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return pem.EncodeToMemory(privateKeyBlock), err
}

// 将Pkcs8 格式byte转为私钥
func ParsePriKeyPkcs8(privateKeyPkcs8 []byte) (*rsa.PrivateKey, error) {
	privateKeyBlock, _ := pem.Decode(privateKeyPkcs8)
	if privateKeyBlock == nil {
		return nil, fmt.Errorf("Pkcs8无效的私钥")
	}
	// if privateKeyBlock == nil || privateKeyBlock.Type != "RSA PRIVATE KEY" {
	// 	return nil, fmt.Errorf("Pkcs8无效的私钥")
	// }
	priKey, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, ok := priKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("Pkcs8 contained non-RSA key. Expected RSA key.")
	}
	return rsaPrivateKey, nil
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

func EncodeToString(data []byte) string {
	return hex.EncodeToString(data)
}

func DecodeString(message string) ([]byte, error) {
	return hex.DecodeString(message)
}

// 公钥加密
func RSA_Encrypt(message []byte, pubKey string) (string, error) {
	b, err := RsaEncrypt(message, []byte(pubKey))
	if err != nil {
		return "", err
	}
	return EncodeToString(b), nil
}

// 公钥加密
func RsaEncrypt(message []byte, pubKey []byte) ([]byte, error) {
	publicKey, err := ParsePubKey(pubKey)
	if err != nil {
		return nil, err
	}
	// 使用公钥加密消息
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
}

// 私钥解密
func RSA_Decrypt(encryptedMsg, priKey string) ([]byte, error) {
	b, err := DecodeString(encryptedMsg)
	if err != nil {
		return nil, err
	}
	return RsaDecrypt(b, []byte(priKey))
}

// 私钥解密Pkcs8 Key
func RSA_DecryptPkcs8(encryptedMsg, priKey string) ([]byte, error) {
	b, err := DecodeString(encryptedMsg)
	if err != nil {
		return nil, err
	}
	return RsaDecryptPkcs8(b, []byte(priKey))
}

// 私钥解密 Pkcs8 Key
func RsaDecryptPkcs8(encryptedMsg, priKey []byte) ([]byte, error) {
	privateKey, err := ParsePriKeyPkcs8(priKey)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedMsg)
}

// 私钥解密
func RsaDecrypt(encryptedMsg, priKey []byte) ([]byte, error) {
	privateKey, err := ParsePriKey(priKey)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedMsg)
}

// RsaSign 私钥加签
func RSA_Sign(priKey string, message []byte) ([]byte, error) {
	return RsaSign([]byte(priKey), message)
}

// RsaSign 私钥加签
func RSA_SignPkcs8(priKey string, message []byte) ([]byte, error) {
	return RsaSignPkcs8([]byte(priKey), message)
}

// RsaSign 私钥加签
func RsaSign(priKey, message []byte) ([]byte, error) {
	privateKey, _ := ParsePriKey(priKey)
	return RsaSignKey(privateKey, message)
}

// RsaSign Pkcs8私钥加签
func RsaSignPkcs8(priKey, message []byte) ([]byte, error) {
	privateKey, _ := ParsePriKeyPkcs8(priKey)
	return RsaSignKey(privateKey, message)
}

// 根据私钥加签
func RsaSignKey(pkey *rsa.PrivateKey, message []byte) ([]byte, error) {
	hash := sha256.Sum256(message)
	return rsa.SignPKCS1v15(rand.Reader, pkey, crypto.SHA256, hash[:])
}

func RsaSignWithHash(pkey *rsa.PrivateKey, message []byte, algorithm uint16) ([]byte, error) {
	switch algorithm {
	case 1:
		hash := sha1.Sum(message)
		return rsa.SignPKCS1v15(rand.Reader, pkey, crypto.SHA1, hash[:])
	case 224:
		h224 := crypto.SHA224.New()
		h224.Write(message)
		hash := h224.Sum(nil)
		return rsa.SignPKCS1v15(rand.Reader, pkey, crypto.SHA224, hash[:])
	case 256:
		hash := sha256.Sum256(message)
		return rsa.SignPKCS1v15(rand.Reader, pkey, crypto.SHA256, hash[:])
	case 384:
		h384 := crypto.SHA384.New()
		h384.Write(message)
		hash := h384.Sum(nil)
		return rsa.SignPKCS1v15(rand.Reader, pkey, crypto.SHA384, hash[:])
	case 512:
		hash := sha512.Sum512(message)
		return rsa.SignPKCS1v15(rand.Reader, pkey, crypto.SHA512, hash[:])
	}
	return nil, errors.New("不支持的hash算法")
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

func RsaVerifyWithHash(publicKeyPEM, message []byte, signature []byte, algorithm uint16) error {
	publicKey, _ := ParsePubKey(publicKeyPEM)
	switch algorithm {
	case 1:
		hash := sha1.Sum(message)
		return rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hash[:], signature)
	case 224:
		h224 := crypto.SHA224.New()
		h224.Write(message)
		hash := h224.Sum(nil)
		return rsa.VerifyPKCS1v15(publicKey, crypto.SHA224, hash[:], signature)
	case 256:
		hash := sha256.Sum256(message)
		return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
	case 384:
		h384 := crypto.SHA384.New()
		h384.Write(message)
		hash := h384.Sum(nil)
		return rsa.VerifyPKCS1v15(publicKey, crypto.SHA384, hash[:], signature)
	case 512:
		hash := sha512.Sum512(message)
		return rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, hash[:], signature)
	}
	return errors.New("不支持的hash算法")
}

func RsaPriKeyPkcs8To1(priPkcs8Key []byte) (string, error) {
	pk, err := ParsePriKeyPkcs8(priPkcs8Key)
	if err != nil {
		return "", err
	}
	return string(PrivateKeyToPem(pk)), nil
}

func RsaPriKeyPkcs1To8(priPkcs1Key []byte) (string, error) {
	pk, err := ParsePriKey(priPkcs1Key)
	if err != nil {
		return "", err
	}
	b, err := PrivateKeyPkcs8ToPem(pk)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func RsaKeyFmt(key string) string {
	if !strings.Contains(key, "-----\n") {
		key = strings.Replace(key, "-----", "-----\n", 1)
	}
	if !strings.Contains(key, "\n-----END") {
		key = strings.Replace(key, "-----END", "\n-----END", -1)
	}
	if strings.Contains(key, "\t") {
		return strings.Replace(key, "\t", "", -1)
	}
	return key
	// if strings.Contains(key, "\n\r") {
	// 	return key, nil
	// }
	// if strings.Contains(key, "\r") {
	// 	return strings.Replace(key, "\r", "\n\r", -1), nil
	// }
	// if strings.Contains(key, "-----") {
	// 	fk := ""
	// 	arr := strings.Split(key, "-----")
	// 	for i := 0; i < len(arr); i++ {
	// 		if arr[i] == "" {
	// 			continue
	// 		} else if strings.HasPrefix(strings.ToUpper(arr[i]), "BEGIN") {
	// 			fk += "-----" + arr[i] + "-----\n\r"
	// 		} else if len(arr[i]) > 64 {
	// 			cnt := len(arr[i]) / 64
	// 			for j := 0; j < cnt; j++ {
	// 				fk += arr[i][j*64:(j+1)*64] + "\n\r"
	// 			}
	// 			if len(arr[i])%64 != 0 {
	// 				fk += arr[i][cnt*64:] + "\n\r"
	// 			}
	// 		} else {
	// 			fk += "-----" + arr[i] + "-----"
	// 		}
	// 	}
	// 	return fk, nil
	// }
	// return key, nil
}
