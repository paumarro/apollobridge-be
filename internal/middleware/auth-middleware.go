package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
)

var (
	keycloakDomain     = "keycloak-apollo.up.railway.app"
	apollobridgeDomain = "apollobridge.up.railway.app"

	loginPageUrl = fmt.Sprintf(
		"https://%s/realms/apollo/protocol/openid-connect/auth?response_type=code&client_id=apollo-client&redirect_uri=https://%s/auth/callback&scope=openid",
		keycloakDomain,
		apollobridgeDomain,
	)
)

func AuthMiddleware(requiredRole string, clientID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// If Authorization header is missing, try to set it from the cookie
		if authHeader == "" {
			accessToken, err := c.Cookie("access_token")
			if err != nil {
				// Capture the original URL
				originalURL := c.Request.URL.String()
				redirectToLogin(c, originalURL)
				return
			}
			authHeader = "Bearer " + accessToken
			c.Request.Header.Set("Authorization", authHeader)
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
			originalURL := c.Request.URL.String()
			redirectToLogin(c, originalURL)
			return
		}

		if !hasRole(token, requiredRole, clientID) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		// Token is valid, proceed with the request
		c.Next()
	}
}

func hasRole(token *jwt.Token, requiredRole string, clientID string) bool {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Failed to parse claims")
		return false
	}

	resourceAccess, ok := claims["resource_access"].(map[string]interface{})
	if !ok {
		fmt.Println("Failed to extract resource_access")
		return false
	}

	clientRoles, ok := resourceAccess[clientID].(map[string]interface{})["roles"].([]interface{})
	if !ok {
		fmt.Println("Failed to extract client roles")
		return false
	}

	fmt.Println("Client roles found:", clientRoles)

	for _, role := range clientRoles {
		if role == requiredRole {
			fmt.Printf("User has client role %s\n", requiredRole)
			return true
		}
	}
	return false
}

func redirectToLogin(c *gin.Context, originalURL string) {
	loginURL := fmt.Sprintf("%s&state=%s", loginPageUrl, url.QueryEscape(originalURL))
	c.Redirect(http.StatusFound, loginURL)
	c.Abort()
}

func getKeycloakPublicKey() (interface{}, error) {
	jwksURL := "http://keycloak:8080/realms/apollo/protocol/openid-connect/certs"
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
