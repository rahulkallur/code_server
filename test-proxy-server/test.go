package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

func test() {
	target, _ := url.Parse("http://localhost:3000")

	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     100,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host

			// Preserve original Host header (code-server needs this)
			origHost := req.Host
			req.Header.Set("X-Forwarded-Host", origHost)
			req.Header.Set("X-Forwarded-Proto", "http")
			req.Header.Set("X-Proxy-Server", "Go Reverse Proxy")
			if clientIP := req.RemoteAddr; clientIP != "" {
				req.Header.Set("X-Forwarded-For", clientIP)
			}
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
		log.Printf("%s %s (Host: %s)", r.Method, r.URL.String(), r.Host)
		if isWebSocket(r) {
			proxyWebSocket(w, r, target)
			return
		}
		proxy.ServeHTTP(w, r)
	})

	log.Println("Starting proxy server on :8080, forwarding to http://localhost:3000")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Proxy server failed: %v", err)
	}
}

func isWebSocket(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("Connection")) == "upgrade" &&
		strings.ToLower(r.Header.Get("Upgrade")) == "websocket"
}

func proxyWebSocket(w http.ResponseWriter, r *http.Request, target *url.URL) {
	log.Printf("WebSocket upgrade: %s (Host: %s)", r.URL.String(), r.Host)

	backendConn, err := net.DialTimeout("tcp", target.Host, 10*time.Second)
	if err != nil {
		log.Printf("Backend dial error: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer backendConn.Close()

	// Forward the original Host header from the client
	// (code-server requires this, per docs)
	origHost := r.Host
	r.URL.Host = target.Host
	r.URL.Scheme = target.Scheme
	r.Header.Set("X-Forwarded-Host", origHost)
	r.Header.Set("X-Forwarded-Proto", "http")

	if err := r.Write(backendConn); err != nil {
		log.Printf("Error writing request to backend: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		log.Printf("Hijack error: %v", err)
		return
	}
	defer clientConn.Close()

	// Read the backend's 101 response first
	br := bufio.NewReader(backendConn)
	resp, err := http.ReadResponse(br, r)
	if err != nil {
		log.Printf("Error reading backend response: %v", err)
		return
	}
	log.Printf("Backend response: %s", resp.Status)

	// Write the 101 response to the client
	if err := resp.Write(clientConn); err != nil {
		log.Printf("Error writing response to client: %v", err)
		return
	}

	// Now tunnel bidirectional data
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		n, err := io.Copy(clientConn, br)
		log.Printf("Backend->Client tunnel closed: %d bytes, err=%v", n, err)
		clientConn.Close()
		backendConn.Close()
	}()

	go func() {
		defer wg.Done()
		n, err := io.Copy(backendConn, clientConn)
		log.Printf("Client->Backend tunnel closed: %d bytes, err=%v", n, err)
		backendConn.Close()
		clientConn.Close()
	}()

	wg.Wait()
	log.Println("WebSocket tunnel closed")
}
