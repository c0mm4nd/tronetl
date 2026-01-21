package tron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
)

type TronClient struct {
	httpURI string
	jsonURI string
}

// Shared HTTP client; keep timeout tight to avoid hanging when nodes are unresponsive.
var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func NewTronClient(providerURL string) *TronClient {
	httpURI := providerURL
	jsonURI := providerURL
	
	// Handle empty or invalid URLs
	if providerURL == "" {
		// Default to localhost with ports
		httpURI = "http://localhost:8090"
		jsonURI = "http://localhost:50545/jsonrpc"
	} else if strings.HasPrefix(providerURL, "https://") || strings.HasPrefix(providerURL, "http://") {
		// Already has protocol, use as-is
		httpURI = providerURL
		jsonURI = providerURL + "/jsonrpc"
	} else if strings.Contains(providerURL, ":") {
		// Already has port, use as-is
		httpURI = providerURL
		jsonURI = providerURL + "/jsonrpc"
	} else {
		// Legacy behavior: add default ports
		httpURI = providerURL + ":8090"
		jsonURI = providerURL + ":50545/jsonrpc"
	}
	
	return &TronClient{
		httpURI: httpURI,
		jsonURI: jsonURI,
	}
}

func (c *TronClient) GetJSONBlockByNumberWithTxs(number *big.Int) *JSONBlockWithTxs {
	payload, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params": []any{
			toBlockNumArg(number), true,
		},
		"id": rand.Int(),
	})
	chk(err)
	resp, err := httpClient.Post(c.jsonURI, "application/json", bytes.NewBuffer(payload))
	chk(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	chk(err)

	var rpcResp JSONResponse
	var block JSONBlockWithTxs
	err = json.Unmarshal(body, &rpcResp)
	chk(err)
	err = json.Unmarshal(rpcResp.Result, &block)
	chk(err)

	return &block
}

func (c *TronClient) GetJSONBlockByNumberWithTxIDs(number *big.Int) *JSONBlockWithTxIDs {
	payload, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params": []any{
			toBlockNumArg(number), false,
		},
		"id": rand.Int(),
	})
	chk(err)
	resp, err := httpClient.Post(c.jsonURI, "application/json", bytes.NewBuffer(payload))
	chk(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	chk(err)

	var rpcResp JSONResponse
	var block JSONBlockWithTxIDs
	err = json.Unmarshal(body, &rpcResp)
	chk(err)
	err = json.Unmarshal(rpcResp.Result, &block)
	chk(err)

	return &block
}

func (c *TronClient) GetHTTPBlockByNumber(number *big.Int) *HTTPBlock {
	url := c.httpURI + "/wallet/getblockbynum" // + "?visible=true"
	payload, err := json.Marshal(map[string]any{
		"num": number.Uint64(),
		// "visible": true,
	})
	chk(err)
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(payload))
	chk(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	chk(err)

	var block HTTPBlock
	err = json.Unmarshal(body, &block)
	chk(err)

	return &block
}

func (c *TronClient) GetTxInfosByNumber(number uint64) []HTTPTxInfo {
	if number == 0 {
		return []HTTPTxInfo{} // 0 height returns `{}` which is not a list
	}

	url := c.httpURI + "/wallet/gettransactioninfobyblocknum" // + "?visible=true"
	payload, err := json.Marshal(map[string]any{
		"num": number,
		// "visible": true,
	})
	chk(err)
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(payload))
	chk(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	chk(err)

	// Handle empty object response for blocks with no transactions
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "{}" {
		return []HTTPTxInfo{}
	}

	// Try to unmarshal as array first
	var txInfos []HTTPTxInfo
	err = json.Unmarshal(body, &txInfos)
	if err != nil {
		// If it fails, it might be an error response object
		// Try to unmarshal as a single object (error case)
		var singleInfo HTTPTxInfo
		err2 := json.Unmarshal(body, &singleInfo)
		if err2 == nil {
			// Successfully parsed as single object, but this shouldn't happen for valid data
			// Return empty array for now
			return []HTTPTxInfo{}
		}
		// Neither worked, panic with original error
		panic(fmt.Sprintf("failed to unmarshal tx infos for block %d: %v, body preview: %s", 
			number, err, string(body[:min(200, len(body))])))
	}

	return txInfos
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (c *TronClient) GetAccount(address string) *HTTPAccount {
	url := c.httpURI + "/wallet/getaccount" // + "?visible=true"
	payload, err := json.Marshal(map[string]any{
		"address": address,
		"visible": true,
	})
	chk(err)
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(payload))
	chk(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	chk(err)

	var acc HTTPAccount
	err = json.Unmarshal(body, &acc)
	chk(err)

	return &acc
}

func (c *TronClient) GetContract(address string) *HTTPContract {
	url := c.httpURI + "/wallet/getcontract" // + "?visible=true"
	payload, err := json.Marshal(map[string]any{
		"value":   address,
		"visible": true,
	})
	chk(err)
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(payload))
	chk(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	chk(err)

	var contract HTTPContract
	err = json.Unmarshal(body, &contract)
	chk(err)

	return &contract
}

type CallResult struct {
	Result struct {
		Result  bool   `json:"result,omitempty"`
		Code    string `json:"code,omitempty"` // contains "ERROR" when is error
		Message string `json:"message,omitempty"`
	} `json:"result,omitempty"`
	EnergyUsed     int              `json:"energy_used"`
	ConstantResult []string         `json:"constant_result"`
	Transaction    *HTTPTransaction `json:"transaction"`
}

type Address string

// Call is offline
func (c *TronClient) CallContract(contractAddr, callerAddr string, val, feeLimit int64, funcSig string, params ...any) *CallResult {
	url := c.httpURI + "/wallet/triggerconstantcontract" // + "?visible=true"
	u256Params := make([]string, len(params))
	for i, param := range params {
		switch p := param.(type) {
		case uint64:
			u256Params[i] = uint256.NewInt(p).Hex()[2:]
		case Address:
			addr := EnsureHexAddr(string(p))
			addr = addr[2:]

			u, err := uint256.FromHex(addr)
			if err != nil {
				panic(err)
			}
			u256Params[i] = u.Hex()[2:]
		default:
			panic(fmt.Sprintf("unsupported type: %#+v", param))
		}
	}
	payload, err := json.Marshal(map[string]any{
		"contract_address":  contractAddr,
		"function_selector": funcSig,
		"parameter":         "",
		"fee_limit":         feeLimit,
		"call_value":        val,
		"owner_address":     callerAddr, // = caller
		"visible":           true,
	})
	chk(err)
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(payload))
	chk(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	chk(err)

	var result CallResult
	err = json.Unmarshal(body, &result)
	chk(err)

	return &result
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}
