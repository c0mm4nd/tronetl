package main

import (
	"log"
	"math/big"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
	"github.com/jszwec/csvutil"
)

// locateStartBlock is a util for locating the start block by start timestamp
func locateStartBlock(cli *tron.TronClient, startTimestamp uint64) uint64 {
	latestBlock := cli.GetJSONBlockByNumberWithTxIDs(nil)
	top := latestBlock.Number
	half := uint64(*top) / 2
	estimateStartNumber := half
	for {
		block := cli.GetJSONBlockByNumberWithTxIDs(new(big.Int).SetUint64(estimateStartNumber))
		if block == nil {
			break
		}
		log.Println(half, block.Timestamp)
		timestamp := uint64(*block.Timestamp / 1000)
		if timestamp < startTimestamp && startTimestamp-timestamp < 60 {
			break
		}

		//
		if timestamp < startTimestamp {
			log.Printf("%d is too small: %d", estimateStartNumber, timestamp)
			half = half / 2
			estimateStartNumber = estimateStartNumber + half
		} else {
			log.Printf("%d is too large: %d", estimateStartNumber, timestamp)
			half = half / 2
			estimateStartNumber = estimateStartNumber - half
		}

		if half == 0 || estimateStartNumber >= uint64(*top) {
			panic("failed to find the block on that timestamp")
		}
	}

	return estimateStartNumber
}

// locateEndBlock is a util for locating the end block by end timestamp
func locateEndBlock(cli *tron.TronClient, endTimestamp uint64) uint64 {
	latestBlock := cli.GetJSONBlockByNumberWithTxIDs(nil)
	top := latestBlock.Number
	half := uint64(*top) / 2
	estimateEndNumber := half
	for {
		block := cli.GetJSONBlockByNumberWithTxIDs(new(big.Int).SetUint64(estimateEndNumber))
		if block == nil {
			break
		}
		log.Println(half, block.Timestamp)
		timestamp := uint64(*block.Timestamp / 1000)
		if timestamp > endTimestamp && timestamp-endTimestamp < 60 {
			break
		}

		//
		if timestamp < endTimestamp {
			log.Printf("%d is too small: %d", estimateEndNumber, timestamp)
			half = half / 2
			estimateEndNumber = estimateEndNumber + half
		} else {
			log.Printf("%d is too large: %d", estimateEndNumber, timestamp)
			half = half / 2
			estimateEndNumber = estimateEndNumber - half
		}

		if half == 0 || estimateEndNumber >= uint64(*top) {
			panic("failed to find the block on that timestamp")
		}
	}

	return estimateEndNumber
}

func createCSVEncodeCh(enc *csvutil.Encoder, maxWorker uint) chan any {
	ch := make(chan any, maxWorker)
	writeFn := func() {
		for {
			obj := <-ch
			err := enc.Encode(obj)
			chk(err)
		}
	}

	go writeFn()
	return ch
}
