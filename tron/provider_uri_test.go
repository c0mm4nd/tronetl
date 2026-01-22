package tron

import (
	"testing"
)

func TestNewTronClientURLHandling(t *testing.T) {
	tests := []struct {
		name            string
		providerURL     string
		expectedHTTPURI string
		expectedJSONURI string
		description     string
	}{
		{
			name:            "empty string defaults to localhost",
			providerURL:     "",
			expectedHTTPURI: "http://localhost:8090",
			expectedJSONURI: "http://localhost:50545/jsonrpc",
			description:     "Empty provider should default to localhost with ports",
		},
		{
			name:            "http with IP but no port should add default port 8090",
			providerURL:     "http://192.168.8.14",
			expectedHTTPURI: "http://192.168.8.14:8090",
			expectedJSONURI: "http://192.168.8.14:50545/jsonrpc",
			description:     "http://IP without port should add :8090 for HTTP and :50545 for JSON-RPC",
		},
		{
			name:            "https with IP but no port should add default port 8090",
			providerURL:     "https://192.168.8.14",
			expectedHTTPURI: "https://192.168.8.14:8090",
			expectedJSONURI: "https://192.168.8.14:50545/jsonrpc",
			description:     "https://IP without port should add :8090 for HTTP and :50545 for JSON-RPC",
		},
		{
			name:            "http with IP and port should use as-is",
			providerURL:     "http://192.168.8.14:8090",
			expectedHTTPURI: "http://192.168.8.14:8090",
			expectedJSONURI: "http://192.168.8.14:8090/jsonrpc",
			description:     "http://IP:port should be used as-is",
		},
		{
			name:            "https with IP and port should use as-is",
			providerURL:     "https://192.168.8.14:8090",
			expectedHTTPURI: "https://192.168.8.14:8090",
			expectedJSONURI: "https://192.168.8.14:8090/jsonrpc",
			description:     "https://IP:port should be used as-is",
		},
		{
			name:            "IP with port but no protocol should add http://",
			providerURL:     "192.168.8.14:8090",
			expectedHTTPURI: "http://192.168.8.14:8090",
			expectedJSONURI: "http://192.168.8.14:8090/jsonrpc",
			description:     "IP:port without protocol should get http:// prefix",
		},
		{
			name:            "hostname without port should add default port",
			providerURL:     "localhost",
			expectedHTTPURI: "localhost:8090",
			expectedJSONURI: "localhost:50545/jsonrpc",
			description:     "hostname without port should add default ports",
		},
		{
			name:            "hostname with port should add http://",
			providerURL:     "localhost:8090",
			expectedHTTPURI: "http://localhost:8090",
			expectedJSONURI: "http://localhost:8090/jsonrpc",
			description:     "hostname:port without protocol should get http:// prefix",
		},
		{
			name:            "http with hostname but no port should add default port",
			providerURL:     "http://localhost",
			expectedHTTPURI: "http://localhost:8090",
			expectedJSONURI: "http://localhost:50545/jsonrpc",
			description:     "http://hostname without port should add :8090 for HTTP and :50545 for JSON-RPC",
		},
		{
			name:            "http with hostname and port should use as-is",
			providerURL:     "http://localhost:8090",
			expectedHTTPURI: "http://localhost:8090",
			expectedJSONURI: "http://localhost:8090/jsonrpc",
			description:     "http://hostname:port should be used as-is",
		},
		{
			name:            "http with domain but no port should add default port",
			providerURL:     "http://tron.example.com",
			expectedHTTPURI: "http://tron.example.com:8090",
			expectedJSONURI: "http://tron.example.com:50545/jsonrpc",
			description:     "http://domain without port should add :8090 for HTTP and :50545 for JSON-RPC",
		},
		{
			name:            "http with domain and port should use as-is",
			providerURL:     "http://tron.example.com:8090",
			expectedHTTPURI: "http://tron.example.com:8090",
			expectedJSONURI: "http://tron.example.com:8090/jsonrpc",
			description:     "http://domain:port should be used as-is",
		},
		{
			name:            "http with IPv6 address without port should add default port",
			providerURL:     "http://[::1]",
			expectedHTTPURI: "http://[::1]:8090",
			expectedJSONURI: "http://[::1]:50545/jsonrpc",
			description:     "http://[IPv6] without port should add :8090 for HTTP and :50545 for JSON-RPC",
		},
		{
			name:            "http with IPv6 address and port should use as-is",
			providerURL:     "http://[::1]:8090",
			expectedHTTPURI: "http://[::1]:8090",
			expectedJSONURI: "http://[::1]:8090/jsonrpc",
			description:     "http://[IPv6]:port should be used as-is",
		},
		{
			name:            "http with IP and path but no port should add port before path",
			providerURL:     "http://192.168.8.14/api/v1",
			expectedHTTPURI: "http://192.168.8.14:8090/api/v1",
			expectedJSONURI: "http://192.168.8.14:50545/api/v1/jsonrpc",
			description:     "http://IP/path without port should insert port before path",
		},
		{
			name:            "http with hostname and path but no port should add port before path",
			providerURL:     "http://example.com/api/v1",
			expectedHTTPURI: "http://example.com:8090/api/v1",
			expectedJSONURI: "http://example.com:50545/api/v1/jsonrpc",
			description:     "http://hostname/path without port should insert port before path",
		},
		{
			name:            "http with IP, port and path should use as-is",
			providerURL:     "http://192.168.8.14:8090/api/v1",
			expectedHTTPURI: "http://192.168.8.14:8090/api/v1",
			expectedJSONURI: "http://192.168.8.14:8090/api/v1/jsonrpc",
			description:     "http://IP:port/path should be used as-is",
		},
		{
			name:            "IPv6 without protocol and without port should add http and port",
			providerURL:     "[::1]",
			expectedHTTPURI: "http://[::1]:8090",
			expectedJSONURI: "http://[::1]:50545/jsonrpc",
			description:     "[IPv6] without protocol should add http:// and ports",
		},
		{
			name:            "IPv6 without protocol but with port should add http",
			providerURL:     "[::1]:8090",
			expectedHTTPURI: "http://[::1]:8090",
			expectedJSONURI: "http://[::1]:8090/jsonrpc",
			description:     "[IPv6]:port without protocol should add http://",
		},
		{
			name:            "IPv6 full address without protocol and without port",
			providerURL:     "[2001:db8::1]",
			expectedHTTPURI: "http://[2001:db8::1]:8090",
			expectedJSONURI: "http://[2001:db8::1]:50545/jsonrpc",
			description:     "[IPv6] full address without protocol should add http:// and ports",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewTronClient(tt.providerURL)

			if client.httpURI != tt.expectedHTTPURI {
				t.Errorf("httpURI mismatch for %s:\n  got:      %q\n  expected: %q\n  %s",
					tt.name, client.httpURI, tt.expectedHTTPURI, tt.description)
			}

			if client.jsonURI != tt.expectedJSONURI {
				t.Errorf("jsonURI mismatch for %s:\n  got:      %q\n  expected: %q\n  %s",
					tt.name, client.jsonURI, tt.expectedJSONURI, tt.description)
			}
		})
	}
}
