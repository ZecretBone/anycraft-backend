package middleware

import (
	"net/http"
	"os"
	"strings"
)

type CORSOptions struct {
	// Comma-separated list of allowed origins (e.g. "http://localhost:3000,https://your.site")
	AllowOrigins string
	// If true, sets Access-Control-Allow-Credentials: true
	AllowCredentials bool
	// Allowed headers for preflight
	AllowHeaders string // e.g. "Content-Type,Authorization,X-Requested-With"
	// Allowed methods for preflight
	AllowMethods string // e.g. "GET,POST,OPTIONS"
}

func getAllowedOrigin(reqOrigin, allowList string) string {
	if reqOrigin == "" {
		return ""
	}
	if allowList == "" || allowList == "*" {
		return "*"
	}
	for _, o := range strings.Split(allowList, ",") {
		if strings.TrimSpace(o) == reqOrigin {
			return reqOrigin
		}
	}
	return ""
}

// CORS wraps next with CORS headers and handles OPTIONS preflight.
func CORS(opts *CORSOptions) func(http.Handler) http.Handler {
	if opts == nil {
		opts = &CORSOptions{}
	}
	if opts.AllowOrigins == "" {
		// Read from env or fallback to dev defaults
		opts.AllowOrigins = os.Getenv("CORS_ORIGINS")
		if opts.AllowOrigins == "" {
			opts.AllowOrigins = "http://localhost:3000"
		}
	}
	if opts.AllowHeaders == "" {
		opts.AllowHeaders = "Content-Type,Authorization,X-Requested-With"
	}
	if opts.AllowMethods == "" {
		opts.AllowMethods = "GET,POST,OPTIONS"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowed := getAllowedOrigin(origin, opts.AllowOrigins)

			// Always advertise what we allow (some browsers need these on every response)
			if allowed != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowed)
				if opts.AllowCredentials && allowed != "*" {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}
			w.Header().Set("Vary", "Origin") // so proxies don't mix origins
			w.Header().Set("Access-Control-Allow-Methods", opts.AllowMethods)
			w.Header().Set("Access-Control-Allow-Headers", opts.AllowHeaders)

			// Preflight
			if r.Method == http.MethodOptions {
				// If no matching origin, still end preflight to avoid 404
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
