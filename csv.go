package main

import "git.ngx.fi/c0mm4nd/tronetl/tron"

type CsvTransaction struct {
	Hash                 string `csv:"hash"`
	Nonce                string `csv:"nonce"`
	BlockHash            string `csv:"block_hash"`
	BlockNumber          uint64 `csv:"block_number"`
	TransactionIndex     int    `csv:"transaction_index"`
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
}

func NewCsvTransaction(blockTimestamp uint64, txIndex int, tx *tron.Transaction) *CsvTransaction {

	return &CsvTransaction{
		Hash:                 tx.Hash[2:],
		Nonce:                "", //tx.Nonce,
		BlockHash:            tx.BlockHash[2:],
		BlockNumber:          uint64(tx.BlockNumber),
		TransactionIndex:     txIndex,
		FromAddress:          hex2TAddr(tx.From[2:]),
		ToAddress:            hex2TAddr(tx.To[2:]),
		Value:                tx.Value.String(),
		Gas:                  tx.Gas.String(),
		GasPrice:             tx.GasPrice.String(),
		Input:                tx.Input[2:],
		BlockTimestamp:       blockTimestamp,
		MaxFeePerGas:         "", //tx.MaxFeePerGas.String(),
		MaxPriorityFeePerGas: "", //tx.MaxPriorityFeePerGas.String(),
		TransactionType:      tx.Type[2:],

		// Status: tx.,
	}
}

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
}

func NewCsvBlock(block tron.Block) *CsvBlock {
	return &CsvBlock{
		Number:           uint64(block.Number),
		Hash:             block.Hash[2:],
		ParentHash:       block.ParentHash[2:],
		Nonce:            "",
		Sha3Uncles:       "", // block.Sha3Uncles,
		LogsBloom:        block.LogsBloom[2:],
		TransactionsRoot: block.TransactionsRoot[2:],
		StateRoot:        block.StateRoot[2:],
		ReceiptsRoot:     "", // block.ReceiptsRoot
		Miner:            hex2TAddr(block.Miner[2:]),
		Difficulty:       "",
		TotalDifficulty:  "",
		Size:             uint64(block.Size),
		ExtraData:        "",
		GasLimit:         block.GasLimit.String(),
		GasUsed:          block.GasUsed.String(),
		Timestamp:        uint64(block.Timestamp),
		TansactionCount:  len(block.Transactions),
		BaseFeePerGas:    "", // block.BaseFeePerGas,
	}
}
