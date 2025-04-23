package api

import (
	"log"
	"net/http"

	"Ethereum-fund-flow-analysis/internal/config"
)

// SetupRouter sets up the HTTP router with all routes
func SetupRouter(cfg *config.Config) http.Handler {
	// Create a new handler with required services
	handler := NewHandler(cfg)

	// Create a new router
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/beneficiary", handler.BeneficiaryHandler)
	mux.HandleFunc("/payer", handler.PayerHandler)

	// Add middleware for logging, CORS, etc.
	return LoggingMiddleware(mux)
}

// LoggingMiddleware logs all incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		// This is a simple logger - in a production environment,
		// you might want to use a more sophisticated logging package
		// like Zap or Logrus
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
