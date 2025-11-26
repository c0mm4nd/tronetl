package main

import (
	"bytes"
	"encoding/csv"
	"os"
	"strings"
	"testing"
)

func requireProvider(t *testing.T) string {
	t.Helper()
	provider := os.Getenv("TRON_PROVIDER")
	if provider == "" {
		t.Skip("TRON_PROVIDER not set; skip integration export tests")
	}
	return provider
}

// check all non-empty entries in column start with 'T'
func assertColumnTAddr(t *testing.T, rows [][]string, header string) {
	t.Helper()
	if len(rows) == 0 {
		t.Fatalf("no rows to inspect for column %s", header)
	}
	idx := -1
	for i, h := range rows[0] {
		if h == header {
			idx = i
			break
		}
	}
	if idx == -1 {
		t.Fatalf("column %s not found", header)
	}
	for i := 1; i < len(rows); i++ {
		val := rows[i][idx]
		if val == "" {
			continue
		}
		if !strings.HasPrefix(val, "T") {
			t.Fatalf("row %d column %s not T-addr: %s", i, header, val)
		}
	}
}

// Ensure export_blocks_and_transactions outputs T-addr fields.
func TestExportBlocksAndTransactionsOutputsTAddr(t *testing.T) {
	provider := requireProvider(t)

	blockBuf, txBuf, trc10Buf := &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
	opts := &ExportBlocksAndTransactionsOptions{
		ProviderURI: provider,
		StartBlock:  76311000,
		EndBlock:    76311000,
		blksOutput:  blockBuf,
		txsOutput:   txBuf,
		trc10Output: trc10Buf,
	}

	ExportBlocksAndTransactions(opts)

	blocks, _ := csv.NewReader(bytes.NewReader(blockBuf.Bytes())).ReadAll()
	txs, _ := csv.NewReader(bytes.NewReader(txBuf.Bytes())).ReadAll()
	trc10, _ := csv.NewReader(bytes.NewReader(trc10Buf.Bytes())).ReadAll()

	assertColumnTAddr(t, blocks, "miner")

	// transactions: from / to may be empty for contract creation
	assertColumnTAddr(t, txs, "from_address")
	assertColumnTAddr(t, txs, "to_address")

	// trc10 transfers: from / to should be T when present
	if len(trc10) > 0 { // sometimes none in the block
		assertColumnTAddr(t, trc10, "from_address")
		assertColumnTAddr(t, trc10, "to_address")
	}
}

// Ensure export_token_transfers outputs T-addr fields.
func TestExportTokenTransfersOutputsTAddr(t *testing.T) {
	provider := requireProvider(t)

	tfBuf, logBuf, internalBuf, receiptBuf := &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
	opts := &ExportTransferOptions{
		ProviderURI: provider,
		StartBlock:  76311000,
		EndBlock:    76311000,
		tfOutput:    tfBuf,
		logOutput:   logBuf,
		internalTxOutput: internalBuf,
		receiptOutput:    receiptBuf,
	}

	ExportTransfersWithWorkers(opts, 2)

	tfs, _ := csv.NewReader(bytes.NewReader(tfBuf.Bytes())).ReadAll()
	logs, _ := csv.NewReader(bytes.NewReader(logBuf.Bytes())).ReadAll()
	itxs, _ := csv.NewReader(bytes.NewReader(internalBuf.Bytes())).ReadAll()
	receipts, _ := csv.NewReader(bytes.NewReader(receiptBuf.Bytes())).ReadAll()

	if len(tfs) > 0 {
		assertColumnTAddr(t, tfs, "token_address")
		assertColumnTAddr(t, tfs, "from_address")
		assertColumnTAddr(t, tfs, "to_address")
	}
	if len(logs) > 0 {
		assertColumnTAddr(t, logs, "address")
	}
	if len(itxs) > 0 {
		assertColumnTAddr(t, itxs, "caller_address")
		assertColumnTAddr(t, itxs, "transferTo_address")
	}
	if len(receipts) > 0 {
		assertColumnTAddr(t, receipts, "contract_address")
	}
}
