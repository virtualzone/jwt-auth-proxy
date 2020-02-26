package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

type Cert struct {
	PrivateKey  *rsa.PrivateKey
	Certificate *x509.Certificate
	CertBytes   []byte
}

func CertCreateCA() (*Cert, error) {
	serial, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, err
	}
	ca := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"Virtualzone"},
			Country:      []string{"DE"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}
	res := &Cert{
		PrivateKey:  caPrivKey,
		Certificate: ca,
		CertBytes:   caBytes,
	}
	return res, nil
}

func CertCreateSign(caCert *Cert) (*Cert, error) {
	serial, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, err
	}
	cert := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"Virtualzone"},
			Country:      []string{"DE"},
		},
		IPAddresses:  GetConfig().BackendCertIPs,
		DNSNames:     GetConfig().BackendCertHostnames,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert.Certificate, &certPrivKey.PublicKey, caCert.PrivateKey)
	if err != nil {
		return nil, err
	}
	res := &Cert{
		PrivateKey:  certPrivKey,
		Certificate: cert,
		CertBytes:   certBytes,
	}
	return res, nil
}

func (cert *Cert) SavePrivateKey(fileName string) error {
	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outFile.Close()
	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(cert.PrivateKey),
	}
	err = pem.Encode(outFile, privateKey)
	if err != nil {
		return err
	}
	return nil
}

func (cert *Cert) SavePublicKey(fileName string) error {
	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outFile.Close()
	var publicKey = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&cert.PrivateKey.PublicKey),
	}
	err = pem.Encode(outFile, publicKey)
	if err != nil {
		return err
	}
	return nil
}

func (cert *Cert) SaveCertificate(fileName string) error {
	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outFile.Close()
	err = pem.Encode(outFile, &pem.Block{Type: "CERTIFICATE", Bytes: cert.CertBytes})
	if err != nil {
		return err
	}
	return nil
}
