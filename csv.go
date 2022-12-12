package main

import (
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
}

// NewCsvTransaction creates a new CsvTransaction
func NewCsvTransaction(blockTimestamp uint64, txIndex int, jsontx *tron.JSONTransaction, httptx *tron.HTTPTransaction) *CsvTransaction {
	to := ""
	if jsontx.To != "" {
		to = Hex2TAddr(jsontx.To[2:])
	}

	txType := "Unknown"
	if len(httptx.RawData.Contract) > 0 {
		txType = httptx.RawData.Contract[0].ContractType
	}

	return &CsvTransaction{
		Hash:                 jsontx.Hash[2:],
		Nonce:                "", //tx.Nonce,
		BlockHash:            jsontx.BlockHash[2:],
		BlockNumber:          uint64(*jsontx.BlockNumber),
		TransactionIndex:     txIndex,
		FromAddress:          Hex2TAddr(jsontx.From[2:]),
		ToAddress:            to,
		Value:                jsontx.Value.ToInt().String(),
		Gas:                  jsontx.Gas.ToInt().String(),
		GasPrice:             jsontx.GasPrice.ToInt().String(), // https://support.ledger.com/hc/en-us/articles/6331588714141-How-do-Tron-TRX-fees-work-?support=true
		Input:                jsontx.Input[2:],
		BlockTimestamp:       blockTimestamp / 1000, // unit: sec
		MaxFeePerGas:         "",                    //tx.MaxFeePerGas.String(),
		MaxPriorityFeePerGas: "",                    //tx.MaxPriorityFeePerGas.String(),
		TransactionType:      txType,                //jsontx.Type[2:],

		Status: httptx.Ret[0].ContractRet, // can be SUCCESS REVERT

		// appendix
		TransactionTimestamp:  httptx.RawData.Timestamp / 1000,  // float64(httptx.RawData.Timestamp) * 1 / 1000,
		TransactionExpiration: httptx.RawData.Expiration / 1000, // float64(httptx.RawData.Expiration) * 1 / 1000,
		FeeLimit:              httptx.RawData.FeeLimit,
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
		Timestamp:        uint64(*jsonblock.Timestamp) / 1000,
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
	BlockNumber     uint64 `csv:"blockNumber" csv:"block_number"`
	TransactionHash string `csv:"transaction_hash" csv:"transaction_hash"`
	LogIndex        uint   `csv:"logIndex" csv:"log_index"`

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
	CallInfoIndex     uint   `csv:"call_info_index"`
	CallTokenID       string `csv:"call_token_id"`
	CallValue         int64  `csv:"call_value"`
	Note              string `csv:"note"`
	Rejected          bool   `csv:"rejected"`
}

// NewCsvInternalTx creates a new CsvInternalTx
func NewCsvInternalTx(index uint, itx *tron.HTTPInternalTransaction, callInfoIndex uint, tokenID string, value int64) *CsvInternalTx {

	return &CsvInternalTx{
		TransactionHash:   itx.TransactionHash,
		Index:             index,
		CallerAddress:     Hex2TAddr(itx.CallerAddress),
		TransferToAddress: Hex2TAddr(itx.TransferToAddress),
		// CallValueInfo:     strings.Join(callValues, ";"),
		CallInfoIndex: callInfoIndex,
		CallTokenID:   tokenID,
		CallValue:     value,

		Note:     itx.Note,
		Rejected: itx.Rejected,
	}
}

// CsvReceipt is a receipt for tron transaction
type CsvReceipt struct {
	TxHash            string `csv:"transaction_hash"`
	EnergyUsage       int64  `csv:"energy_usage,omitempty"`
	EnergyFee         int64  `csv:"energy_fee,omitempty"`
	OriginEnergyUsage int64  `csv:"origin_energy_usage,omitempty"`
	EnergyUsageTotal  int64  `csv:"energy_usage_total,omitempty"`
	NetUsage          int64  `csv:"net_usage,omitempty"`
	NetFee            int64  `csv:"net_fee,omitempty"`
	Result            string `csv:"result"`
}

func NewCsvReceipt(txHash string, r *tron.HTTPReceipt) *CsvReceipt {

	return &CsvReceipt{
		TxHash:            txHash,
		EnergyUsage:       r.EnergyUsage,
		EnergyFee:         r.EnergyFee,
		OriginEnergyUsage: r.OriginEnergyUsage,
		EnergyUsageTotal:  r.EnergyUsageTotal,
		NetUsage:          r.NetUsage,
		NetFee:            r.NetFee,
		Result:            r.Result,
	}
}
