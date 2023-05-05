package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func TestNewValidator(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %s", err)
	}

	validator := NewValidator(&privateKey.PublicKey)

	if err != nil {
		t.Fatalf("Failed to create validator: %s", err)
	}

	if validator == nil {
		t.Fatal("Validator is nil")
	}
}

func TestValidate(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %s", err)
	}

	validator := NewValidator(&privateKey.PublicKey)

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

	validator := NewValidator(&privateKey.PublicKey)

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

	validator := NewValidator(&privateKey.PublicKey)

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

	validator := NewValidator(&privateKey.PublicKey)
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

	validator := NewValidator(&privateKey.PublicKey)
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
