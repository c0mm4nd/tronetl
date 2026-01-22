package tron

import (
	"testing"
)

func TestNewTronClientWithOverrides(t *testing.T) {
	tests := []struct {
		name               string
		providerURL        string
		httpProviderURI    string
		jsonrpcProviderURI string
		expectedHTTPURI    string
		expectedJSONURI    string
	}{
		{
			name:               "No overrides",
			providerURL:        "http://example.com",
			httpProviderURI:    "",
			jsonrpcProviderURI: "",
			expectedHTTPURI:    "http://example.com",
			expectedJSONURI:    "http://example.com/jsonrpc",
		},
		{
			name:               "HTTP override only",
			providerURL:        "http://example.com",
			httpProviderURI:    "http://custom-http.com:8090",
			jsonrpcProviderURI: "",
			expectedHTTPURI:    "http://custom-http.com:8090",
			expectedJSONURI:    "http://example.com/jsonrpc",
		},
		{
			name:               "JSONRPC override only",
			providerURL:        "http://example.com",
			httpProviderURI:    "",
			jsonrpcProviderURI: "http://custom-jsonrpc.com:50545",
			expectedHTTPURI:    "http://example.com",
			expectedJSONURI:    "http://custom-jsonrpc.com:50545",
		},
		{
			name:               "Both overrides",
			providerURL:        "http://example.com",
			httpProviderURI:    "http://custom-http.com:8090",
			jsonrpcProviderURI: "http://custom-jsonrpc.com:50545",
			expectedHTTPURI:    "http://custom-http.com:8090",
			expectedJSONURI:    "http://custom-jsonrpc.com:50545",
		},
		{
			name:               "Overrides with localhost default",
			providerURL:        "localhost",
			httpProviderURI:    "http://custom-http.com:8090",
			jsonrpcProviderURI: "http://custom-jsonrpc.com:50545",
			expectedHTTPURI:    "http://custom-http.com:8090",
			expectedJSONURI:    "http://custom-jsonrpc.com:50545",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewTronClientWithOverrides(tt.providerURL, tt.httpProviderURI, tt.jsonrpcProviderURI)
			
			if client.httpURI != tt.expectedHTTPURI {
				t.Errorf("httpURI = %v, want %v", client.httpURI, tt.expectedHTTPURI)
			}
			
			if client.jsonURI != tt.expectedJSONURI {
				t.Errorf("jsonURI = %v, want %v", client.jsonURI, tt.expectedJSONURI)
			}
		})
	}
}

func TestNewTronClientBackwardCompatibility(t *testing.T) {
	// Test that NewTronClient still works as before
	client := NewTronClient("http://example.com")
	
	if client.httpURI != "http://example.com" {
		t.Errorf("httpURI = %v, want %v", client.httpURI, "http://example.com")
	}
	
	if client.jsonURI != "http://example.com/jsonrpc" {
		t.Errorf("jsonURI = %v, want %v", client.jsonURI, "http://example.com/jsonrpc")
	}
}
