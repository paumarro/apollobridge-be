package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
)

var loginPageUrl = "http://localhost:8080/realms/apollo/protocol/openid-connect/auth?response_type=code&client_id=apollo-client&redirect_uri=http://localhost:3000/artworks&scope=openid"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			redirectToLogin(c)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		pubKey, err := getKeycloakPublicKey()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get public key"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return pubKey, nil
		})

		if err != nil || !token.Valid {
			redirectToLogin(c)
			return
		}

		// Token is valid, proceed with the request
		c.Next()
	}
}

func getKeycloakPublicKey() (interface{}, error) {
	jwksURL := "http://localhost:8080/auth/realms/apollo/protocol/openid-connect/certs"
	set, err := jwk.Fetch(context.Background(), jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWK set: %v", err)
	}

	key, ok := set.Get(0)
	if !ok {
		return nil, fmt.Errorf("failed to get key from JWK set")
	}

	var pubKey interface{}
	if err := key.Raw(&pubKey); err != nil {
		return nil, fmt.Errorf("failed to get raw key: %v", err)
	}

	return pubKey, nil
}

func redirectToLogin(c *gin.Context) {
	c.Redirect(http.StatusFound, loginPageUrl)
	c.Abort()
}
