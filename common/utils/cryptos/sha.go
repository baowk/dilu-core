package cryptos

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// MD5
func MD5(data []byte) string {
	h := md5.Sum(data)
	return hex.EncodeToString(h[:])
}

// sha256
func SHA256(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// 生成密码
func GenPwd(pwd string) (enPwd string, err error) {
	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost); err != nil {
		return
	} else {
		enPwd = string(hash)
	}
	return
}

// 验证密码
func CompPwd(hashPwd, srcPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(srcPwd)); err != nil {
		return false
	}
	return true
}

/*
获取文件的MD5
*/
func MD5File(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
		return ""
	}
	md5 := md5.New()
	io.Copy(md5, f)
	MD5Str := hex.EncodeToString(md5.Sum(nil))
	f.Close()
	return MD5Str
}
