package core

import (
	"Panda/utils"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"
)

// Crypto ...
type Crypto interface {
	Encrypt([]byte) []byte
	Decrypt([]byte) []byte
	EncodeWrite(io.ReadWriter, []byte) (int, error)
	DecodeRead(io.ReadWriter, []byte) (int, error)
}

func getCrypt(method byte) Crypto {
	var res Crypto
	switch method {
	case 0x00:
		res = NoCrpt{}
	case 0x80:
		res = AES256{}
	default:
		utils.Logger.Fatal("不支持该加密方式")
	}

	return res
}

// NoCrpt ...
type NoCrpt struct {
}

// Encrypt ...
func (crypto NoCrpt) Encrypt(plaintext []byte) []byte {
	return plaintext
}

// Decrypt ...
func (crypto NoCrpt) Decrypt(ciphertext []byte) []byte {
	return ciphertext
}

// EncodeWrite ...
func (crypto NoCrpt) EncodeWrite(r io.ReadWriter, plaintext []byte) (int, error) {
	res := crypto.Encrypt(plaintext)
	return r.Write(res)
}

// DecodeRead ...
func (crypto NoCrpt) DecodeRead(r io.ReadWriter, plaintext []byte) (int, error) {
	n, err := r.Read(plaintext)
	if err != nil && err != io.EOF {
		return 0, err
	}
	plaintext = crypto.Decrypt(plaintext)

	return n, nil
}

// AES256 ...
type AES256 struct {
}

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

//aes的加密字符串
var key = []byte("astaxie12798akljzmknm.ahkjkljl;k")

// Decrypt ...
func (crypto AES256) Decrypt(ciphertext []byte) []byte {
	// 创建加密算法aes
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}

	// 解密字符串
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plaintextCopy := make([]byte, len(ciphertext))
	cfbdec.XORKeyStream(plaintextCopy, ciphertext)
	fmt.Printf("%x=>%s\n", ciphertext, plaintextCopy)

	return plaintextCopy
}

// Encrypt ...
func (crypto AES256) Encrypt(plaintext []byte) []byte {
	// 创建加密算法aes
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key), err)
		os.Exit(-1)
	}

	//加密字符串
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	fmt.Printf("%s=>%x\n", plaintext, ciphertext)

	return ciphertext
}

// EncodeWrite ...
func (crypto AES256) EncodeWrite(r io.ReadWriter, plaintext []byte) (int, error) {
	res := crypto.Encrypt(plaintext)
	return r.Write(res)
}

// DecodeRead ...
func (crypto AES256) DecodeRead(r io.ReadWriter, plaintext []byte) (int, error) {
	n, err := r.Read(plaintext)
	if err != nil && err != io.EOF {
		return 0, err
	}
	plaintext = crypto.Decrypt(plaintext)

	return n, nil
}
