package main

import (
	"strings"
	"testing"
)

func TestCollectAllAddrs(t *testing.T) {
	// two valid T-addrs (34 chars), plus an invalid entry that should be ignored
	src := strings.NewReader("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t,invalid,THPvaUhoh2Qn2y9THCZML3H815hhFhn5YC\n")

	options := &ExportAddressDetailsOptions{
		addrSource: src,
		Addresses:  []string{"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"}, // duplicate on purpose
	}

	all := collectAllAddrs(options)

	if len(all) != 2 { // two unique addresses (duplicate is removed)
		t.Fatalf("expected 2 addresses, got %d", len(all))
	}

	// Check that all addresses are valid T-addresses
	expectedAddrs := map[string]bool{
		"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t": false,
		"THPvaUhoh2Qn2y9THCZML3H815hhFhn5YC": false,
	}

	for _, addr := range all {
		if len(addr) != 34 || addr[0] != 'T' {
			t.Fatalf("addr %s is not a valid T-address", addr)
		}
		if _, ok := expectedAddrs[addr]; !ok {
			t.Fatalf("unexpected address: %s", addr)
		}
		expectedAddrs[addr] = true
	}

	// Ensure both expected addresses were found
	for addr, found := range expectedAddrs {
		if !found {
			t.Fatalf("expected address %s not found", addr)
		}
	}
}

func TestCollectAllAddrs_HexDeduplication(t *testing.T) {
	// Test that hex addresses are properly converted to T-addr and deduplicated
	// 41a614f803b6fd780986a42c78ec9c7f77e6ded13c is the hex for TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t
	options := &ExportAddressDetailsOptions{
		Addresses: []string{
			"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",                 // T-addr
			"41a614f803b6fd780986a42c78ec9c7f77e6ded13c",         // Hex (same address)
			"a614f803b6fd780986a42c78ec9c7f77e6ded13c",           // Hex without 41 prefix (same address)
		},
	}

	all := collectAllAddrs(options)

	if len(all) != 1 { // should be deduplicated to 1 address
		t.Fatalf("expected 1 address after deduplication, got %d", len(all))
	}

	if all[0] != "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" {
		t.Fatalf("expected TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t, got %s", all[0])
	}
}
