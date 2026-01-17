package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashedPassword returns the bcrypt hash of the password at cost 12
func HashedPassword(password string) (string, error) {
	// get the hash string from bcrypt
	passHashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(passHashed), nil
}

// checks user inputted password against stored hashed password
func CheckPasswordvsHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// create a Json Web Token (JWT) for user authentication
func CreateJWT(userID uuid.UUID, tokenSecret string, expiration time.Duration) (string, error) {
	// create token payload with standard claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "doc_appointment",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiration)),
		Subject:   userID.String(),
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

// ValidateJWT verifies signature & standard claims, then returns the userID from Subject.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		// Ensure we actually got HS256 (or a method you expect)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}

		// converts the string key to the appropriate type to use by the hash method to compare outputs
		return []byte(tokenSecret), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	// get the user uuid
	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("could not parse id")
	}
	return userId, nil
}

// retreive the JWT token sent by the client
func GetBearerToken(headers http.Header) (string, error) {
	//get the value sign the key header name
	value := headers.Get("Authorization")
	if value == "" {
		return "", errors.New("Authorization header missing")
	}

	// check for prefix
	const prefix = "Bearer "
	if !strings.HasPrefix(value, prefix) {
		return "", errors.New("authorization header format must be 'Bearer <token>'")
	}

	// trim prefix
	tokenStr := strings.TrimSpace(strings.TrimPrefix(value, prefix))
	if tokenStr == "" {
		return "", errors.New("token missing after 'Bearer '")
	}

	return tokenStr, nil
}

// function to create encoded 32-byte string for refresh token
func MakeRefreshToken() (string, error) {
	// create random 32-byte
	key := make([]byte, 32)
	// fills key with random data
	rand.Read(key)
	// encode byte to hex string
	token := hex.EncodeToString(key)

	return token, nil
}

// get API key from Webhook event request header
func GetAPIKey(headers http.Header) (string, error) {
	// get header value
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("Authorization header missing")
	}

	// check prefix value
	prefix := "ApiKey "
	if !strings.HasPrefix(val, prefix) {
		return "", errors.New("authorization header format must be 'ApiKey <token>'")
	}

	// trim prefix
	keyStr := strings.TrimSpace(strings.TrimPrefix(val, prefix))
	if keyStr == "" {
		return "", errors.New("token missing after 'ApiKey '")
	}

	return keyStr, nil
}
