package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not found"})
		return
	}

	originalURL := c.Query("state")

	kcClientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	tokenURL := "http://keycloak:8080/realms/apollo/protocol/openid-connect/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", "http://localhost:3000/auth/callback")
	data.Set("client_id", "apollo-client")
	data.Set("client_secret", kcClientSecret)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get token", "details": err.Error()})
		return
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Failed to get token",
			"status":   resp.StatusCode,
			"response": string(bodyBytes),
		})
		return
	}

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode token response"})
		return
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token response"})
		return
	}

	c.SetCookie("access_token", accessToken, 3600, "/", "localhost", false, true)

	// Redirect to the original URL
	if originalURL != "" {
		c.Redirect(http.StatusFound, originalURL)
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Token stored in cookie"})
	}
}
