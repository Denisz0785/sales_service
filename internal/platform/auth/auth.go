package auth

import (
	"crypto/rsa"

	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type KeyLookupFunc func(string) (*rsa.PublicKey, error)

// NewSimpleKeyLookupFunc creates a KeyLookupFunc that always returns the same
// public key, identified by the given activeKID.
//
// activeKID: the identifier of the active key.
func NewSimpleKeyLookupFunc(activeKID string, publicKey *rsa.PublicKey) KeyLookupFunc {
	// This function defines the KeyLookupFunc that always returns the same
	// public key, identified by the given activeKID.
	f := func(kid string) (*rsa.PublicKey, error) {
		// If the requested key ID is not the active one, an error is returned.
		if activeKID != kid {
			return nil, fmt.Errorf("unrecognized key: %v", kid)
		}

		// If the requested key ID is the active one, the public key is returned.
		return publicKey, nil
	}
	return f
}

type Authenticator struct {
	PrivateKey       *rsa.PrivateKey
	activeKID        string
	algorithm        string
	pubKeyLookupFunc KeyLookupFunc
	parser           *jwt.Parser
}

func NewAuthenticator(privateKey *rsa.PrivateKey, activeKID string, algorithm string, pubKeyLookupFunc KeyLookupFunc) (*Authenticator, error) {
	if privateKey == nil {
		return nil, errors.New("private key cannot be nil")
	}

	if activeKID == "" {
		return nil, errors.New("active kid cannot be empty")
	}

	if jwt.GetSigningMethod(algorithm) == nil {
		return nil, fmt.Errorf("invalid signing method: %s", algorithm)
	}

	if pubKeyLookupFunc == nil {
		return nil, errors.New("public key lookup function cannot be nil")
	}

	parser := jwt.Parser{
		ValidMethods: []string{algorithm},
	}

	a := &Authenticator{
		PrivateKey:       privateKey,
		activeKID:        activeKID,
		algorithm:        algorithm,
		pubKeyLookupFunc: pubKeyLookupFunc,
		parser:           &parser,
	}
	return a, nil
}

// GenerateToken generates a JWT token using the provided claims and returns the
// token as a string. It uses the Authenticator's private key to sign the token.
func (a *Authenticator) GenerateToken(claims Claims) (string, error) {
	// Get the signing method based on the algorithm specified in the authenticator.
	method := jwt.GetSigningMethod(a.algorithm)

	// Create a new token with the specified method and claims.
	tkn := jwt.NewWithClaims(method, claims)

	// Add the active kid to the token's header.
	tkn.Header["kid"] = a.activeKID

	// Sign the token using the authenticator's private key.
	str, err := tkn.SignedString(a.PrivateKey)

	// If there was an error signing the token, return the error.
	if err != nil {
		return "", errors.Wrap(err, "generating token")
	}

	// Return the token as a string.
	return str, nil
}

// ParseClaims parses a JWT token and returns the claims.
func (a *Authenticator) ParseClaims(tokenString string) (Claims, error) {
	// Parse the token using the authenticator's parser.
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, fmt.Errorf("expecting JWT header to have string kid")
		}
		userKID, ok := kid.(string)
		if !ok {
			return nil, fmt.Errorf("kid must be string")
		}
		return a.pubKeyLookupFunc(userKID)
	}
	var claims Claims
	token, err := a.parser.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		return Claims{}, errors.Wrap(err, "parsing token")
	}

	// If the token is not valid, return an error.
	if !token.Valid {
		return Claims{}, errors.New("token is not valid")
	}

	return claims, nil
}
