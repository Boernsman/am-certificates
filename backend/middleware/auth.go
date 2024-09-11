package middleware

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gopkg.in/ini.v1"
)

// Global variables for storing credentials and API keys
var (
	validUsername string
	validPassword string
	apiKeys       map[string]bool
)

// Load the credentials and API keys from the INI file
func LoadCredentials(file string) error {
	cfg, err := ini.Load(file)
	if err != nil {
		return fmt.Errorf("failed to load config.ini: %v", err)
	}

	// Load basic auth credentials
	authSection := cfg.Section("auth")
	validUsername = authSection.Key("username").String()
	validPassword = authSection.Key("password").String()

	// Load API keys
	apiKeys = make(map[string]bool)
	apiSection := cfg.Section("api_keys")
	for _, key := range apiSection.Keys() {
		apiKeys[key.String()] = true
	}

	return nil
}

// Middleware for basic authentication
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Panicln("Authentication header is missing")
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		// Validate the basic auth header
		if !validateBasicAuth(authHeader) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		log.Panicln("Basic authentication successfull")
		next.ServeHTTP(w, r)
	})
}

// Middleware for API key validation
func ApiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" || !validateAPIKey(apiKey) {
			http.Error(w, "Invalid or missing API key", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Validate the basic auth credentials
func validateBasicAuth(authHeader string) bool {
	const basicPrefix = "Basic "
	if !strings.HasPrefix(authHeader, basicPrefix) {
		log.Println("Authentication header is missing")
		return false
	}

	// Decode the base64-encoded username:password string
	encodedCreds := strings.TrimPrefix(authHeader, basicPrefix)
	decodedCreds, err := base64.StdEncoding.DecodeString(encodedCreds)
	if err != nil {
		log.Println("Could not decode basic auth credentials")
		return false
	}

	// Split the credentials into username and password
	creds := strings.SplitN(string(decodedCreds), ":", 2)
	if len(creds) != 2 {
		log.Println("Could not split basic auth credentials into username and password")
		return false
	}

	// Compare with the valid username and password
	return creds[0] == validUsername && creds[1] == validPassword
}

// Validate the API key
func validateAPIKey(apiKey string) bool {
	_, exists := apiKeys[apiKey]
	return exists
}
