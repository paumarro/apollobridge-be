package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
)

var (
	kcDomain           = os.Getenv("KEYCLOAK_DOMAIN")
	apollobridgeDomain = os.Getenv("APOLLO_DOMAIN")
	jwksURL            = os.Getenv("JWKS_URL")
	loginPageUrl       = fmt.Sprintf(
		"https://%s/realms/apollo/protocol/openid-connect/auth?response_type=code&client_id=apollo-client&redirect_uri=https://%s/auth/callback&scope=openid",
		kcDomain,
		apollobridgeDomain,
	)
)

func refreshAccessToken(refreshToken string) (string, string, error) {
	kcClientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	kcClientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	kcDomain := os.Getenv("KEYCLOAK_DOMAIN")

	tokenURL := fmt.Sprintf("https://%s/realms/apollo/protocol/openid-connect/token", kcDomain)

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", kcClientID)
	data.Set("client_secret", kcClientSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", fmt.Errorf("failed to create refresh token request: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to refresh token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("failed to refresh token: %s", string(bodyBytes))
	}

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", "", fmt.Errorf("failed to decode refresh token response: %v", err)
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid refresh token response: missing access token")
	}

	newRefreshToken, ok := tokenResponse["refresh_token"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid refresh token response: missing refresh token")
	}

	return accessToken, newRefreshToken, nil
}

func Auth(requiredRole string, clientID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// If Authorization header is missing, try to set it from the cookie
		if authHeader == "" {
			accessToken, err := c.Cookie("access_token")
			if err != nil {
				fmt.Println("Error fetching access_token cookie:", err)
				originalURL := c.Request.URL.String()
				redirectToLogin(c, originalURL)
				return
			}

			authHeader = "Bearer " + accessToken
			c.Request.Header.Set("Authorization", authHeader)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Mitigation: Validate length of token to avoid excessive memory allocation
		if len(tokenString) > 2024 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "JWT too large"})
			return
		}

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodRS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return getKeycloakPublicKey(token)
		})

		if err != nil || !token.Valid {
			// Check if the error is due to token expiry
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
					fmt.Println("Access token expired, attempting to refresh")

					// Fetch the refresh token
					refreshToken, err := c.Cookie("refresh_token")
					if err != nil {
						fmt.Println("Error fetching refresh_token cookie:", err)
						redirectToLogin(c, c.Request.URL.String())
						return
					}

					// Attempt to refresh the token
					newAccessToken, newRefreshToken, err := refreshAccessToken(refreshToken)
					if err != nil {
						fmt.Println("Failed to refresh token:", err)
						redirectToLogin(c, c.Request.URL.String())
						return
					}

					// Update cookies with the new tokens
					c.SetCookie("access_token", newAccessToken, 3600, "/", apollobridgeDomain, true, true)
					c.SetCookie("refresh_token", newRefreshToken, 3600*24, "/", apollobridgeDomain, true, true)

					// Retry the request with the new access token
					c.Request.Header.Set("Authorization", "Bearer "+newAccessToken)
					c.Next()
					return
				}
			}

			// Other token errors
			fmt.Printf("Token parsing error: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if !hasRole(token, requiredRole, clientID) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		c.Next()
	}
}

func hasRole(token *jwt.Token, requiredRole string, clientID string) bool {
	claims, ok := token.Claims.(jwt.MapClaims)
	if requiredRole == "" {
		return true
	}
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

func getKeycloakPublicKey(token *jwt.Token) (interface{}, error) {
	// Extract the "kid" from the token header
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("missing kid in token header")
	}

	// Fetch the JWKS
	set, err := jwk.Fetch(context.Background(), jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWK set: %v", err)
	}

	// Find the key with the matching "kid"
	keys, _ := set.LookupKeyID(kid)
	if keys == nil {
		return nil, fmt.Errorf("no matching key found for kid: %s", kid)
	}

	// Extract the public key
	var pubKey interface{}
	if err := keys.Raw(&pubKey); err != nil {
		return nil, fmt.Errorf("failed to extract public key: %v", err)
	}

	return pubKey, nil
}
