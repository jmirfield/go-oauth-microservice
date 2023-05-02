package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"oauth/config"
	"os"

	"github.com/golang-jwt/jwt"
)

type Keys struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

func SetupKeyPair(cfg *config.Config) (*Keys, error) {
	keyFile, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read key file: %s", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyFile)
	if err != nil {
		return nil, fmt.Errorf("unabled to parse private key: %s", err)
	}

	return &Keys{Private: privateKey, Public: &privateKey.PublicKey}, nil
}

func (k *Keys) PrivateBytes() []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(k.Private),
	})
}

func (k *Keys) PublicBytes() ([]byte, error) {
	publicKeyPEM, err := x509.MarshalPKIXPublicKey(k.Public)
	if err != nil {
		return nil, fmt.Errorf("error encoding public key: %s", err)
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyPEM,
	}), nil
}
