package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/olegshulyakov/go-briefly-bot/pkg/db"
	"github.com/olegshulyakov/go-briefly-bot/pkg/utils"
)

// Context keys for storing validated data
type contextKey string

const (
	// ContextClientAppName is the key used to store the validated client app name in the request context.
	ContextClientAppName contextKey = "clientAppName"
	// ContextClientAppID is the key used to store the validated client app ID in the request context.
	ContextClientAppID contextKey = "clientAppID"
)

// AuthMiddleware creates a middleware function that validates the API key.
// It expects the configuration to be available, likely through a closure or global access.
// For this implementation, we'll assume config is accessible.
// A more robust approach might be to pass required dependencies (like the config or DB manager) to the middleware factory.
func AuthMiddleware(config *utils.Config, logger *utils.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the API key from the Authorization header (Bearer token style)
			// e.g., "Authorization: Bearer YOUR_API_KEY"
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Authentication failed: Missing Authorization header", "path", r.URL.Path, "method", r.Method)
				http.Error(w, `{"error": "Missing Authorization header. Expected 'Authorization: Bearer <API_KEY>'"}`, http.StatusUnauthorized)
				return
			}

			// Check for "Bearer " prefix
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				logger.Warn("Authentication failed: Invalid Authorization header format", "path", r.URL.Path, "method", r.Method)
				http.Error(w, `{"error": "Invalid Authorization header format. Expected 'Bearer <API_KEY>'"}`, http.StatusUnauthorized)
				return
			}

			// Extract the API key
			apiKey := strings.TrimPrefix(authHeader, bearerPrefix)

			// Validate the API key against the configured one
			// In a more complex system, you might validate against a database or external service.
			if apiKey != config.APIKey {
				logger.Warn("Authentication failed: Invalid API key provided", "path", r.URL.Path, "method", r.Method)
				http.Error(w, `{"error": "Invalid API key"}`, http.StatusUnauthorized)
				return
			}

			// If valid, call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func CORSMiddleware(config *utils.Config, logger *utils.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS middleware
			next.ServeHTTP(w, r)
		})
	}
}

func LoggingMiddleware(config *utils.Config, logger *utils.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Request logging middleware
			next.ServeHTTP(w, r)
		})
	}
}

// ClientAppValidationMiddleware creates a middleware that validates the client app name
// provided in the URL path parameter against the DictClientApps table.
func ClientAppValidationMiddleware(dbManager *db.DBManager, logger *utils.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract client app name from the URL path.
			// This part is highly dependent on your router.
			// Assuming `gorilla/mux` style path variables:
			// vars := mux.Vars(r)
			// clientAppName := vars["client_app"]
			//
			// For standard library, you might need to parse it from r.URL.Path
			// This is a simplified example. You'd need to adapt based on your routing library.
			//
			// Let's assume a helper function or a way to get the path parameter.
			// For now, we'll use a placeholder. You need to replace this with actual path extraction logic.
			// Example with gorilla/mux:
			// clientAppName := mux.Vars(r)["client_app"]
			//
			// Since we don't have the router import here, let's simulate getting the path param.
			// You will need to adjust this part based on your actual routing implementation.
			// --- Simulate path parameter extraction ---
			// This is a placeholder. Replace with actual logic based on your router (e.g., Gorilla Mux).
			clientAppName := getClientAppFromPath(r) // You need to implement this function
			// --- End Simulation ---

			if clientAppName == "" {
				// If client app name is not in the path or is empty, it's a bad request for routes that require it.
				// This check might be redundant if your router ensures the parameter exists for specific routes.
				// But it's good practice to check.
				logger.Warn("Client App Validation failed: Missing client_app parameter in path", "path", r.URL.Path, "method", r.Method)
				http.Error(w, `{"error": "Missing client_app parameter in path"}`, http.StatusBadRequest)
				return
			}

			// Get the shared database connection to validate the client app
			sharedDB, err := dbManager.GetPrimaryDB()
			if err != nil {
				logger.Error("Client App Validation failed: Could not connect to shared database", "error", err)
				http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
				return
			}

			// Validate the client app name exists in the DictClientApps table
			// We can use the existing GetClientAppID function and check if it returns an error.
			// If it returns an error like "not found", the app is invalid.
			clientAppID, err := db.GetClientAppID(sharedDB, clientAppName)
			if err != nil {
				// Check if it's a "not found" error
				if strings.Contains(strings.ToLower(err.Error()), "not found") {
					logger.Warn("Client App Validation failed: Invalid client app name", "client_app", clientAppName, "path", r.URL.Path, "method", r.Method)
					http.Error(w, fmt.Sprintf(`{"error": "Invalid client app name: %s"}`, clientAppName), http.StatusBadRequest)
					return
				}
				// Any other error is likely an internal DB issue
				logger.Error("Client App Validation failed: Error querying DictClientApps", "error", err, "client_app", clientAppName)
				http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
				return
			}

			// If valid, add the client app name and ID to the request context for downstream handlers
			ctx := context.WithValue(r.Context(), ContextClientAppName, clientAppName)
			ctx = context.WithValue(ctx, ContextClientAppID, clientAppID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// getClientAppFromPath is a placeholder function to demonstrate how you might extract
// the client app name from the request path.
// You MUST replace this with actual logic based on your chosen router.
// Example using Gorilla Mux:
// import "github.com/gorilla/mux"
//
//	func getClientAppFromPath(r *http.Request) string {
//	    vars := mux.Vars(r)
//	    return vars["client_app"] // Assumes the path parameter is named "client_app"
//	}
//
// Example using standard library path parsing (less robust):
// import "strings"
//
//	func getClientAppFromPath(r *http.Request) string {
//	    pathParts := strings.Split(r.URL.Path, "/")
//	    // Assuming path is like /v1/message/new/telegram or /v1/results/telegram
//	    // Find the index of "new" or "results" and take the next part
//	    for i, part := range pathParts {
//	        if (part == "new" || part == "results") && i+1 < len(pathParts) {
//	            return pathParts[i+1]
//	        }
//	    }
//	    return "" // Not found
//	}
//
// For this implementation, we'll provide a basic version assuming standard library routing
// and specific path patterns. You will need to adapt this.
func getClientAppFromPath(r *http.Request) string {
	// --- Basic Path Parsing Logic (Replace with your router's method) ---
	// This is a simplified example and might not cover all edge cases.
	// It assumes paths like:
	// POST /v1/message/new/{client_app}
	// GET /v1/results/{client_app}
	path := r.URL.Path
	// Normalize path by removing trailing slash
	path = strings.TrimSuffix(path, "/")

	parts := strings.Split(path, "/")
	// Expected structure: ["", "v1", "message", "new", "client_app_name"] or ["", "v1", "results", "client_app_name"]
	// We need to find the position of "new" or "results" and take the next part.

	for i, part := range parts {
		if part == "new" && i+1 < len(parts) {
			return parts[i+1]
		}
		if part == "results" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	// If not found in expected positions, return empty string
	// This might happen for other endpoints like /v1/status which don't require client_app validation
	return ""
}

// GetClientAppNameFromContext retrieves the validated client app name from the request context.
// This should be used by handlers that require client app validation.
func GetClientAppNameFromContext(ctx context.Context) (string, bool) {
	clientAppName, ok := ctx.Value(ContextClientAppName).(string)
	return clientAppName, ok
}

// GetClientAppIDFromContext retrieves the validated client app ID from the request context.
// This should be used by handlers that require the numeric client app ID.
func GetClientAppIDFromContext(ctx context.Context) (int8, bool) {
	clientAppID, ok := ctx.Value(ContextClientAppID).(int8)
	return clientAppID, ok
}
