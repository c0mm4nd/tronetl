package tron

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"

	"golang.org/x/crypto/sha3"
)

// FunctionSelector returns 4-byte selector for a function signature, with 0x prefix.
func FunctionSelector(sig string) string {
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(sig))
	sum := h.Sum(nil)
	return "0x" + hex.EncodeToString(sum[:4])
}

// EthCall executes eth_call via JSON-RPC and returns hex output (0x...)
func (c *TronClient) EthCall(toAddr, data string) (string, error) {
	if data != "" && !strings.HasPrefix(data, "0x") {
		data = "0x" + data
	}
	toHex := EnsureHexAddr(toAddr)
	if !strings.HasPrefix(toHex, "0x") {
		toHex = "0x" + toHex
	}

	payload, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"method":  "eth_call",
		"params": []any{
			map[string]any{"to": toHex, "data": data},
			"latest",
		},
		"id": rand.Int(),
	})
	if err != nil {
		return "", err
	}

	resp, err := httpClient.Post(c.jsonURI, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var rpcResp struct {
		JSONRPC string  `json:"jsonrpc"`
		ID      int     `json:"id"`
		Result  *string `json:"result"`
		Error   *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return "", fmt.Errorf("eth_call unmarshal error: %w", err)
	}
	if rpcResp.Error != nil {
		return "", fmt.Errorf("eth_call error: %s", rpcResp.Error.Message)
	}
	if rpcResp.Result == nil {
		return "", errors.New("eth_call empty result")
	}
	return *rpcResp.Result, nil
}
