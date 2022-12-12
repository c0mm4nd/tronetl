package tron

import (
	"encoding/json"
	"math/big"
)

// HTTPTxInfo is a TxInfo result from HTTP RESTful API
// the struct follows https://github.com/tronprotocol/protocol/blob/2351aa6c2d708bf5ef47baf70410b3bc87d65fa7/core/Tron.proto#L341
type HTTPTxInfo struct {
	ID                            string                     `json:"id"`
	Fee                           int                        `json:"fee,omitempty"`
	BlockNumber                   int                        `json:"blockNumber"`
	BlockTimeStamp                int64                      `json:"blockTimeStamp"`
	ContractResult                []string                   `json:"contractResult"`
	ContractAddress               string                     `json:"contract_address"`
	Receipt                       *HTTPReceipt               `json:"receipt"`
	Log                           []*HTTPTxInfoLog           `json:"log,omitempty"`
	Result                        any                        `json:"result,omitempty"` // enum code { SUCESS = 0; FAILED = 1; }
	ResMessage                    string                     `json:"resMessage,omitempty"`
	AssetIssueID                  string                     `json:"assetIssueID,omitempty"`
	WithdrawAmount                int64                      `json:"withdraw_amount,omitempty"`
	UnfreezeAmount                int64                      `json:"unfreeze_amount,omitempty"`
	InternalTransactions          []*HTTPInternalTransaction `json:"internal_transactions,omitempty"`
	ExchangeReceivedAmount        int64                      `json:"exchange_received_amount,omitempty"`
	ExchangeInjectAnotherAmount   int64                      `json:"exchange_inject_another_amount,omitempty"`
	ExchangeWithdrawAnotherAmount int64                      `json:"exchange_withdraw_another_amount,omitempty"`
	ExchangeID                    int64                      `json:"exchange_id,omitempty"`
	ShieldedTransactionFee        int64                      `json:"shielded_transaction_fee,omitempty"`
}

type HTTPReceipt struct {
	EnergyUsage       int64  `json:"energy_usage,omitempty"`
	EnergyFee         int64  `json:"energy_fee,omitempty"`
	OriginEnergyUsage int64  `json:"origin_energy_usage,omitempty"`
	EnergyUsageTotal  int64  `json:"energy_usage_total,omitempty"`
	NetUsage          int64  `json:"net_usage,omitempty"`
	NetFee            int64  `json:"net_fee,omitempty"`
	Result            string `json:"result"`
}

// HTTPTxInfoLog is a Log result from HTTP RESTful API
type HTTPTxInfoLog struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
}

// HTTPInternalTransaction is a internal transaction result from HTTP RESTful API
type HTTPInternalTransaction struct {
	TransactionHash   string                                  `json:"hash"`
	CallerAddress     string                                  `json:"caller_address"`
	TransferToAddress string                                  `json:"transferTo_address"`
	CallValueInfo     []*HTTPInternalTransactionCallValueInfo `json:"callValueInfo,omitempty"`
	Note              string                                  `json:"note"`
	Rejected          bool                                    `json:"rejected"`
}

// HTTPInternalTransactionCallValueInfo is a field in HTTPInternalTransaction
// https://github.com/tronprotocol/protocol/blob/2351aa6c2d708bf5ef47baf70410b3bc87d65fa7/core/Tron.proto#L509
type HTTPInternalTransactionCallValueInfo struct {
	CallValue int64  `json:"callValue,omitempty"`
	TokenID   string `json:"tokenId,omitempty"`
}

// HTTPInternalTransaction is a Block result from HTTP RESTful API
// https://github.com/tronprotocol/protocol/blob/2351aa6c2d708bf5ef47baf70410b3bc87d65fa7/core/Tron.proto#L406
type HTTPBlock struct {
	BlockID      string            `json:"blockID"`
	BlockHeader  *HTTPBlockHeader  `json:"block_header"`
	Transactions []HTTPTransaction `json:"transactions"`
}

// HTTPBlockHeader represents the block header from the Block result from HTTP RESTful API
type HTTPBlockHeader struct {
	RawData struct {
		Timestamp  int64  `json:"timestamp,omitempty"`
		TxTrieRoot string `json:"txTrieRoot,omitempty"`
		ParentHash string `json:"parentHash"`

		Number           int    `json:"number"`
		WitnessAddress   string `json:"witness_address"`
		Version          int    `json:"version,omitempty"`
		AccountStateRoot string `json:"accountStateRoot,omitempty"`
	} `json:"raw_data"`
	WitnessSignature string `json:"witness_signature"`
}

// HTTPTransaction represents the Transaction result from HTTP RESTful API
type HTTPTransaction struct {
	Ret []struct {
		ContractRet string `json:"contractRet"`
	} `json:"ret"`
	Signature []string `json:"signature"`
	TxID      string   `json:"txID"`
	RawData   struct {
		Data          string          `json:"data"`
		Contract      []*ContractCall `json:"contract"`
		RefBlockBytes string          `json:"ref_block_bytes"`
		RefBlockHash  string          `json:"ref_block_hash"`
		Expiration    uint64          `json:"expiration"`
		Timestamp     uint64          `json:"timestamp"`
		FeeLimit      uint64          `json:"fee_limit"`
	} `json:"raw_data"`
	RawDataHex string `json:"raw_data_hex"`
}

// ContractCall represents a tron native contract call inside the Transaction
// Details of Parameter in https://tronprotocol.github.io/documentation-en/mechanism-algorithm/system-contracts/
type ContractCall struct {
	ContractType string `json:"type"`
	Parameter    struct {
		Value   json.RawMessage `json:"value"`
		TypeURL string          `json:"type_url"`
	} `json:"parameter"` // google.any decode with ContractType
	Provider     string `json:"provider"`
	PermissionID int32  `json:"Permission_id"`
}

// TRC10TransferParams can be the params of the following calls:
// - TransferAssetContract
// - TransferContract
type TRC10TransferParams struct {
	AssetName    string   `json:"asset_name"`
	Amount       *big.Int `json:"amount"`
	OwnerAddress string   `json:"owner_address"`
	ToAddress    string   `json:"to_address"`
}
