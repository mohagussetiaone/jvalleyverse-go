package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"jvalleyverse/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a new JWT token for a user
func GenerateJWT(userID string, role string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "role":    role,
        "exp":     time.Now().Add(config.AppConfig.JWTExpiry).Unix(),
        "iat":     time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

// ParseJWT validates and parses a JWT token, returns userID and role
func ParseJWT(tokenString string) (string, string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(config.AppConfig.JWTSecret), nil
    })
    if err != nil {
        return "", "", err
    }
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID, ok1 := claims["user_id"].(string)
        role, ok2 := claims["role"].(string)
        if !ok1 || !ok2 {
            return "", "", errors.New("invalid claims")
        }
        return userID, role, nil
    }
    return "", "", errors.New("invalid token")
}

// GenerateXSRFToken creates a random string for XSRF protection
func GenerateXSRFToken() string {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        // fallback to timestamp based (not secure but for demo)
        return hex.EncodeToString([]byte(time.Now().String()))
    }
    return hex.EncodeToString(bytes)
}