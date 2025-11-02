package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

type AccessClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

func LoadRSAPublicKey() *rsa.PublicKey {
	data, err := os.ReadFile("pkg/middleware/rsa_public.pem")
	if err != nil {
		panic(fmt.Errorf("failed to read public key file: %w", err))
	}

	block, _ := pem.Decode(data)
	if block == nil {
		panic(fmt.Errorf("failed to parse PEM block from public key file"))
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(fmt.Errorf("failed to parse RSA public key: %w", err))
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		panic(fmt.Errorf("key is not RSA public key"))
	}

	return rsaPub
}

func JWTAccessMiddleware(rsaPubKey *rsa.PublicKey, metrics *Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken := c.GetHeader("Authorization")
		if jwtToken == "" {
			metrics.JWTValidationTotal.WithLabelValues("no_token").Inc()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}
		jwtToken = jwtToken[len("Bearer "):]

		token, err := jwt.ParseWithClaims(jwtToken, &AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return rsaPubKey, nil
		})
		if err != nil {
			metrics.JWTValidationTotal.WithLabelValues("invalid_claims").Inc()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*AccessClaims)
		if !ok || !token.Valid {
			metrics.JWTValidationTotal.WithLabelValues("invalid_claims").Inc()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		metrics.JWTValidationTotal.WithLabelValues("success").Inc()
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
