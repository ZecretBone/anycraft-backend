package middleware

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// responseWriter wraps http.ResponseWriter so we can capture status and bytes.
type responseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		// default status if Write called without WriteHeader
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

type LoggerOptions struct {
	// Where to write logs. If nil, logs go to the default logger (stdout).
	Writer io.Writer
	// If true, logs request/response headers (no bodies).
	Verbose bool
}

// RequestLogger returns a middleware that logs every request.
func RequestLogger(opts *LoggerOptions) func(http.Handler) http.Handler {
	if opts == nil {
		opts = &LoggerOptions{}
	}
	// If a custom writer is provided, wire it into the std logger.
	if opts.Writer != nil {
		log.SetOutput(opts.Writer)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Basic request identity
			method := r.Method
			path := r.URL.Path
			ua := r.UserAgent()

			// Try to get a request id from upstream/proxy, or use short time-based fallback
			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = time.Now().Format("150405.000") // e.g., 165312.123
			}

			// Attempt to get real client IP (behind proxy)
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.RemoteAddr
			} else {
				// X-Forwarded-For can be "client, proxy1, proxy2"
				ip = strings.TrimSpace(strings.Split(ip, ",")[0])
			}

			// Wrap writer to capture status/size
			lw := &responseWriter{ResponseWriter: w}

			// Call next
			next.ServeHTTP(lw, r)

			dur := time.Since(start)

			// Access log line
			log.Printf(
				`[%s] %s %s -> %d (%dB) %s ip=%s ua="%s"`,
				reqID, method, path, lw.status, lw.bytes, dur, ip, ua,
			)

			// Optional: verbose headers (no bodies)
			if opts.Verbose {
				for k, v := range r.Header {
					log.Printf("[%s] req.hdr %s=%v", reqID, k, v)
				}
				for k, v := range lw.Header() {
					log.Printf("[%s] res.hdr %s=%v", reqID, k, v)
				}
			}
		})
	}
}
