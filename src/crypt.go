package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func Encrypt(passphrase, s string) (string, error) {
	c, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	res := gcm.Seal(nonce, nonce, []byte(s), nil)
	return base64.StdEncoding.EncodeToString(res), nil
}

func Decrypt(passphrase, s string) (string, error) {
	s2, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	c, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(s2) < nonceSize {
		return "", err
	}
	nonce, s2 := s2[:nonceSize], s2[nonceSize:]
	res, err := gcm.Open(nil, []byte(nonce), []byte(s2), nil)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
