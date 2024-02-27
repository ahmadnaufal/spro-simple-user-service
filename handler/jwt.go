package handler

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type JWT struct {
	publicKey  []byte
	privateKey []byte
}

type JWTCustomClaims struct {
	ID          string `json:"id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

func NewJWT(publicKey, privateKey []byte) JWT {
	return JWT{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func (j JWT) CreateToken(payload JWTCustomClaims) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "error parsing private key")
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, payload).SignedString(key)
	if err != nil {
		return "", errors.Wrap(err, "error returning signed string")
	}

	return token, nil
}

func (j JWT) ValidateToken(token string) (string, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse key")
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return "", errors.Wrap(err, "error parsing token")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return "", errors.New("invalid claims")
	}

	return claims["id"].(string), nil
}
