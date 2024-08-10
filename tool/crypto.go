/*
 * This file was last modified at 2024-07-29 19:00 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * crypto.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

var (
	ErrEncryptAES = fmt.Errorf("encrypt with AES")
	ErrEncryptRSA = fmt.Errorf("encrypt with RSA")
	ErrDecryptRSA = fmt.Errorf("decrypt with RSA")
)

func EncryptAES(secretKey, plain []byte) ([]byte, error) {

	cipherBlock, err := aes.NewCipher(secretKey)

	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(plain))
	cipherBlock.Encrypt(ciphertext, plain)

	return ciphertext, nil
}

func DecryptAES(secretKey, bytes []byte) ([]byte, error) {

	cipherBlock, err := aes.NewCipher(secretKey)

	if err != nil {
		return nil, err
	}
	plain := make([]byte, len(bytes))
	cipherBlock.Decrypt(plain, bytes)

	return plain, nil
}

func EncryptRSA(rsaPublicKey *rsa.PublicKey, plain []byte) ([]byte, error) {

	if rsaPublicKey != nil {
		if result, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, plain); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
	return nil, ErrEncryptRSA
}

func DecryptRSA(rsaPrivateKey *rsa.PrivateKey, bytes []byte) ([]byte, error) {

	if rsaPrivateKey != nil {
		if result, err := rsa.DecryptPKCS1v15(nil, rsaPrivateKey, bytes); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
	return nil, ErrDecryptRSA
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
