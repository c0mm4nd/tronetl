package main

import (
	"encoding/csv"
	"io"
	"log"
	"math/big"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
	"github.com/jszwec/csvutil"
	"golang.org/x/exp/slices"
)

type ExportTransferOptions struct {
	// outputType string // failed to tfOutput to a conn with
	tfOutput         io.Writer
	logOutput        io.Writer
	internalTxOutput io.Writer
	// withBlockOutput io.Writer

	ProviderURI string `json:"provider_uri,omitempty"`
	StartBlock  uint64 `json:"start_block,omitempty"`
	EndBlock    uint64 `json:"end_block,omitempty"`

	// extension
	StartTimestamp uint64 `json:"start_timestamp,omitempty"`
	EndTimestamp   uint64 `json:"end_timestamp,omitempty"`

	Contracts []string `json:"contracts,omitempty"`
}

// ExportTransfers is the main func for handling export_transfers command
func ExportTransfers(options *ExportTransferOptions) {
	cli := tron.NewTronClient(options.ProviderURI)

	tfWriter := csv.NewWriter(options.tfOutput)
	defer tfWriter.Flush()
	tfEncoder := csvutil.NewEncoder(tfWriter)

	logWriter := csv.NewWriter(options.tfOutput)
	defer logWriter.Flush()
	logEncoder := csvutil.NewEncoder(logWriter)

	internalTxWriter := csv.NewWriter(options.tfOutput)
	defer internalTxWriter.Flush()
	internalTxEncoder := csvutil.NewEncoder(internalTxWriter)

	filterLogContracts := make([]string, len(options.Contracts))
	for i, addr := range options.Contracts {
		filterLogContracts[i] = Tstring2HexAddr(addr)[2:] // hex addr with 41 prefix
	}

	if options.StartTimestamp != 0 {
		// fast locate estimate start height

		estimateStartNumber := locateStartBlock(cli, options.StartTimestamp)

		for number := estimateStartNumber; ; number++ {
			block := cli.GetJSONBlockByNumberWithTxIDs(new(big.Int).SetUint64(number))
			if block == nil {
				break
			}

			blockTime := uint64(*block.Timestamp) / 1000

			if blockTime < options.StartTimestamp {
				log.Printf("passed start block %d: %d", number, *block.Timestamp)
				continue
			}

			options.StartBlock = number
			break
		}
	}

	if options.EndTimestamp != 0 {
		estimateEndNumber := locateEndBlock(cli, options.EndTimestamp)

		for number := estimateEndNumber; ; number-- {
			block := cli.GetJSONBlockByNumberWithTxIDs(new(big.Int).SetUint64(number))
			if block == nil {
				break
			}

			blockTime := uint64(*block.Timestamp) / 1000

			if blockTime > options.EndBlock {
				log.Printf("passed end block %d: %d", number, *block.Timestamp)
				continue
			}

			options.StartBlock = number
			break
		}
	}

	log.Printf("try parsing token transfers from block %d to %d", options.StartBlock, options.EndBlock)

	if options.StartBlock != 0 && options.EndBlock != 0 {
		for number := options.StartBlock; number <= options.EndBlock; number++ {
			txInfos := cli.GetTxInfosByNumber(number)
			for _, txInfo := range txInfos {
				txHash := txInfo.ID
				for logIndex, log := range txInfo.Log {
					if len(filterLogContracts) != 0 && !slices.Contains(filterLogContracts, log.Address) {
						continue
					}

					if options.tfOutput != nil {
						tf := ExtractTransferFromLog(log.Topics, log.Data, log.Address, uint(logIndex), txHash, number)
						if tf != nil {
							err := tfEncoder.Encode(tf)
							chk(err)
						}
					}

					if options.logOutput != nil {
						err := logEncoder.Encode(NewCsvLog(number, txHash, uint(logIndex), log))
						chk(err)
					}

				}

				if options.internalTxOutput != nil {
					for internalIndex, internalTx := range txInfo.InternalTransactions {
						err := internalTxEncoder.Encode(NewCsvInternalTx(uint(internalIndex), internalTx))
						chk(err)
					}
				}

			}

			log.Printf("parsed block %d", number)
		}

		return
	}

}
