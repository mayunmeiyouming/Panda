package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func main() {
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	key := []byte("0123456789ABCDEF")
	plaintext := []byte("Apple is very good")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if false {
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			panic(err.Error())
		}
	}

	fmt.Printf("nonce: %x\n", nonce)

	nr := len(plaintext)
	bufA :=   make([]byte, 2+aesgcm.Overhead()+(16*1024 - 1)+aesgcm.Overhead())
	payloadBuf := bufA[2+aesgcm.Overhead() : 2+aesgcm.Overhead()+(16*1024 - 1)]

	buf := bufA[:2+aesgcm.Overhead()+nr+aesgcm.Overhead()]
	payloadBuf = append(payloadBuf[:0], plaintext[:nr]...)
	buf[0], buf[1] = byte(nr>>8), byte(nr)

	aesgcm.Seal(buf[:0], nonce, buf[:2], nil)
	fmt.Printf("nonce1:%x\n", nonce)
	increment(nonce)
	fmt.Printf("nonce2:%x\n", nonce)
	fmt.Printf("cipher1:%x\n", buf)

	ciphertext := aesgcm.Seal(payloadBuf[:0], nonce, payloadBuf, nil)
	fmt.Printf("nonce3:%x\n", nonce)
	increment(nonce)
	fmt.Printf("nonce4:%x\n", nonce)

	// ciphertext := aesgcm.Seal(plaintext[:0], nonce, plaintext, nil)	

	fmt.Println("nr length: ", nr)
	fmt.Println("Overhead length: ", aesgcm.Overhead())
	fmt.Println("nonce length: ", aesgcm.NonceSize())
	fmt.Printf("plaintext:%s\n", plaintext)
	fmt.Printf("nonce5:%x\n", nonce)
	fmt.Printf("cipher2:%x\n", buf)
	fmt.Printf("ciphertext:%x\n", ciphertext)
}

func increment(b []byte) {
	for i := range b {
		b[i]++
		fmt.Printf("increment:%x\n", b)
		if b[i] != 0 {
			return
		}
	}
}
