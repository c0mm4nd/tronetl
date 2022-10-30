package main

import (
	"encoding/csv"
	"io"
	"log"
	"math/big"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
	"github.com/jszwec/csvutil"
)

type ExportBlocksAndTransactionsOptions struct {
	// outputType string // failed to output to a conn with
	blksOutput io.Writer
	txsOutput  io.Writer
	// withBlockOutput io.Writer

	ProviderURI string `json:"provider_uri,omitempty"`
	StartBlock  uint64 `json:"start_block,omitempty"`
	EndBlock    uint64 `json:"end_block,omitempty"`

	// extension
	StartTimestamp uint64 `json:"start_timestamp,omitempty"`
	EndTimestamp   uint64 `json:"end_timestamp,omitempty"`
}

func exportBlocksAndTransactions(options *ExportBlocksAndTransactionsOptions) {
	cli := tron.NewTronClient(options.ProviderURI)

	blksCsvWriter := csv.NewWriter(options.blksOutput)
	defer blksCsvWriter.Flush()
	blksCsvEncoder := csvutil.NewEncoder(blksCsvWriter)

	txsCsvWriter := csv.NewWriter(options.txsOutput)
	defer txsCsvWriter.Flush()
	txsCsvEncoder := csvutil.NewEncoder(txsCsvWriter)

	for number := options.StartBlock; number <= options.EndBlock; number++ {
		num := new(big.Int).SetUint64(number)
		jsonblock := cli.GetJSONBlockByNumber(num)
		httpblock := cli.GetHTTPBlockByNumber(num)
		csvBlock := NewCsvBlock(jsonblock, httpblock)
		for txIndex, transaction := range jsonblock.Transactions {
			csvTx := NewCsvTransaction(uint64(*jsonblock.Timestamp), txIndex, &transaction, &httpblock.Transactions[txIndex])
			txsCsvEncoder.Encode(csvTx)
		}

		blksCsvEncoder.Encode(csvBlock)
		log.Printf("parsed block %d", number)
	}
}
