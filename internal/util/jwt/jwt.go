package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/TRIOSYS-Software/auth-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	jwtKeyID    string
	jwtIssuer   string
	jwtAudience string
)

func InitializeJWT(privateKeyPath, publicKeyPath, keyID, issuer, audience string) error {
	privBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}
	pubBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return err
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return err
	}
	jwtKeyID = keyID
	jwtIssuer = issuer
	jwtAudience = audience

	return nil
}

func GenerateToken(userID int64, username, email string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = jwtKeyID

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(config.Cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}
	return claims, nil
}

func ServeJWKS() (map[string]any, error) {
	log := config.Cfg.Logger
	if publicKey == nil || jwtKeyID == "" {
		log.Error("Public key or Key ID not initialized for JWKS")
		return nil, fmt.Errorf("public key or key ID not initialized")
	}
	n := base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes())
	jwk := map[string]any{
		"kty": "RSA",
		"kid": jwtKeyID,
		"use": "sig",
		"alg": "RS256",
		"n":   n,
		"e":   e,
	}
	jwks := map[string]any{
		"keys": []any{jwk},
	}
	log.Info("ServeJWKS(): JWKS served successfully")
	return jwks, nil
}
