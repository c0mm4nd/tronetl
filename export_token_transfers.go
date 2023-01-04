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
	receiptOutput    io.Writer

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

	var tfEncoder, logEncoder, internalTxEncoder, receiptEncoder *csvutil.Encoder

	if options.tfOutput != nil {
		tfWriter := csv.NewWriter(options.tfOutput)
		defer tfWriter.Flush()
		tfEncoder = csvutil.NewEncoder(tfWriter)
	}

	if options.logOutput != nil {
		logWriter := csv.NewWriter(options.logOutput)
		defer logWriter.Flush()
		logEncoder = csvutil.NewEncoder(logWriter)
	}

	if options.internalTxOutput != nil {
		internalTxWriter := csv.NewWriter(options.internalTxOutput)
		defer internalTxWriter.Flush()
		internalTxEncoder = csvutil.NewEncoder(internalTxWriter)
	}

	if options.receiptOutput != nil {
		receiptWriter := csv.NewWriter(options.receiptOutput)
		defer receiptWriter.Flush()
		receiptEncoder = csvutil.NewEncoder(receiptWriter)
	}

	filterLogContracts := make([]string, len(options.Contracts))
	for i, addr := range options.Contracts {
		filterLogContracts[i] = tron.EnsureHexAddr(addr)[2:] // hex addr with 41 prefix
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
			for txIndex, txInfo := range txInfos {
				txHash := txInfo.ID

				if receiptEncoder != nil {
					receiptEncoder.Encode(NewCsvReceipt(number, txHash, uint(txIndex), txInfo.ContractAddress, txInfo.Receipt))
				}

				for logIndex, log := range txInfo.Log {
					if len(filterLogContracts) != 0 && !slices.Contains(filterLogContracts, log.Address) {
						continue
					}

					if tfEncoder != nil {
						tf := ExtractTransferFromLog(log.Topics, log.Data, log.Address, uint(logIndex), txHash, number)
						if tf != nil {
							err := tfEncoder.Encode(tf)
							chk(err)
						}
					}

					if logEncoder != nil {
						err := logEncoder.Encode(NewCsvLog(number, txHash, uint(logIndex), log))
						chk(err)
					}

				}

				if internalTxEncoder != nil {
					for internalIndex, internalTx := range txInfo.InternalTransactions {
						for callInfoIndex, callInfo := range internalTx.CallValueInfo {
							err := internalTxEncoder.Encode(NewCsvInternalTx(number, txHash, uint(internalIndex), internalTx, uint(callInfoIndex), callInfo.TokenID, callInfo.CallValue))
							chk(err)
						}

					}
				}

			}

			log.Printf("parsed block %d", number)
		}

		return
	}

}
