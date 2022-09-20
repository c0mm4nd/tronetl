package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type TronHTTPClient struct {
	providerURL string
}

func NewTronHTTPClient(providerURL string) *TronHTTPClient {
	return &TronHTTPClient{
		providerURL: providerURL,
	}
}

type Transaction struct {
	Ret []struct {
		ContractRet string `json:"contractRet"`
	} `json:"ret"`
	Signature  []string `json:"signature"`
	TxID       string   `json:"txID"`
	RawDataHex string   `json:"raw_data_hex"`
	RawData    struct {
		Contract []struct {
			Parameter struct {
				Value struct {
					Data            string `json:"data"`
					OwnerAddress    string `json:"owner_address"`
					ContractAddress string `json:"contract_address"`
				} `json:"value"`
				TypeURL string `json:"type_url"`
			} `json:"parameter"`
			Type string `json:"type"`
		} `json:"contract"`
		RefBlockBytes  string `json:"ref_block_bytes"`
		RefBlockHash   string `json:"ref_block_hash"`
		Expiration     int64  `json:"expiration"`
		Timestamp      int64  `json:"timestamp"`
		FeeLimit       int    `json:"fee_limit"`
		Number         int    `json:"number"`
		TxTrieRoot     string `json:"txTrieRoot"`
		WitnessAddress string `json:"witness_address"`
		ParentHash     string `json:"parentHash"`
		Version        int    `json:"version"`
	} `json:"raw_data,omitempty"`
}

type Block struct {
	BlockID     string `json:"blockID"`
	BlockHeader struct {
		RawData struct {
			Number         int    `json:"number"`
			TxTrieRoot     string `json:"txTrieRoot"`
			WitnessAddress string `json:"witness_address"`
			ParentHash     string `json:"parentHash"`
			Version        int    `json:"version"`
			Timestamp      int64  `json:"timestamp"`
		} `json:"raw_data"`
		WitnessSignature string `json:"witness_signature"`
	} `json:"block_header"`
	Transactions []Transaction `json:"transactions"`
}

func (c *TronHTTPClient) GetBlockByNumber(number uint64) *Block {
	url := c.providerURL + "/wallet/getblockbynum"
	payload, err := json.Marshal(map[string]any{
		"num": number,
	})
	chk(err)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	chk(err)
	body, err := io.ReadAll(resp.Body)
	chk(err)

	var block Block
	err = json.Unmarshal(body, &block)
	chk(err)

	return &block
}

type TxInfo struct {
	Log []struct {
		Address string   `json:"address"`
		Data    string   `json:"data"`
		Topics  []string `json:"topics"`
	} `json:"log,omitempty"`
	Fee            int      `json:"fee,omitempty"`
	BlockNumber    int      `json:"blockNumber"`
	ContractResult []string `json:"contractResult"`
	BlockTimeStamp int64    `json:"blockTimeStamp"`
	Receipt        struct {
		Result            string `json:"result"`
		NetFee            int    `json:"net_fee"`
		EnergyUsageTotal  int    `json:"energy_usage_total"`
		OriginEnergyUsage int    `json:"origin_energy_usage"`
	} `json:"receipt"`
	ID                   string `json:"id"`
	ContractAddress      string `json:"contract_address,omitempty"`
	InternalTransactions []struct {
		CallerAddress     string `json:"caller_address"`
		Note              string `json:"note"`
		TransferToAddress string `json:"transferTo_address"`
		CallValueInfo     []struct {
		} `json:"callValueInfo"`
		Hash string `json:"hash"`
	} `json:"internal_transactions,omitempty"`
}

func (c *TronHTTPClient) GetTxInfosByNumber(number uint64) []TxInfo {
	url := c.providerURL + "/wallet/gettransactioninfobyblocknum"
	payload, err := json.Marshal(map[string]any{
		"num": number,
	})
	chk(err)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	chk(err)
	body, err := io.ReadAll(resp.Body)
	chk(err)

	var txInfos []TxInfo
	err = json.Unmarshal(body, &txInfos)
	chk(err)

	return txInfos
}
