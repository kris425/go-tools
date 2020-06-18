// @Title	encrypt
// @Description 加密、签名、算法
// @Author	kris
// @CreateTime 	2020-06-16
package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/des"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// @title    		MD5
// @description  	MD5加密
// @auth      		kris	2020-06-16
func MD5(s string, salt []byte) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(salt))
}

// @title    		Sha1
// @description  	SHA1加密
// @auth      		kris	2020-06-16
func Sha1(data []byte, salt []byte) string {
	s := sha1.New()
	s.Write(data)
	return hex.EncodeToString(s.Sum(salt))
}

// @title    		EncryptMD5
// @description  	BASE64加密
// @auth      		kris	2020-06-16
func Base64Encode(data []byte) string {
	s := base64.StdEncoding.EncodeToString(data)
	return s
}

// @title    		EncryptMD5
// @description  	BASE64加密
// @auth      		kris	2020-06-16
func Base64Decode(data string) ([]byte, error) {
	s, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return []byte{}, err
	}
	return s, nil
}

// @title    		AESEncryptECB
// @description  	AES加密，ECB模式，使用PKCS5Padding填充
// @auth      		kris 	2020-06-16
// @Param			src 	待加密密文
// @Param			key 	加密key
func AESEncryptECB(src []byte, key []byte) []byte {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(src) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, src)
	pad := byte(len(plain) - len(src))
	for i := len(src); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted := make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(src); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	return encrypted
}

// @title    		AESDecryptECB
// @description  	AES解密，ECB模式，使用PKCS5Padding填充
// @auth      		kris 	2020-06-16
// @Param			encrypted 加密密文
// @Param			key 加密key
func AESDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}

// @title    		EncryptDESECB
// @description  	AES加密，ECB模式，使用PKCS5Padding填充
// @auth      		kris 	2020-06-16
// @Param
func DESEncryptECB(src, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	src = PKCS5Padding(src, bs)
	if len(src)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

// @title    		DecryptDESECB
// @description  	AES加密，ECB模式，使用PKCS5Padding填充
// @auth      		kris 	2020-06-16
// @Param
func DESDecryptECB(src, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = PKCS5UnPadding(out)
	return out, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func HexEncode(str []byte) string {
	param := ""
	for i := 0; i < len(str); i++ {
		strHex := fmt.Sprintf("%x", str[i]&0xFF)
		if len(strHex) == 1 {
			param = param + "0" + strHex
		} else {
			param = param + strHex
		}
	}
	return strings.Trim(param, "")
}
