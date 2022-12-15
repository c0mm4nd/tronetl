package tron

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestCall(t *testing.T) {
	cli := NewTronClient("http://localhost")
	fmt.Println(cli.CallContract(
		Tstring2HexAddr("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"),
		Tstring2HexAddr("THPvaUhoh2Qn2y9THCZML3H815hhFhn5YC"),
		0,
		0,
		"totalSupply()",
	).ConstantResult)
	result := cli.CallContract(
		Tstring2HexAddr("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"),
		Tstring2HexAddr("THPvaUhoh2Qn2y9THCZML3H815hhFhn5YC"),
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
