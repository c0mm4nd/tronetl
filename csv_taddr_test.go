package main

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
)

// helper constants: TR7NH... (USDT) known mapping
const (
	tAddr = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	// hex with 0x prefix for JSON fields
	hexWith0x = "0x41a614f803b6fd780986a42c78ec9c7f77e6ded13c"
	// raw hex without prefix for HTTP/contract fields
	rawHex = "41a614f803b6fd780986a42c78ec9c7f77e6ded13c"
)

func TestEnsureTAddr_IdempotentAndHexConversion(t *testing.T) {
	if tron.EnsureTAddr(tAddr) != tAddr {
		t.Fatalf("EnsureTAddr should keep T-addr unchanged")
	}
	if got := tron.EnsureTAddr(rawHex); got != tAddr {
		t.Fatalf("EnsureTAddr should convert hex to T-addr, got %s", got)
	}
}

func TestNewCsvTransaction_UsesTAddr(t *testing.T) {
	bn := hexutil.Uint64(1)
	val := hexutil.Big(*big.NewInt(1))
	gas := hexutil.Big(*big.NewInt(1))
	gasPrice := hexutil.Big(*big.NewInt(1))
	js := &tron.JSONTransaction{
		BlockHash:        "0xdead",
		BlockNumber:      &bn,
		From:             hexWith0x,
		To:               hexWith0x,
		Gas:              &gas,
		GasPrice:         &gasPrice,
		Hash:             "0xhash",
		Input:            "0x",
		TransactionIndex: &bn,
		Type:             "0x0",
		Value:            &val,
	}
	httpTx := &tron.HTTPTransaction{}

	c := NewCsvTransaction(0, 0, js, httpTx)
	if c.FromAddress != tAddr || c.ToAddress != tAddr {
		t.Fatalf("addresses should be T-addr, got from=%s to=%s", c.FromAddress, c.ToAddress)
	}
}

func TestNewCsvTRC10Transfer_UsesTAddr(t *testing.T) {
	bn := hexutil.Uint64(1)
	val := hexutil.Big(*big.NewInt(1))
	gas := hexutil.Big(*big.NewInt(1))
	gasPrice := hexutil.Big(*big.NewInt(1))
	js := &tron.JSONTransaction{
		BlockHash:        "0xdead",
		BlockNumber:      &bn,
		From:             hexWith0x,
		To:               hexWith0x,
		Gas:              &gas,
		GasPrice:         &gasPrice,
		Hash:             "0xhash",
		Input:            "0x",
		TransactionIndex: &bn,
		Type:             "0x0",
		Value:            &val,
	}
	tf := &tron.TRC10TransferParams{
		OwnerAddress: rawHex,
		ToAddress:    rawHex,
		Amount:       big.NewInt(1),
		AssetName:    "TRX",
	}
	httpTx := &tron.HTTPTransaction{TxID: "tx"}
	c := NewCsvTRC10Transfer("blockhash", 1, 0, 0, js, httpTx, tf)
	if c.FromAddress != tAddr || c.ToAddress != tAddr {
		t.Fatalf("TRC10 transfer addresses should be T-addr, got %s %s", c.FromAddress, c.ToAddress)
	}
}

func TestNewCsvLog_UsesTAddr(t *testing.T) {
	log := &tron.HTTPTxInfoLog{Address: rawHex}
	c := NewCsvLog(1, "tx", 0, log)
	if c.Address != tAddr {
		t.Fatalf("log address should be T-addr, got %s", c.Address)
	}
}

func TestNewCsvInternalTx_UsesTAddr(t *testing.T) {
	itx := &tron.HTTPInternalTransaction{
		InternalTransactionHash: "hash",
		CallerAddress:           rawHex,
		TransferToAddress:       rawHex,
	}
	c := NewCsvInternalTx(1, "tx", 0, itx, 0, "", 0)
	if c.CallerAddress != tAddr || c.TransferToAddress != tAddr {
		t.Fatalf("internal tx addresses should be T-addr, got %s %s", c.CallerAddress, c.TransferToAddress)
	}
}

func TestNewCsvReceipt_UsesTAddr(t *testing.T) {
	r := &tron.HTTPReceipt{}
	c := NewCsvReceipt(1, "tx", 0, rawHex, r)
	if c.ContractAddress != tAddr {
		t.Fatalf("receipt contract address should be T-addr, got %s", c.ContractAddress)
	}
}

func TestNewCsvAccount_UsesTAddr(t *testing.T) {
	acc := &tron.HTTPAccount{Address: rawHex}
	c := NewCsvAccount(acc)
	if c.Address != tAddr {
		t.Fatalf("account address should be T-addr, got %s", c.Address)
	}
}

func TestNewCsvContract_UsesTAddr(t *testing.T) {
	ctr := &tron.HTTPContract{
		ContractAddress: rawHex,
		OriginAddress:   rawHex,
	}
	c := NewCsvContract(ctr)
	if c.Address != tAddr || c.OriginAddress != tAddr {
		t.Fatalf("contract addresses should be T-addr, got %s %s", c.Address, c.OriginAddress)
	}
}

func TestNewCsvReceipt_AllowsEmptyContractAddr(t *testing.T) {
	r := &tron.HTTPReceipt{}
	c := NewCsvReceipt(1, "tx", 0, "", r)
	if c.ContractAddress != "" {
		t.Fatalf("empty contract address should stay empty, got %s", c.ContractAddress)
	}
}
