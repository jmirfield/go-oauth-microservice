package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

func GetPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	keyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read key file: %s", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyFile)
	if err != nil {
		return nil, fmt.Errorf("unabled to parse private key: %s", err)
	}

	return privateKey, nil
}

func PrivateBytes(p *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(p),
	})
}

func PublicBytes(p *rsa.PublicKey) ([]byte, error) {
	publicKeyPEM, err := x509.MarshalPKIXPublicKey(p)
	if err != nil {
		return nil, fmt.Errorf("error encoding public key: %s", err)
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyPEM,
	}), nil
}
