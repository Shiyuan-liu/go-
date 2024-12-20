package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 小写
func Md5Encode(password string) string {

	h := md5.New()
	h.Write([]byte(password))
	temp := h.Sum(nil)
	return hex.EncodeToString(temp)
}

// 大写
func MD5Encode(password string) string {
	return strings.ToUpper(Md5Encode(password))
}

// 加密
func MakePassword(password, salt string) string {
	return MD5Encode(password + salt)
}

// 解密
func ValidPassword(password, salt, makepassword string) bool {
	return MD5Encode(password+salt) == makepassword
}
