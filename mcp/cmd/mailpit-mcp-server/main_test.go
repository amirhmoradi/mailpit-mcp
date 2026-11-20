package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/amirhmoradi/mailpit-mcp/internal/client"
	"github.com/amirhmoradi/mailpit-mcp/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestStreamableHTTP(t *testing.T) {
	mailpit := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/message/fixture" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ID":"fixture","Subject":"SMTP fixture","Text":"Captured test message"}`)
	}))
	defer mailpit.Close()
	c := client.New(client.Config{BaseURL: mailpit.URL, Timeout: time.Second})
	httpServer := httptest.NewServer(httpHandler(server.New(c)))
	defer httpServer.Close()

	for _, version := range []string{"2026-07-28", "2025-11-25"} {
		t.Run(version, func(t *testing.T) {
			ctx := t.Context()
			cl := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
			session, err := cl.Connect(ctx, &mcp.StreamableClientTransport{
				Endpoint: httpServer.URL + "/mcp",
			}, &mcp.ClientSessionOptions{ProtocolVersion: version})
			if err != nil {
				t.Fatal(err)
			}
			defer session.Close()
			if got := session.InitializeResult().ProtocolVersion; got != version {
				t.Fatalf("negotiated %s, want %s", got, version)
			}
			list, err := session.ListTools(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(list.Tools) != 24 {
				t.Fatalf("got %d tools, want 24", len(list.Tools))
			}
			result, err := session.CallTool(ctx, &mcp.CallToolParams{
				Name: "get_message", Arguments: map[string]any{"id": "fixture"},
			})
			if err != nil {
				t.Fatal(err)
			}
			if result.IsError || len(result.Content) != 1 {
				t.Fatalf("unexpected tool result: %+v", result)
			}
			text, ok := result.Content[0].(*mcp.TextContent)
			if !ok || !strings.Contains(text.Text, "Captured test message") {
				t.Fatalf("missing fixture content: %+v", result.Content[0])
			}
		})
	}
}

func TestMCPHTTPRejectsUntrustedBrowserRequests(t *testing.T) {
	handler := httpHandler(server.New(client.New(client.Config{})))
	for _, tc := range []struct{ host, origin string }{
		{host: "untrusted.example:3000"},
		{host: "127.0.0.1:3000", origin: "https://untrusted.example"},
	} {
		req := httptest.NewRequest(http.MethodPost, "http://"+tc.host+"/mcp", nil)
		if tc.origin != "" {
			req.Header.Set("Origin", tc.origin)
		}
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		if resp.Code != http.StatusForbidden {
			t.Errorf("host %q origin %q: status %d, want 403", tc.host, tc.origin, resp.Code)
		}
	}
}
