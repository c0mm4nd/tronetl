package main

import (
	"bytes"
	"encoding/csv"
	"os"
	"testing"
)

// Integration tests that hit a real TRON fullnode (set TRON_PROVIDER).
// They verify both correctness of key fields and that addresses remain T-addr.

func requireProviderFunctional(t *testing.T) string {
	t.Helper()
	p := os.Getenv("TRON_PROVIDER")
	if p == "" {
		t.Skip("TRON_PROVIDER not set; skipping integration export tests")
	}
	return p
}

// small csv helpers
func csvRows(buf *bytes.Buffer) [][]string {
	rows, _ := csv.NewReader(bytes.NewReader(buf.Bytes())).ReadAll()
	return rows
}

func col(rows [][]string, name string) int {
	for i, h := range rows[0] {
		if h == name {
			return i
		}
	}
	return -1
}

func findRow(rows [][]string, pred func([]string) bool) []string {
	for i := 1; i < len(rows); i++ {
		if pred(rows[i]) {
			return rows[i]
		}
	}
	return nil
}

// Known fixtures for block 76311000 on mainnet (stable historical data)
const (
	block76311000       = 76311000
	txUSDTTransferHash  = "bcf03ecf35947b6a7a65af92a5d15a73d0e6d3ff5f89e4a9af545bf9196e1688"
	usdtTAddr           = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	usdtSenderTAddr     = "TBbjGN8V4TuA1b1A1WHRe4FYfxdkNcKCy8"
	usdtReceiverTAddr   = "TTSFgXBu5YMjDWZvQdMxsaqJhWNAKhbrey"
	internalTxHash      = "b233ecea785c680bab83a63a78038016fe3a2ebd593091e240b700d590be1172"
	internalCallerTAddr = "TDQaYrhQynYV9aXTYj63nwLAafRffWSEj6"
	internalToTAddr     = "TSSMHYeV2uE9qYH95DqyoCuNCzEL1NvU3S"
)

func TestExportBlocksAndTransactions_FunctionalValues(t *testing.T) {
	provider := requireProviderFunctional(t)

	blkBuf, txBuf, trc10Buf := &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
	opts := &ExportBlocksAndTransactionsOptions{
		ProviderURI: provider,
		StartBlock:  block76311000,
		EndBlock:    block76311000,
		blksOutput:  blkBuf,
		txsOutput:   txBuf,
		trc10Output: trc10Buf,
	}
	ExportBlocksAndTransactions(opts)

	txs := csvRows(txBuf)
	hashCol := col(txs, "hash")
	fromCol := col(txs, "from_address")
	toCol := col(txs, "to_address")
	typeCol := col(txs, "transaction_type")
	if hashCol == -1 || fromCol == -1 || toCol == -1 || typeCol == -1 {
		t.Fatalf("transaction columns missing")
	}
	row := findRow(txs, func(r []string) bool { return r[hashCol] == txUSDTTransferHash })
	if row == nil {
		t.Fatalf("known tx %s not found", txUSDTTransferHash)
	}
	if row[fromCol] != usdtSenderTAddr || row[toCol] != usdtTAddr {
		t.Fatalf("tx addresses mismatch: from=%s to=%s", row[fromCol], row[toCol])
	}
	if row[typeCol] != "TriggerSmartContract" {
		t.Fatalf("expected TriggerSmartContract, got %s", row[typeCol])
	}
}

func TestExportTokenTransfers_FunctionalValues(t *testing.T) {
	provider := requireProviderFunctional(t)

	tfBuf, logBuf, intBuf, recBuf := &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
	opts := &ExportTransferOptions{
		ProviderURI:     provider,
		StartBlock:      block76311000,
		EndBlock:        block76311000,
		tfOutput:        tfBuf,
		logOutput:       logBuf,
		internalTxOutput: intBuf,
		receiptOutput:    recBuf,
	}
	ExportTransfersWithWorkers(opts, 2)

	// token_transfers
	tfs := csvRows(tfBuf)
	thCol := col(tfs, "transaction_hash")
	tokenCol := col(tfs, "token_address")
	fromCol := col(tfs, "from_address")
	toCol := col(tfs, "to_address")
	valCol := col(tfs, "value")
	if thCol == -1 || tokenCol == -1 || fromCol == -1 || toCol == -1 || valCol == -1 {
		t.Fatalf("token transfer columns missing")
	}
	tfRow := findRow(tfs, func(r []string) bool { return r[thCol] == txUSDTTransferHash })
	if tfRow == nil {
		t.Fatalf("token transfer for tx %s not found", txUSDTTransferHash)
	}
	if tfRow[tokenCol] != usdtTAddr || tfRow[fromCol] != usdtSenderTAddr || tfRow[toCol] != usdtReceiverTAddr {
		t.Fatalf("token transfer addresses mismatch: token=%s from=%s to=%s", tfRow[tokenCol], tfRow[fromCol], tfRow[toCol])
	}
	if tfRow[valCol] != "11000000" {
		t.Fatalf("unexpected transfer value %s", tfRow[valCol])
	}

	// receipts
	recs := csvRows(recBuf)
	rcCol := col(recs, "transaction_hash")
	ctrCol := col(recs, "contract_address")
	resCol := col(recs, "result")
	if rcCol == -1 || ctrCol == -1 || resCol == -1 {
		t.Fatalf("receipt columns missing")
	}
	recRow := findRow(recs, func(r []string) bool { return r[rcCol] == txUSDTTransferHash })
	if recRow == nil {
		t.Fatalf("receipt for tx %s not found", txUSDTTransferHash)
	}
	if recRow[ctrCol] != usdtTAddr {
		t.Fatalf("receipt contract address mismatch: %s", recRow[ctrCol])
	}
	if recRow[resCol] != "SUCCESS" {
		t.Fatalf("receipt result unexpected: %s", recRow[resCol])
	}

	// logs
	logs := csvRows(logBuf)
	lhCol := col(logs, "transaction_hash")
	addrCol := col(logs, "address")
	if lhCol == -1 || addrCol == -1 {
		t.Fatalf("log columns missing")
	}
	logRow := findRow(logs, func(r []string) bool { return r[lhCol] == txUSDTTransferHash })
	if logRow == nil || logRow[addrCol] != usdtTAddr {
		t.Fatalf("log for tx %s not found or addr mismatch", txUSDTTransferHash)
	}

	// internal transactions
	itxs := csvRows(intBuf)
	ithCol := col(itxs, "transaction_hash")
	callerCol := col(itxs, "caller_address")
	toIntCol := col(itxs, "transferTo_address")
	if ithCol == -1 || callerCol == -1 || toIntCol == -1 {
		t.Fatalf("internal tx columns missing")
	}
	itRow := findRow(itxs, func(r []string) bool { return r[ithCol] == internalTxHash })
	if itRow == nil {
		t.Fatalf("internal tx %s not found", internalTxHash)
	}
	if itRow[callerCol] != internalCallerTAddr || itRow[toIntCol] != internalToTAddr {
		t.Fatalf("internal tx addresses mismatch: caller=%s to=%s", itRow[callerCol], itRow[toIntCol])
	}
}
