package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(userID int, ttl time.Duration, secretJWTKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = userID
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := token.SignedString([]byte(secretJWTKey))
	if err != nil {
		return "", fmt.Errorf("generating JWT Token failed: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(token string, signedJWTKey string) (int, error) {
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return []byte(signedJWTKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("invalidate token: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return 0, fmt.Errorf("invalid token claim")
	}

	userIDFloat64, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in token")
	}

	userID := int(userIDFloat64)

	return userID, nil
}
