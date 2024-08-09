/*
 * This file was last modified at 2024-07-29 15:38 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * load_rsa_key.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"
)

func LoadPrivateKey(fileName string) *rsa.PrivateKey {
	if len(fileName) > 1 {
		file, err := os.Open(fileName)
		if err != nil {
			return nil
		}
		defer FileClose(file)
		buf, err := io.ReadAll(file)
		if err != nil {
			return nil
		}
		if block := readPEMString(string(buf)); block != nil {
			privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil
			}
			return privateKey
		}
	}
	return nil
}

func LoadPublicKey(fileName string) *rsa.PublicKey {
	if len(fileName) > 1 {
		file, err := os.Open(fileName)
		if err != nil {
			return nil
		}
		defer FileClose(file)
		buf, err := io.ReadAll(file)
		if err != nil {
			return nil
		}
		if block := readPEMString(string(buf)); block != nil {
			publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				return nil
			}
			return publicKey
		}
	}
	return nil
}

func readPEMString(p string) *pem.Block {
	result, _ := pem.Decode([]byte(p))
	return result
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
