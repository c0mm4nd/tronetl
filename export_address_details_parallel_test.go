package main

import (
	"bytes"
	"encoding/csv"
	"os"
	"testing"
	"time"
)

// runExport is a helper to execute export with given workers and addresses, returning elapsed time and CSV buffers.
func runExport(t *testing.T, workers uint, addrs []string) (time.Duration, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	accBuf := &bytes.Buffer{}
	ctrBuf := &bytes.Buffer{}
	provider := os.Getenv("TRON_PROVIDER")

	opts := &ExportAddressDetailsOptions{
		Addresses:       addrs,
		ProviderURI:     provider,
		accountsOutput:  accBuf,
		contractsOutput: ctrBuf,
		tokensOutput:    nil, // keep tokens disabled for speed
	}

	start := time.Now()
	ExportAddressDetailsWithWorkers(opts, workers)
	return time.Since(start), accBuf, ctrBuf
}

// Ensures parallel workers bring speed-up versus single worker when talking to a real node.
func TestExportAddressDetailsWithWorkersSpeed(t *testing.T) {
	// repeat same known contract to generate enough remote calls
	addrs := make([]string, 0, 40)
	for i := 0; i < 40; i++ {
		addrs = append(addrs, "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t") // USDT contract
	}

	t1, _, _ := runExport(t, 1, addrs)
	t4, _, _ := runExport(t, 4, addrs)

	t.Logf("elapsed single=%v, parallel=%v", t1, t4)

	if t4 >= t1 {
		t.Fatalf("expected 4 workers to be faster or equal, got single=%v parallel=%v", t1, t4)
	}
}

// Verifies real-node outputs have rows for accounts and contracts paths.
func TestExportAddressDetailsWithWorkersOutputs(t *testing.T) {
	addrs := []string{
		"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", // contract
		"THPvaUhoh2Qn2y9THCZML3H815hhFhn5YC", // regular account
	}

	_, accBuf, ctrBuf := runExport(t, 2, addrs)

	accRows, err := csv.NewReader(bytes.NewReader(accBuf.Bytes())).ReadAll()
	if err != nil {
		t.Fatalf("read accounts csv failed: %v", err)
	}
	ctrRows, err := csv.NewReader(bytes.NewReader(ctrBuf.Bytes())).ReadAll()
	if err != nil {
		t.Fatalf("read contracts csv failed: %v", err)
	}

	if len(accRows) != len(addrs)+1 { // header + rows
		t.Fatalf("accounts rows mismatch: want %d got %d", len(addrs)+1, len(accRows))
	}

	if len(ctrRows) < 2 { // header + at least the contract row
		t.Fatalf("expected contract rows, got %d", len(ctrRows))
	}

	// ensure the USDT contract row exists and is flagged as ERC20
	found := false
	for i := 1; i < len(ctrRows); i++ {
		if ctrRows[i][0] == "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" {
			found = true
			if ctrRows[i][3] != "true" {
				t.Fatalf("contract is_erc20 expected true, got %s", ctrRows[i][3])
			}
			break
		}
	}
	if !found {
		t.Fatalf("USDT contract row not found in contracts CSV")
	}
}
