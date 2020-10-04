package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

func main() {

	key := []byte("0123456789ABCDEF")

	ciphertext, _ := hex.DecodeString("08f24c28f0fc9aef5812a35ce66235bc2488d6c29b") //加密生成的结果

	nonce, _ := hex.DecodeString("000000000000000000000000") //加密用的nonce

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(plaintext))
}
