package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Validator -
type Validator struct {
	key *rsa.PublicKey
}

// NewValidator -
func NewValidator(publicKey []byte) (*Validator, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse key: %s", err)
	}

	return &Validator{
		key: key,
	}, nil
}

// Validate validates a JWT
func (v *Validator) Validate(token string) (*jwt.Token, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return v.key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %s", err)
	}

	claims, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token has expired")
	}

	return t, nil
}
