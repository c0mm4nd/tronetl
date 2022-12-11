package main

import (
	"strconv"
	"strings"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
)

// CsvTransaction represents a tron tx csv output, not trc10
// 1 TRX = 1000000 sun
type CsvTransaction struct {
	Hash             string `csv:"hash"`
	Nonce            string `csv:"nonce"`
	BlockHash        string `csv:"block_hash"`
	BlockNumber      uint64 `csv:"block_number"`
	TransactionIndex int    `csv:"transaction_index"`

	FromAddress          string `csv:"from_address"`
	ToAddress            string `csv:"to_address"`
	Value                string `csv:"value"`
	Gas                  string `csv:"gas"`
	GasPrice             string `csv:"gas_price"`
	Input                string `csv:"input"`
	BlockTimestamp       uint64 `csv:"block_timestamp"`
	MaxFeePerGas         string `csv:"max_fee_per_gas"`
	MaxPriorityFeePerGas string `csv:"max_priority_fee_per_gas"`
	TransactionType      string `csv:"transaction_type"`

	Status string

	// appendix
	TransactionTimestamp  uint64 `csv:"transaction_timestamp"`
	TransactionExpiration uint64 `csv:"transaction_expiration"`
	FeeLimit              uint64 `csv:"fee_limit"`
	ContractCallCount     int    `csv:"contract_calls"`
}

// NewCsvTransaction creates a new CsvTransaction
func NewCsvTransaction(blockTimestamp uint64, txIndex int, jsontx *tron.JSONTransaction, httptx *tron.HTTPTransaction) *CsvTransaction {
	to := ""
	if jsontx.To != "" {
		to = Hex2TAddr(jsontx.To[2:])
	}

	return &CsvTransaction{
		Hash:                 jsontx.Hash[2:],
		Nonce:                "", //tx.Nonce,
		BlockHash:            jsontx.BlockHash[2:],
		BlockNumber:          uint64(*jsontx.BlockNumber),
		TransactionIndex:     txIndex,
		FromAddress:          Hex2TAddr(jsontx.From[2:]),
		ToAddress:            to,
		Value:                jsontx.Value.String(),
		Gas:                  jsontx.Gas.String(),
		GasPrice:             jsontx.GasPrice.String(), // https://support.ledger.com/hc/en-us/articles/6331588714141-How-do-Tron-TRX-fees-work-?support=true
		Input:                jsontx.Input[2:],
		BlockTimestamp:       blockTimestamp,
		MaxFeePerGas:         "", //tx.MaxFeePerGas.String(),
		MaxPriorityFeePerGas: "", //tx.MaxPriorityFeePerGas.String(),
		TransactionType:      jsontx.Type[2:],

		Status: httptx.Ret[0].ContractRet,

		// appendix
		TransactionTimestamp:  httptx.RawData.Timestamp,
		TransactionExpiration: httptx.RawData.Expiration,
		FeeLimit:              httptx.RawData.FeeLimit,
		ContractCallCount:     len(httptx.RawData.Contract),
	}
}

// CsvBlock represents a tron block output
type CsvBlock struct {
	Number           uint64 `csv:"number"`
	Hash             string `csv:"hash"`
	ParentHash       string `csv:"parent_hash"`
	Nonce            string `csv:"nonce"`
	Sha3Uncles       string `csv:"sha3_uncles"`
	LogsBloom        string `csv:"logs_bloom"`
	TransactionsRoot string `csv:"transaction_root"`
	StateRoot        string `csv:"state_root"`
	ReceiptsRoot     string `csv:"receipts_root"`
	Miner            string `csv:"miner"`
	Difficulty       string `csv:"difficulty"`
	TotalDifficulty  string `csv:"total_difficulty"`
	Size             uint64 `csv:"size"`
	ExtraData        string `csv:"extra_data"`
	GasLimit         string `csv:"gas_limit"`
	GasUsed          string `csv:"gas_used"`
	Timestamp        uint64 `csv:"timestamp"`
	TansactionCount  int    `csv:"transaction_count"`
	BaseFeePerGas    string `csv:"base_fee_per_gas"`

	// append
	WitnessSignature string `csv:"witness_signature"`
}

// NewCsvBlock creates a new CsvBlock
func NewCsvBlock(jsonblock *tron.JSONBlockWithTxs, httpblock *tron.HTTPBlock) *CsvBlock {
	return &CsvBlock{
		Number:           uint64(*jsonblock.Number),
		Hash:             jsonblock.Hash[2:],
		ParentHash:       jsonblock.ParentHash[2:],
		Nonce:            "",
		Sha3Uncles:       "", // block.Sha3Uncles,
		LogsBloom:        jsonblock.LogsBloom[2:],
		TransactionsRoot: jsonblock.TransactionsRoot[2:],
		StateRoot:        jsonblock.StateRoot[2:],
		ReceiptsRoot:     "",                             // block.ReceiptsRoot
		Miner:            Hex2TAddr(jsonblock.Miner[2:]), // = WitnessAddress
		Difficulty:       "",
		TotalDifficulty:  "",
		Size:             uint64(*jsonblock.Size),
		ExtraData:        "",
		GasLimit:         jsonblock.GasLimit.ToInt().String(),
		GasUsed:          jsonblock.GasUsed.ToInt().String(),
		Timestamp:        uint64(*jsonblock.Timestamp),
		TansactionCount:  len(jsonblock.Transactions),
		BaseFeePerGas:    "", // block.BaseFeePerGas,

		//append
		WitnessSignature: httpblock.BlockHeader.WitnessSignature,
	}
}

// CsvTRC10Transfer is a trc10 transfer output
// https://developers.tron.network/docs/trc10-transfer-in-smart-contracts
// https://tronprotocol.github.io/documentation-en/mechanism-algorithm/system-contracts/
// It represents:
// - TransferContract
// - TransferAssetContract
type CsvTRC10Transfer struct {
	BlockNumber       uint64 `csv:"block_number"`
	BlockHash         string `csv:"block_hash"`
	TransactionHash   string `csv:"transaction_hash"`
	TransactionIndex  int    `csv:"transaction_index"`
	ContractCallIndex int    `csv:"contract_call_index"`

	AssetName   string `csv:"asset_name"` // do not omit => empty means trx
	FromAddress string `csv:"from_address"`
	ToAddress   string `csv:"to_address"`
	Value       string `csv:"value"`
}

// NewCsvTRC10Transfer creates a new CsvTRC10Transfer
func NewCsvTRC10Transfer(blockNum uint64, txIndex, callIndex int, httpTx *tron.HTTPTransaction, tfParams *tron.TRC10TransferParams) *CsvTRC10Transfer {

	return &CsvTRC10Transfer{
		TransactionHash:   httpTx.TxID,
		BlockHash:         httpTx.RawData.RefBlockHash,
		BlockNumber:       blockNum,
		TransactionIndex:  txIndex,
		ContractCallIndex: callIndex,

		AssetName:   tfParams.AssetName,
		FromAddress: Hex2TAddr(tfParams.OwnerAddress),
		ToAddress:   Hex2TAddr(tfParams.ToAddress),
		Value:       tfParams.Amount.String(),
	}
}

// CsvLog is a EVM smart contract event log output
type CsvLog struct {
	BlockNumber     uint64 `json:"blockNumber" csv:"block_number"`
	TransactionHash string `json:"transaction_hash" csv:"transaction_hash"`
	LogIndex        uint   `json:"logIndex" csv:"log_index"`

	Address string `csv:"address"`
	Topics  string `csv:"topics"`
	Data    string `csv:"data"`
}

// NewCsvLog creates a new CsvLog
func NewCsvLog(blockNumber uint64, txHash string, logIndex uint, log *tron.HTTPTxInfoLog) *CsvLog {
	return &CsvLog{
		BlockNumber:     blockNumber,
		TransactionHash: txHash,
		LogIndex:        logIndex,

		Address: Hex2TAddr(log.Address),
		Topics:  strings.Join(log.Topics, ";"),
		Data:    log.Data,
	}
}

// CsvInternalTx is a EVM smart contract internal transaction
type CsvInternalTx struct {
	TransactionHash   string `csv:"transaction_hash"`
	Index             uint   `csv:"internal_index"`
	CallerAddress     string `csv:"caller_address"`
	TransferToAddress string `csv:"transferTo_address"`
	CallValueInfo     string `csv:"callValueInfo"`
	Note              string `csv:"note"`
	Rejected          bool   `csv:"rejected"`
}

// NewCsvInternalTx creates a new CsvInternalTx
func NewCsvInternalTx(index uint, itx *tron.HTTPInternalTransaction) *CsvInternalTx {
	callValues := make([]string, len(itx.CallValueInfo))
	for i, callValue := range itx.CallValueInfo {
		callValues[i] = callValue.TokenID + ":" + strconv.FormatInt(callValue.CallValue, 10)
	}

	return &CsvInternalTx{
		TransactionHash:   itx.TransactionHash,
		Index:             index,
		CallerAddress:     Hex2TAddr(itx.CallerAddress),
		TransferToAddress: Hex2TAddr(itx.TransferToAddress),
		CallValueInfo:     strings.Join(callValues, ";"),
		Note:              itx.Note,
		Rejected:          itx.Rejected,
	}
}
