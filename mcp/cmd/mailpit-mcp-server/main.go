// Package main provides the entry point for the Mailpit MCP server.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/amirhmoradi/mailpit-mcp/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Load configuration from environment
	cfg := server.DefaultConfig()
	cfg.LoadFromEnv()

	// Create Mailpit API client
	mailpitClient := cfg.NewMailpitClient()

	// Create MCP server
	mcpServer := server.New(mailpitClient)

	// Run with the appropriate transport
	switch cfg.Transport {
	case "stdio":
		runSTDIO(mcpServer)
	case "http":
		runHTTP(mcpServer, cfg)
	default:
		log.Fatalf("Unknown transport: %s (use 'stdio' or 'http')", cfg.Transport)
	}
}

// runSTDIO runs the MCP server over standard input/output.
func runSTDIO(s *mcp.Server) {
	// Log to stderr since stdout is used for MCP communication
	log.SetOutput(os.Stderr)
	log.Printf("Starting Mailpit MCP server (STDIO transport, version %s)", server.Version)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Shutting down...")
		cancel()
	}()

	if err := s.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func httpHandler(s *mcp.Server, allowedHosts []string) http.Handler {
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return s
	}, &mcp.StreamableHTTPOptions{
		Stateless:                  true,
		JSONResponse:               true,
		MaxRequestBodyBytes:        64 << 20,
		DisableLocalhostProtection: true, // replaced by the explicit Host check below
	})
	protected := http.NewCrossOriginProtection().Handler(handler)

	mux := http.NewServeMux()
	mux.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		// Docker can deliver published-port traffic via a bridge address.
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
		}
		host = strings.Trim(host, "[]")
		if strings.EqualFold(host, "localhost") || net.ParseIP(host).IsLoopback() {
			protected.ServeHTTP(w, r)
			return
		}
		for _, allowed := range allowedHosts {
			allowed = strings.TrimSpace(allowed)
			if allowed != "" && strings.EqualFold(host, allowed) {
				protected.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Forbidden: unrecognized Host", http.StatusForbidden)
	})
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})
	return mux
}

// runHTTP runs the server over stateless Streamable HTTP.
func runHTTP(s *mcp.Server, cfg *server.Config) {
	addr := net.JoinHostPort(cfg.HTTPHost, strconv.Itoa(cfg.HTTPPort))
	log.Printf("Starting Mailpit MCP server (Streamable HTTP, version %s) on %s", server.Version, addr)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           httpHandler(s, cfg.HTTPAllowedHosts),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			httpServer.Close()
		}
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
	<-done
}
