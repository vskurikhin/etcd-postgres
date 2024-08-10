/*
 * This file was last modified at 2024-07-29 16:49 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * crypto_test.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypto(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "positive test #0 AES",
			fRun: testAESPositiveCase,
		},
		{
			name: "negative test #1 AES case #1",
			fRun: testAESNegativeCase1,
		},
		{
			name: "negative test #2 AES case #2",
			fRun: testAESNegativeCase2,
		},
		{
			name: "positive test #3 RSA",
			fRun: testRSAPositiveCase,
		},
		{
			name: "negative test #4 RSA",
			fRun: testRSANegativeCase,
		},
	}
	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testAESPositiveCase(t *testing.T) {
	secret := make([]byte, 32) // 32 bytes to select AES-256.
	if _, err := rand.Reader.Read(secret); err != nil {
		t.Fail()
	}
	expected := "This information"
	encrypt, err := EncryptAES(secret, []byte(expected))
	assert.Nil(t, err)
	got, err := DecryptAES(secret, encrypt)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(got))
}

func testAESNegativeCase1(t *testing.T) {
	expected := "test"
	encrypt, err := EncryptAES(nil, []byte(expected))
	assert.NotNil(t, err)
	got, err := DecryptAES(nil, encrypt)
	assert.NotNil(t, err)
	assert.NotEqual(t, expected, string(got))
}

func testAESNegativeCase2(t *testing.T) {
	secret := make([]byte, 15)
	expected := "test"
	encrypt, err := EncryptAES(secret, []byte(expected))
	assert.NotNil(t, err)
	assert.Equal(t, aes.KeySizeError(15), err)
	got, err := DecryptAES(secret, encrypt)
	assert.NotNil(t, err)
	assert.NotEqual(t, expected, string(got))
}

func testRSAPositiveCase(t *testing.T) {
	expected := "Supercalifragilisticexpialidocio"
	encrypt, err := EncryptRSA(LoadPublicKey("test_public-key.pem"), []byte(expected))
	assert.Nil(t, err)
	b64 := base64.StdEncoding.EncodeToString(encrypt)
	_, _ = fmt.Fprintf(os.Stderr, "b64: %s\n", b64)
	got, err := DecryptRSA(LoadPrivateKey("test_private-key.pem"), encrypt)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(got))
}

func testRSANegativeCase(t *testing.T) {
	expected := ""
	encrypt, err := EncryptRSA(nil, []byte(expected))
	assert.NotNil(t, err)
	got, err := DecryptRSA(nil, encrypt)
	assert.NotNil(t, err)
	assert.Equal(t, expected, string(got))
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
