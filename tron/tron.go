package tron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
)

type TronClient struct {
	httpURI string
	jsonURI string
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func NewTronClient(providerURL string) *TronClient {
	return &TronClient{
		httpURI: providerURL + ":8090",
		jsonURI: providerURL + ":50545/jsonrpc",
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
	resp, err := http.Post(c.jsonURI, "application/json", bytes.NewBuffer(payload))
	chk(err)
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
	resp, err := http.Post(c.jsonURI, "application/json", bytes.NewBuffer(payload))
	chk(err)
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
		// "visable": true,
	})
	chk(err)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	chk(err)
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
		// "visable": true,
	})
	chk(err)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	chk(err)
	body, err := io.ReadAll(resp.Body)
	chk(err)

	var txInfos []HTTPTxInfo
	err = json.Unmarshal(body, &txInfos)
	chk(err)

	return txInfos
}

func (c *TronClient) GetAccount(address string) *HTTPAccount {
	url := c.httpURI + "/wallet/getaccount" // + "?visible=true"
	payload, err := json.Marshal(map[string]any{
		"address": address,
		// "visable": true,
	})
	chk(err)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	chk(err)
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
		"value": address,
		// "visable": true,
	})
	chk(err)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	chk(err)
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
		// "visable": true,
	})
	chk(err)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	chk(err)
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
