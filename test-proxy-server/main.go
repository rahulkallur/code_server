package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func main() {
	target, _ := url.Parse("http://localhost:3000")

	// Configure a transport with timeouts and connection pooling
	transport := &http.Transport{
		// Connection pooling settings
		MaxIdleConns:        100,              // Total idle connections across all hosts
		MaxIdleConnsPerHost: 10,               // Idle connections per backend host
		MaxConnsPerHost:     100,              // Max concurrent connections per host
		IdleConnTimeout:     90 * time.Second, // 90 seconds

		// Timeouts for connection establishment
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,

		// TLS handshake timeout
		TLSHandshakeTimeout: 10 * time.Second,

		// Response header timeout
		ResponseHeaderTimeout: 10 * time.Second,
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// Set the target Host and Scheme
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host

			// Add a custom header to identify the proxied request

			req.Header.Set("X-Proxy-Server", "Go Reverse Proxy")
			// req.Header.Set("X-Forwarded-Host", req.Host)

			// Add the original clientIP for logging or debugging purposes
			if clientIP := req.RemoteAddr; clientIP != "" {
				req.Header.Set("X-Forwarded-For", clientIP)
			}

			// Log the request
			log.Printf("Proxying request: %s %s", req.Method, req.URL.String())
		},
		ModifyResponse: func(resp *http.Response) error {
			resp.Header.Set("X-Content-Type-Options", "nosniff")
			resp.Header.Set("X-XSS-Protection", "1; mode=block")
			resp.Header.Set("Proxy-Server", "Go Reverse Proxy")
			return nil
		},
		Transport: transport,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s (Host: %s) %s %s", r.Method, r.URL.String(), r.Host, r.Header.Get("Connection"), r.Header.Get("Upgrade"))
		// if isWebSocketConnection(r) {

		// }
		proxy.ServeHTTP(w, r)
	})

	log.Println("Starting proxy server on :8080, forwarding to http://localhost:5500")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Proxy server failed: %v", err)
	}

}

func isWebSocketConnection(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("Connection")) == "upgrade" &&
		strings.ToLower(r.Header.Get("Upgrade")) == "websocket"
}
