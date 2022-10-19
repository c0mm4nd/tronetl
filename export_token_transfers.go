package main

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
	"github.com/jszwec/csvutil"
)

func exportTransfers(providerURI string, start uint64, end uint64, output string) {
	cli := tron.NewTronClient(providerURI)

	outFile, err := os.Create(output)
	chk(err)

	w := csv.NewWriter(outFile)
	defer w.Flush()
	enc := csvutil.NewEncoder(w)

	loc, _ := time.LoadLocation("Asia/Shanghai")

	for number := start; number <= end; number++ {
		// temp block time filter
		block := cli.GetBlockByNumber(number)
		blockTime := time.Unix(int64(*block.Timestamp)/1000, 0)
		y, m, d := blockTime.In(loc).Date()
		if !(y == 2022 && m == 06 && (d == 6 || d == 5)) {
			log.Printf("passed block %d: %d = %d-%d-%d", number, *block.Timestamp, y, m, d)
			continue
		}

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

		log.Printf("parsed block %d", number)
	}
}
