package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func TestNewValidator(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %s", err)
	}
	publicKeyPEM, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("Error encoding public key: %s", err)
	}

	validator, err := NewValidator(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyPEM,
	}))
	if err != nil {
		t.Fatalf("Failed to create validator: %s", err)
	}

	if validator == nil {
		t.Fatal("Validator is nil")
	}
}

func TestNewValidatorWithInvalidPublicKey(t *testing.T) {
	// Creating a non-RSA public key
	invalidPublicKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate invalid public key: %s", err)
	}

	// Marshaling the public key to PEM format
	invalidPublicKeyPEM, err := x509.MarshalPKIXPublicKey(&invalidPublicKey.PublicKey)
	if err != nil {
		t.Fatalf("Error encoding invalid public key: %s", err)
	}

	// Trying to create a validator with an invalid public key
	_, err = NewValidator(pem.EncodeToMemory(&pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: invalidPublicKeyPEM,
	}))

	if err == nil {
		t.Fatal("Expected error but got nil")
	}
}

func TestValidate(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %s", err)
	}

	publicKeyPEM, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("Error encoding public key: %s", err)
	}

	validator, err := NewValidator(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyPEM,
	}))

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "1234567890",
		"iat": 1516239022,
		"exp": time.Now().Add(1 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("Failed to sign token: %s", err)
	}

	validatedToken, err := validator.Validate(tokenString)
	if err != nil {
		t.Fatalf("Failed to validate token: %s", err)
	}

	if validatedToken == nil {
		t.Fatal("Validated token is nil")
	}

	if !validatedToken.Valid {
		t.Fatal("Validated token is not valid")
	}

	claims, ok := validatedToken.Claims.(*jwt.StandardClaims)
	if !ok {
		t.Fatal("Failed to get claims from validated token")
	}

	if sub := claims.Subject; !ok || sub != "1234567890" {
		t.Fatalf("Failed to get sub claim from validated token: %v", sub)
	}

	if iat := claims.IssuedAt; !ok || iat != 1516239022 {
		t.Fatalf("Failed to get iat claim from validated token: %v", iat)
	}
}

func TestValidateWithExpiredToken(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %s", err)
	}

	publicKeyPEM, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("Error encoding public key: %s", err)
	}

	validator, err := NewValidator(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyPEM,
	}))

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "1234567890",
		"iat": 1516239022,
		"exp": 1516239022,
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("Failed to sign token: %s", err)
	}

	_, err = validator.Validate(tokenString)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestValidateWithInvalidClaims(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %s", err)
	}

	publicKeyPEM, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("Error encoding public key: %s", err)
	}

	validator, err := NewValidator(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyPEM,
	}))

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"invalid": "1234567890",
	})

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("Failed to sign token: %s", err)
	}

	_, err = validator.Validate(tokenString)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestValidateWithInvalidToken(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %s", err)
	}

	publicKeyPEM, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("Error encoding public key: %s", err)
	}

	validator, err := NewValidator(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyPEM,
	}))
	if err != nil {
		t.Fatalf("Failed to create validator: %s", err)
	}

	_, err = validator.Validate("invalid-token")
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
}

func TestValidateWrongAlgorithm(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %s", err)
	}

	publicKeyPEM, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("Error encoding public key: %s", err)
	}

	validator, err := NewValidator(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyPEM,
	}))
	if err != nil {
		t.Fatalf("Failed to create validator: %s", err)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("Failed to sign token: %s", err)
	}

	_, err = validator.Validate(tokenString)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
