package tron

import (
	"bytes"
	"encoding/json"
	"io"
	"math/big"
	"math/rand"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
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

func (c *TronClient) GetTxInfosByNumber(number uint64) []HTTPTxInfo {
	url := c.httpURI + "/wallet/gettransactioninfobyblocknum"
	payload, err := json.Marshal(map[string]any{
		"num": number,
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

func (c *TronClient) GetJSONBlockByNumber(number *big.Int, requireDetail bool) *JSONBlock {
	payload, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params": []any{
			toBlockNumArg(number), requireDetail,
		},
		"id": rand.Int(),
	})
	chk(err)
	resp, err := http.Post(c.jsonURI, "application/json", bytes.NewBuffer(payload))
	chk(err)
	body, err := io.ReadAll(resp.Body)
	chk(err)

	var rpcResp JSONResponse
	var block JSONBlock
	err = json.Unmarshal(body, &rpcResp)
	chk(err)
	err = json.Unmarshal(rpcResp.Result, &block)
	chk(err)

	return &block
}

func (c *TronClient) GetHTTPBlockByNumber(number *big.Int) *HTTPBlock {
	url := c.httpURI + "/wallet/getblockbynum"
	payload, err := json.Marshal(map[string]any{
		"num": number.Uint64(),
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
