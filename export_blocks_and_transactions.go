package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"math/big"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
	"github.com/jszwec/csvutil"
)

// ExportBlocksAndTransactionsOptions is the option for ExportBlocksAndTransactions func
type ExportBlocksAndTransactionsOptions struct {
	blksOutput  io.Writer
	txsOutput   io.Writer
	trc10Output io.Writer

	ProviderURI string `json:"provider_uri,omitempty"`
	StartBlock  uint64 `json:"start_block,omitempty"`
	EndBlock    uint64 `json:"end_block,omitempty"`

	// extension
	StartTimestamp uint64 `json:"start_timestamp,omitempty"`
	EndTimestamp   uint64 `json:"end_timestamp,omitempty"`
}

// ExportBlocksAndTransactions is the main func for handling export_blocks_and_transactions command
func ExportBlocksAndTransactions(options *ExportBlocksAndTransactionsOptions) {
	cli := tron.NewTronClient(options.ProviderURI)

	var blksCsvEncoder, txsCsvEncoder, trc10CsvEncoder *csvutil.Encoder
	if options.blksOutput != nil {
		blksCsvWriter := csv.NewWriter(options.blksOutput)
		defer blksCsvWriter.Flush()
		blksCsvEncoder = csvutil.NewEncoder(blksCsvWriter)
	}

	if options.txsOutput != nil {
		txsCsvWriter := csv.NewWriter(options.txsOutput)
		defer txsCsvWriter.Flush()
		txsCsvEncoder = csvutil.NewEncoder(txsCsvWriter)
	}

	if options.trc10Output != nil {
		trc10CsvWriter := csv.NewWriter(options.trc10Output)
		defer trc10CsvWriter.Flush()
		trc10CsvEncoder = csvutil.NewEncoder(trc10CsvWriter)
	}

	for number := options.StartBlock; number <= options.EndBlock; number++ {
		num := new(big.Int).SetUint64(number)

		jsonblock := cli.GetJSONBlockByNumberWithTxs(num)
		httpblock := cli.GetHTTPBlockByNumber(num)
		csvBlock := NewCsvBlock(jsonblock, httpblock)
		if options.txsOutput != nil || options.trc10Output != nil {
			for txIndex, jsontx := range jsonblock.Transactions {
				httptx := httpblock.Transactions[txIndex]
				if options.txsOutput != nil {
					csvTx := NewCsvTransaction(uint64(*jsonblock.Timestamp), txIndex, &jsontx, &httptx)
					err := txsCsvEncoder.Encode(csvTx)
					chk(err)
				}

				if options.trc10Output != nil {
					for callIndex, contractCall := range httptx.RawData.Contract {
						if contractCall.ContractType == "TransferAssetContract" ||
							contractCall.ContractType == "TransferContract" {
							var tfParams tron.TRC10TransferParams

							err := json.Unmarshal(contractCall.Parameter.Value, &tfParams)
							chk(err)
							csvTf := NewCsvTRC10Transfer(number, txIndex, callIndex, &httpblock.Transactions[txIndex], &tfParams)
							err = trc10CsvEncoder.Encode(csvTf)
							chk(err)
						}
					}

				}

			}
		}

		err := blksCsvEncoder.Encode(csvBlock)
		chk(err)

		log.Printf("parsed block %d", number)
	}
}
