/*
 * This file was last modified at 2024-07-29 13:43 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * load_rsa_key_test.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	cRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	mRand "math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadKeys(t *testing.T) {

	rnd := mRand.New(mRand.NewSource(time.Now().UnixNano()))
	id := rnd.Uint32() % 9999

	privateKey, publicKey := generateRsaKeyPair()

	testPrivateKeyFileName := fmt.Sprintf("%s/test_private_key_%04d.pem", os.TempDir(), id)
	privateStr := exportRsaPrivateKeyAsPemStr(privateKey)
	privateKeyFile, err := os.OpenFile(testPrivateKeyFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0640)
	IfErrorThenPanic(err)
	defer FileClose(privateKeyFile)
	_, err = privateKeyFile.Write([]byte(privateStr))
	IfErrorThenPanic(err)
	rsaPrivateKey := LoadPrivateKey(testPrivateKeyFileName)
	assert.Equal(t, privateKey, rsaPrivateKey)
	_ = os.RemoveAll(testPrivateKeyFileName)

	testPublicKeyFileName := fmt.Sprintf("%s/test_public_key_%04d.pem", os.TempDir(), id)
	publicKeyStr, err := exportRsaPublicKeyAsPemStr(publicKey)
	IfErrorThenPanic(err)
	publicKeyFile, err := os.OpenFile(testPublicKeyFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0640)
	IfErrorThenPanic(err)
	defer FileClose(publicKeyFile)
	_, err = publicKeyFile.Write([]byte(publicKeyStr))
	IfErrorThenPanic(err)
	rsaPublicKey := LoadPublicKey(testPublicKeyFileName)
	assert.Equal(t, publicKey, rsaPublicKey)
	_ = os.RemoveAll(testPublicKeyFileName)
}

func exportRsaPrivateKeyAsPemStr(privateKey *rsa.PrivateKey) string {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)
	return string(privateKeyPEM)
}

func exportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(pubkey)
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)
	return string(publicKeyPEM), nil
}

func generateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {

	privateKey, err := rsa.GenerateKey(cRand.Reader, 1024)
	IfErrorThenPanic(err)
	return privateKey, &privateKey.PublicKey
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
