package main

import (
	"strings"
	"testing"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
)

func TestCollectAllAddrs(t *testing.T) {
	// two valid T-addrs (34 chars), plus an invalid entry that should be ignored
	src := strings.NewReader("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t,invalid,THPvaUhoh2Qn2y9THCZML3H815hhFhn5YC\n")

	options := &ExportAddressDetailsOptions{
		addrSource: src,
		Addresses:  []string{"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"}, // duplicate on purpose
	}

	all := collectAllAddrs(options)

	if len(all) != 3 { // two from source + one duplicate from Addresses slice
		t.Fatalf("expected 3 addresses, got %d", len(all))
	}

	expectedHex := tron.EnsureHexAddr("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
	if all[0] != expectedHex {
		t.Fatalf("first addr unexpected, want %s got %s", expectedHex, all[0])
	}

	for _, addr := range all {
		if len(addr) != 42 || !strings.HasPrefix(addr, "41") {
			t.Fatalf("addr %s is not a 21-byte hex Tron address", addr)
		}
	}
}
