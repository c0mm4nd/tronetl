package tron

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"
)

func TestCall(t *testing.T) {
	provider := os.Getenv("TRON_PROVIDER")
	if provider == "" {
		t.Skip("TRON_PROVIDER not set; skipping integration test")
	}

	cli := NewTronClient(provider)
	fmt.Println(cli.CallContract(
		EnsureHexAddr("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"),
		EnsureHexAddr("THPvaUhoh2Qn2y9THCZML3H815hhFhn5YC"),
		0,
		0,
		"totalSupply()",
	).ConstantResult)
	result := cli.CallContract(
		EnsureHexAddr("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"),
		EnsureHexAddr("THPvaUhoh2Qn2y9THCZML3H815hhFhn5YC"),
		0,
		0,
		"symbol()",
	).ConstantResult[0]
	for i := 0; i+64 <= len(result); i += 64 {
		fmt.Println(result[i : i+64])
		decoded, _ := hex.DecodeString(result[i : i+64])
		fmt.Println(new(big.Int).SetBytes(decoded))
		fmt.Println(string(decoded))
	}
}
