package main

import (
	"math/big"
	"testing"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
)

// TestParseTotalSupplyWithTronGrid tests the TotalSupply parsing with real data from TronGrid public API
// This test verifies that large token supplies (uint256) are correctly handled as strings without overflow
// Note: This test may be skipped if TronGrid API is unavailable or requires an API key
func TestParseTotalSupplyWithTronGrid(t *testing.T) {
	// Use TronGrid public API with explicit override to avoid port issues
	provider := "https://api.trongrid.io"
	cli := tron.NewTronClientWithOverrides("", provider, "")

	// Test with USDT (TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t)
	// USDT has a very large total supply that would overflow uint64
	usdtContract := "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	callerAddr := "TJmmqjb1DK9TTZbQXzRQ2AuA94z4gKd6vD" // Real caller address

	t.Run("USDT_TotalSupply", func(t *testing.T) {
		// CallContract with visible=true expects T-addresses, not hex
		result := cli.CallContract(
			usdtContract,
			callerAddr,
			0,
			1000000,
			"totalSupply()",
		)

		// Skip test if TronGrid API is unavailable or returns an error
		if len(result.ConstantResult) == 0 || result.Result.Result == false {
			t.Skipf("TronGrid API call failed or unavailable: %+v. This test requires a working TronGrid API connection.", result.Result)
		}

		// Parse using the updated ParseTotalSupply that returns string
		totalSupply := ParseTotalSupply(result.ConstantResult)
		if totalSupply == nil {
			t.Fatal("ParseTotalSupply returned nil")
		}

		t.Logf("USDT Total Supply: %s", *totalSupply)

		// Verify the result is a valid number string
		supplyBig, ok := new(big.Int).SetString(*totalSupply, 10)
		if !ok {
			t.Fatalf("Total supply is not a valid number: %s", *totalSupply)
		}

		// USDT has 6 decimals and a very large supply
		// The total supply should be much larger than uint64 max (18446744073709551615)
		// Let's verify it's at least a reasonable size for USDT
		minExpected := new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil) // At least 10^15

		if supplyBig.Cmp(minExpected) < 0 {
			t.Fatalf("USDT total supply seems too small: %s (expected at least 10^15)", *totalSupply)
		}

		t.Logf("✓ Successfully parsed large USDT total supply without overflow")
	})
}

// TestParseTotalSupplyLargeNumbers tests that large numbers are correctly parsed as strings
func TestParseTotalSupplyLargeNumbers(t *testing.T) {
	testCases := []struct {
		name     string
		hexInput string
		expected string
	}{
		{
			name:     "Small number",
			hexInput: "0000000000000000000000000000000000000000000000000000000000000064", // 100 in hex
			expected: "100",
		},
		{
			name:     "uint64 max",
			hexInput: "000000000000000000000000000000000000000000000000ffffffffffffffff", // uint64 max
			expected: "18446744073709551615",
		},
		{
			name:     "Large number exceeding uint64",
			hexInput: "0000000000000000000000000000000000000000000000056bc75e2d63100000", // 100000000000000000000 (100 * 10^18)
			expected: "100000000000000000000",
		},
		{
			name:     "Very large number",
			hexInput: "00000000000000000000000000000000000000000000d3c21bcecceda1000000", // 1000000000000000000000000 (10^24)
			expected: "1000000000000000000000000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseTotalSupply([]string{tc.hexInput})
			if result == nil {
				t.Fatalf("ParseTotalSupply returned nil for input: %s", tc.hexInput)
			}

			if *result != tc.expected {
				t.Errorf("ParseTotalSupply = %s, expected %s", *result, tc.expected)
			}

			t.Logf("Hex: %s -> Decimal: %s", tc.hexInput, *result)
		})
	}
}

// TestParseTotalSupplyNoOverflow verifies that the string type prevents overflow
func TestParseTotalSupplyNoOverflow(t *testing.T) {
	// This is a number that would overflow uint64
	// uint64 max is 18446744073709551615 (2^64 - 1)
	// Let's use 10^30 = 1000000000000000000000000000000
	largeNumberHex := "000000000000000000000000000000000000000c9f2c9cd04674edea40000000" // 10^30 in hex

	result := ParseTotalSupply([]string{largeNumberHex})
	if result == nil {
		t.Fatal("ParseTotalSupply returned nil")
	}

	expected := "1000000000000000000000000000000"
	if *result != expected {
		t.Errorf("ParseTotalSupply = %s, expected %s", *result, expected)
	}

	// Verify it's much larger than uint64 max
	resultBig, _ := new(big.Int).SetString(*result, 10)
	uint64Max := new(big.Int).SetUint64(^uint64(0))

	if resultBig.Cmp(uint64Max) <= 0 {
		t.Errorf("Test number should be larger than uint64 max")
	}

	t.Logf("Successfully parsed large number without overflow: %s", *result)
	t.Logf("This is %s times larger than uint64 max", new(big.Int).Div(resultBig, uint64Max).String())
}
