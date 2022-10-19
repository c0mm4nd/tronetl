package main

import (
	"encoding/csv"
	"log"
	"math/big"
	"os"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
	"github.com/jszwec/csvutil"
)

type ExportBlocksAndTransactionsOptions struct {
}

func exportBlocksAndTransactions(providerURL string, start uint64, end uint64, blksOutput, txsOutput string) {
	cli := tron.NewTronClient(providerURL)

	blksOutFile, err := os.Create(blksOutput)
	chk(err)

	blksCsvWriter := csv.NewWriter(blksOutFile)
	defer blksCsvWriter.Flush()
	blksCsvEncoder := csvutil.NewEncoder(blksCsvWriter)

	txsOutFile, err := os.Create(txsOutput)
	chk(err)

	txsCsvWriter := csv.NewWriter(txsOutFile)
	defer txsCsvWriter.Flush()
	txsCsvEncoder := csvutil.NewEncoder(txsCsvWriter)

	for number := start; number <= end; number++ {
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
