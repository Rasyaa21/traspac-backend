package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"gin-backend-app/internal/models"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
    UserID uuid.UUID `json:"user_id"`
    Email  string    `json:"email"`
    Name   string    `json:"name"`
    jwt.RegisteredClaims
}

type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
}

var (
    jwtSecretBytes []byte
    jwtOnce        sync.Once
)

func initJWTSecret() {
    jwtOnce.Do(func() {
        secret := os.Getenv("JWT_SECRET")
        if secret == "" {
            // Development fallback - CHANGE IN PRODUCTION
            secret = "your_super_secret_jwt_key_change_in_production_min_32_characters"
        }
        jwtSecretBytes = []byte(secret)
    })
}

// Generate access token for user
func GenerateAccessToken(user *models.User) (string, error) {
    initJWTSecret()

    if len(jwtSecretBytes) == 0 {
        return "", errors.New("JWT_SECRET not set")
    }

    expireTime := time.Now().Add(24 * time.Hour) // 24 hours for access token

    claims := Claims{
        UserID: user.ID,
        Email:  user.Email,
        Name:   user.Name,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expireTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "traspac-backend",
            Subject:   user.ID.String(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtSecretBytes)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// Generate refresh token (random string)
func GenerateRefreshToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

// Generate both access and refresh tokens
func GenerateTokenPair(user *models.User) (*TokenPair, error) {
    accessToken, err := GenerateAccessToken(user)
    if err != nil {
        return nil, err
    }

    refreshToken, err := GenerateRefreshToken()
    if err != nil {
        return nil, err
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    24 * 60 * 60, // 24 hours in seconds
    }, nil
}

// Legacy function for backward compatibility
func GenerateToken(user *models.User) (string, error) {
    return GenerateAccessToken(user)
}

// CreateToken alias for consistency
func CreateToken(user *models.User) (string, error) {
    return GenerateAccessToken(user)
}

// Validate JWT token
func ValidateToken(tokenString string) (*Claims, error) {
    initJWTSecret()
    if len(jwtSecretBytes) == 0 {
        return nil, errors.New("JWT_SECRET not set")
    }

    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Validate signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("invalid signing method")
        }
        return jwtSecretBytes, nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}

// Get claims from Gin context Authorization header
func GetClaimsFromHeader(c *gin.Context) (*Claims, error) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        return nil, errors.New("authorization header is required")
    }

    parts := strings.SplitN(authHeader, " ", 2)
    if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || parts[1] == "" {
        return nil, errors.New("authorization header format must be Bearer {token}")
    }

    return ValidateToken(parts[1])
}

// Helper functions for Gin context
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
    userID, exists := c.Get("user_id")
    if !exists {
        return uuid.Nil, errors.New("user not authenticated")
    }

    id, ok := userID.(uuid.UUID)
    if !ok {
        return uuid.Nil, errors.New("invalid user ID type")
    }

    return id, nil
}

func GetUserEmailFromContext(c *gin.Context) (string, error) {
    email, exists := c.Get("user_email")
    if !exists {
        return "", errors.New("user email not found")
    }

    emailStr, ok := email.(string)
    if !ok {
        return "", errors.New("invalid email type")
    }

    return emailStr, nil
}

func GetUserNameFromContext(c *gin.Context) (string, error) {
    name, exists := c.Get("user_name")
    if !exists {
        return "", errors.New("user name not found")
    }

    nameStr, ok := name.(string)
    if !ok {
        return "", errors.New("invalid name type")
    }

    return nameStr, nil
}