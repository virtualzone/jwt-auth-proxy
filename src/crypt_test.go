package main

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	s := "This is a test"
	passphrase := "nNoh417qsUa0pdTA"

	res, err := Encrypt(passphrase, s)
	if err != nil {
		t.Fatalf("Encryption error: %s", err)
	}
	checkStringNotEmpty(t, res)

	res, err = Decrypt(passphrase, res)
	if err != nil {
		t.Fatalf("Encryption error: %s", err)
	}
	checkTestString(t, s, res)
}

func TestEncryptDecryptShortString(t *testing.T) {
	s := "T"
	passphrase := "nNoh417qsUa0pdTA"

	res, err := Encrypt(passphrase, s)
	if err != nil {
		t.Fatalf("Encryption error: %s", err)
	}
	checkStringNotEmpty(t, res)

	res, err = Decrypt(passphrase, res)
	if err != nil {
		t.Fatalf("Encryption error: %s", err)
	}
	checkTestString(t, s, res)
}

func TestEncryptDecryptShortPassphrase(t *testing.T) {
	s := "This is a test"
	passphrase := "123456789012345"

	_, err := Encrypt(passphrase, s)
	if err == nil {
		t.Fatal("Expected minimum error due to short key length (15 bytes)")
	}
}
