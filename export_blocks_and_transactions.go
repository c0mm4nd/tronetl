package main

func exportBlocksAndTransactions(providerURL string, start uint64, end uint64, blksOutput, txsOutput string) {
	// cli := tron.NewTronClient(providerURL)

	// blksOutFile, err := os.Create(blksOutput)
	// chk(err)

	// blksCsvWriter := csv.NewWriter(blksOutFile)
	// defer blksCsvWriter.Flush()
	// blksCsvEncoder := csvutil.NewEncoder(blksCsvWriter)

	// txsOutFile, err := os.Create(txsOutput)
	// chk(err)

	// txsCsvWriter := csv.NewWriter(txsOutFile)
	// defer txsCsvWriter.Flush()
	// txsCsvEncoder := csvutil.NewEncoder(txsCsvWriter)

	// for number := start; number <= end; number++ {
	// 	block := cli.GetBlockByNumber(number)
	// 	csvBlock := NewCsvBlock(block)
	// }
}
