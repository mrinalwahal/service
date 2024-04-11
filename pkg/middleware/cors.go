package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

type CORSConfig struct {

	// AllowedOrigins is the list of origins that are allowed to access the resource.
	// Default: `[]string{"*"}`
	//
	// This field is optional.
	AllowedOrigins []string

	// AllowedMethods is the list of methods that are allowed to access the resource.
	// Default: `[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}`
	//
	// This field is optional.
	AllowedMethods []string

	// AllowedHeaders is the list of headers that are allowed to access the resource.
	// Default: `[]string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization
	// "accept", "origin", "Cache-Control", "X-Requested-With"}`
	//
	// This field is optional.
	AllowedHeaders []string

	// AllowCredentials is the flag that determines if the resource allows credentials.
	// Default: `false`
	//
	// This field is optional.
	AllowCredentials bool
}

// CORS middleware adds the CORS headers to the response.
func CORS(config *CORSConfig) Middleware {

	// Set the default configuration.

	if config.AllowedOrigins == nil {
		config.AllowedOrigins = []string{"*"}
	}

	if config.AllowedMethods == nil {
		config.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}

	if config.AllowedHeaders == nil {
		config.AllowedHeaders = []string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"accept",
			"origin",
			"Cache-Control",
			"X-Requested-With",
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", strings.Join(config.AllowedOrigins, ","))
			w.Header().Add("Access-Control-Allow-Credentials", fmt.Sprint(config.AllowCredentials))
			w.Header().Add("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ","))
			w.Header().Add("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ","))

			if r.Method == http.MethodOptions {
				http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
