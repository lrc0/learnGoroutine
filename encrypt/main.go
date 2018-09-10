package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

func main() {
	plaintext := []byte("matrix:MbJAXoEDsHKW9uwBlR2WkhIYlGdz_1Ti@tcp(172.17.123.185:3306)/bp_matrix?charset=utf8")
	keyText := "abcdefghijklmnopqrstuvwxyzABCDEF"
	block, _ := initCryptoCipher(keyText)

	secret := enCrypto(block, commonIV, plaintext)
	target := deCrypto(block, commonIV, []byte(secret))
	fmt.Println("secret: ", fmt.Sprintf("%x", secret))
	fmt.Println("target: ", fmt.Sprintf("%s", target))
}

func initCryptoCipher(key string) (cipher.Block, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key), err)
		return nil, err
	}
	return c, nil
}

func enCrypto(block cipher.Block, iv, target []byte) []byte {
	streamCFB := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(target))
	streamCFB.XORKeyStream(cipherText, target)
	fmt.Printf("原文:%s => 密文:%x\n", target, cipherText)
	return cipherText
}

func deCrypto(block cipher.Block, iv, target []byte) []byte {
	streamCFB := cipher.NewCFBDecrypter(block, iv)
	cipherText := make([]byte, len(target))
	streamCFB.XORKeyStream(cipherText, target)
	fmt.Printf("密文:%x => 原文:%s\n", target, cipherText)
	return cipherText
}
