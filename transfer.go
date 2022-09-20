package main

import (
	"encoding/csv"
	"os"

	"github.com/jszwec/csvutil"
)

func exportTransfers(providerURL string, start uint64, end uint64, output string) {
	cli := NewTronHTTPClient(providerURL)

	outFile, err := os.Create(output)
	chk(err)

	w := csv.NewWriter(outFile)
	defer w.Flush()
	enc := csvutil.NewEncoder(w)

	for number := start; number <= end; number++ {
		txInfos := cli.GetTxInfosByNumber(number)
		for _, txInfo := range txInfos {
			txHash := txInfo.ID
			for logIndex, log := range txInfo.Log {
				tf := ExtractTransferFromLog(log.Topics, log.Data, log.Address, uint(logIndex), txHash, number)
				if tf != nil {
					err := enc.Encode(tf)
					chk(err)
				}

			}
		}
	}
}
