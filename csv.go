package main

import (
	"encoding/hex"
	"math/big"
	"strings"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
	"golang.org/x/crypto/sha3"
)

// CsvTransaction represents a tron tx csv output, not trc10
// 1 TRX = 1000000 sun
type CsvTransaction struct {
	Hash             string `csv:"hash"`
	Nonce            string `csv:"nonce"`
	BlockHash        string `csv:"block_hash"`
	BlockNumber      uint64 `csv:"block_number"`
	TransactionIndex int    `csv:"transaction_index"`

	FromAddress          string `csv:"from_address"`
	ToAddress            string `csv:"to_address"`
	Value                string `csv:"value"`
	Gas                  string `csv:"gas"`
	GasPrice             string `csv:"gas_price"`
	Input                string `csv:"input"`
	BlockTimestamp       uint64 `csv:"block_timestamp"`
	MaxFeePerGas         string `csv:"max_fee_per_gas"`
	MaxPriorityFeePerGas string `csv:"max_priority_fee_per_gas"`
	TransactionType      string `csv:"transaction_type"`

	Status string

	// appendix
	TransactionTimestamp  uint64 `csv:"transaction_timestamp"`
	TransactionExpiration uint64 `csv:"transaction_expiration"`
	FeeLimit              uint64 `csv:"fee_limit"`
}

// NewCsvTransaction creates a new CsvTransaction
func NewCsvTransaction(blockTimestamp uint64, txIndex int, jsontx *tron.JSONTransaction, httptx *tron.HTTPTransaction) *CsvTransaction {
	to := ""
	if jsontx.To != "" {
		to = tron.Hex2TAddr(jsontx.To[2:])
	}

	txType := "Unknown"
	if len(httptx.RawData.Contract) > 0 {
		txType = httptx.RawData.Contract[0].ContractType
	}

	status := ""
	if len(httptx.Ret) > 0 {
		status = httptx.Ret[0].ContractRet
	}

	return &CsvTransaction{
		Hash:                 jsontx.Hash[2:],
		Nonce:                "", //tx.Nonce,
		BlockHash:            jsontx.BlockHash[2:],
		BlockNumber:          uint64(*jsontx.BlockNumber),
		TransactionIndex:     txIndex,
		FromAddress:          tron.Hex2TAddr(jsontx.From[2:]),
		ToAddress:            to,
		Value:                jsontx.Value.ToInt().String(),
		Gas:                  jsontx.Gas.ToInt().String(),
		GasPrice:             jsontx.GasPrice.ToInt().String(), // https://support.ledger.com/hc/en-us/articles/6331588714141-How-do-Tron-TRX-fees-work-?support=true
		Input:                jsontx.Input[2:],
		BlockTimestamp:       blockTimestamp / 1000, // unit: sec
		MaxFeePerGas:         "",                    //tx.MaxFeePerGas.String(),
		MaxPriorityFeePerGas: "",                    //tx.MaxPriorityFeePerGas.String(),
		TransactionType:      txType,                //jsontx.Type[2:],

		Status: status, // can be SUCCESS REVERT

		// appendix
		TransactionTimestamp:  httptx.RawData.Timestamp / 1000,  // float64(httptx.RawData.Timestamp) * 1 / 1000,
		TransactionExpiration: httptx.RawData.Expiration / 1000, // float64(httptx.RawData.Expiration) * 1 / 1000,
		FeeLimit:              httptx.RawData.FeeLimit,
	}
}

// CsvBlock represents a tron block output
type CsvBlock struct {
	Number           uint64 `csv:"number"`
	Hash             string `csv:"hash"`
	ParentHash       string `csv:"parent_hash"`
	Nonce            string `csv:"nonce"`
	Sha3Uncles       string `csv:"sha3_uncles"`
	LogsBloom        string `csv:"logs_bloom"`
	TransactionsRoot string `csv:"transaction_root"`
	StateRoot        string `csv:"state_root"`
	ReceiptsRoot     string `csv:"receipts_root"`
	Miner            string `csv:"miner"`
	Difficulty       string `csv:"difficulty"`
	TotalDifficulty  string `csv:"total_difficulty"`
	Size             uint64 `csv:"size"`
	ExtraData        string `csv:"extra_data"`
	GasLimit         string `csv:"gas_limit"`
	GasUsed          string `csv:"gas_used"`
	Timestamp        uint64 `csv:"timestamp"`
	TansactionCount  int    `csv:"transaction_count"`
	BaseFeePerGas    string `csv:"base_fee_per_gas"`

	// append
	WitnessSignature string `csv:"witness_signature"`
}

// NewCsvBlock creates a new CsvBlock
func NewCsvBlock(jsonblock *tron.JSONBlockWithTxs, httpblock *tron.HTTPBlock) *CsvBlock {
	return &CsvBlock{
		Number:           uint64(*jsonblock.Number),
		Hash:             jsonblock.Hash[2:],
		ParentHash:       jsonblock.ParentHash[2:],
		Nonce:            "",
		Sha3Uncles:       "", // block.Sha3Uncles,
		LogsBloom:        jsonblock.LogsBloom[2:],
		TransactionsRoot: jsonblock.TransactionsRoot[2:],
		StateRoot:        jsonblock.StateRoot[2:],
		ReceiptsRoot:     "",                                  // block.ReceiptsRoot
		Miner:            tron.Hex2TAddr(jsonblock.Miner[2:]), // = WitnessAddress
		Difficulty:       "",
		TotalDifficulty:  "",
		Size:             uint64(*jsonblock.Size),
		ExtraData:        "",
		GasLimit:         jsonblock.GasLimit.ToInt().String(),
		GasUsed:          jsonblock.GasUsed.ToInt().String(),
		Timestamp:        uint64(*jsonblock.Timestamp) / 1000,
		TansactionCount:  len(jsonblock.Transactions),
		BaseFeePerGas:    "", // block.BaseFeePerGas,

		//append
		WitnessSignature: httpblock.BlockHeader.WitnessSignature,
	}
}

// CsvTRC10Transfer is a trc10 transfer output
// https://developers.tron.network/docs/trc10-transfer-in-smart-contracts
// https://tronprotocol.github.io/documentation-en/mechanism-algorithm/system-contracts/
// It represents:
// - TransferContract
// - TransferAssetContract
type CsvTRC10Transfer struct {
	BlockNumber       uint64 `csv:"block_number"`
	BlockHash         string `csv:"block_hash"`
	TransactionHash   string `csv:"transaction_hash"`
	TransactionIndex  int    `csv:"transaction_index"`
	ContractCallIndex int    `csv:"contract_call_index"`

	AssetName   string `csv:"asset_name"` // do not omit => empty means trx
	FromAddress string `csv:"from_address"`
	ToAddress   string `csv:"to_address"`
	Value       string `csv:"value"`
}

// NewCsvTRC10Transfer creates a new CsvTRC10Transfer
func NewCsvTRC10Transfer(blockNum uint64, txIndex, callIndex int, httpTx *tron.HTTPTransaction, tfParams *tron.TRC10TransferParams) *CsvTRC10Transfer {

	return &CsvTRC10Transfer{
		TransactionHash:   httpTx.TxID,
		BlockHash:         httpTx.RawData.RefBlockHash,
		BlockNumber:       blockNum,
		TransactionIndex:  txIndex,
		ContractCallIndex: callIndex,

		AssetName:   tfParams.AssetName,
		FromAddress: tron.Hex2TAddr(tfParams.OwnerAddress),
		ToAddress:   tron.Hex2TAddr(tfParams.ToAddress),
		Value:       tfParams.Amount.String(),
	}
}

// CsvLog is a EVM smart contract event log output
type CsvLog struct {
	BlockNumber     uint64 `csv:"block_number"`
	TransactionHash string `csv:"transaction_hash"`
	LogIndex        uint   `csv:"log_index"`

	Address string `csv:"address"`
	Topics  string `csv:"topics"`
	Data    string `csv:"data"`
}

// NewCsvLog creates a new CsvLog
func NewCsvLog(blockNumber uint64, txHash string, logIndex uint, log *tron.HTTPTxInfoLog) *CsvLog {
	return &CsvLog{
		BlockNumber:     blockNumber,
		TransactionHash: txHash,
		LogIndex:        logIndex,

		Address: tron.Hex2TAddr(log.Address),
		Topics:  strings.Join(log.Topics, ";"),
		Data:    log.Data,
	}
}

// CsvInternalTx is a EVM smart contract internal transaction
type CsvInternalTx struct {
	TransactionHash   string `csv:"transaction_hash"`
	Index             uint   `csv:"internal_index"`
	CallerAddress     string `csv:"caller_address"`
	TransferToAddress string `csv:"transferTo_address"`
	CallInfoIndex     uint   `csv:"call_info_index"`
	CallTokenID       string `csv:"call_token_id"`
	CallValue         int64  `csv:"call_value"`
	Note              string `csv:"note"`
	Rejected          bool   `csv:"rejected"`
}

// NewCsvInternalTx creates a new CsvInternalTx
func NewCsvInternalTx(index uint, itx *tron.HTTPInternalTransaction, callInfoIndex uint, tokenID string, value int64) *CsvInternalTx {

	return &CsvInternalTx{
		TransactionHash:   itx.TransactionHash,
		Index:             index,
		CallerAddress:     tron.Hex2TAddr(itx.CallerAddress),
		TransferToAddress: tron.Hex2TAddr(itx.TransferToAddress),
		// CallValueInfo:     strings.Join(callValues, ";"),
		CallInfoIndex: callInfoIndex,
		CallTokenID:   tokenID,
		CallValue:     value,

		Note:     itx.Note,
		Rejected: itx.Rejected,
	}
}

// CsvReceipt is a receipt for tron transaction
type CsvReceipt struct {
	TxHash  string `csv:"transaction_hash"`
	TxIndex uint   `csv:"transaction_index"`
	// BlockHash         string `csv:"block_hash"` // cannot get this
	BlockNumber       uint64 `csv:"block_number"`
	ContractAddress   string `csv:"contract_address"`
	EnergyUsage       int64  `csv:"energy_usage,omitempty"`
	EnergyFee         int64  `csv:"energy_fee,omitempty"`
	OriginEnergyUsage int64  `csv:"origin_energy_usage,omitempty"`
	EnergyUsageTotal  int64  `csv:"energy_usage_total,omitempty"`
	NetUsage          int64  `csv:"net_usage,omitempty"`
	NetFee            int64  `csv:"net_fee,omitempty"`
	Result            string `csv:"result"`
}

func NewCsvReceipt(blockNum uint64, txHash string, txIndex uint, contractAddr string, r *tron.HTTPReceipt) *CsvReceipt {

	return &CsvReceipt{
		TxHash:  txHash,
		TxIndex: txIndex,
		// BlockHash:         blockHash,
		BlockNumber:       blockNum,
		ContractAddress:   contractAddr,
		EnergyUsage:       r.EnergyUsage,
		EnergyFee:         r.EnergyFee,
		OriginEnergyUsage: r.OriginEnergyUsage,
		EnergyUsageTotal:  r.EnergyUsageTotal,
		NetUsage:          r.NetUsage,
		NetFee:            r.NetFee,
		Result:            r.Result,
	}
}

// CsvAccount is a tron account
type CsvAccount struct {
	AccountName string `csv:"account_name"`
	Address     string `csv:"address"`
	Type        string `csv:"type"`
	CreateTime  int64  `csv:"create_time"`

	// DecodedName string `csv:decoded_name`
}

func NewCsvAccount(acc *tron.HTTPAccount) *CsvAccount {
	name, _ := hex.DecodeString(acc.AccountName)
	return &CsvAccount{
		AccountName: string(name),
		Address:     tron.Hex2TAddr(acc.Address),
		Type:        acc.AccountType,
		CreateTime:  acc.CreateTime / 1000,
	}
}

// CsvContract is a standard EVM contract
type CsvContract struct {
	Address           string `csv:"address"`
	Bytecode          string `csv:"bytecode"`
	FunctionSighashes string `csv:"function_sighashes"`
	IsErc20           bool   `csv:"is_erc20"`
	IsErc721          bool   `csv:"is_erc721"`
	BlockNumber       uint64 `csv:"block_number"`

	// append some...
	ContractName               string
	ConsumeUserResourcePercent int
	OriginAddress              string
	OriginEnergyLimit          int64
}

var keccakHasher = sha3.NewLegacyKeccak256()

func NewCsvContract(c *tron.HTTPContract) *CsvContract {
	hashes := make([]string, 0, len(c.Abi.Entrys))
	for _, abi := range c.Abi.Entrys {
		if strings.ToLower(abi.Type) == "function" {
			content := abi.Name + "("
			types := make([]string, 0, len(abi.Inputs))
			for _, input := range abi.Inputs {
				types = append(types, input.Type)
			}
			funchash := keccakHasher.Sum([]byte(content + strings.Join(types, ",") + ")"))
			hashes = append(hashes, hex.EncodeToString(funchash))
		}
	}

	isErc20 := implementsAnyOf(hashes, "totalSupply()") &&
		implementsAnyOf(hashes, "balanceOf(address)") &&
		implementsAnyOf(hashes, "transfer(address,uint256)") &&
		implementsAnyOf(hashes, "transferFrom(address,address,uint256)") &&
		implementsAnyOf(hashes, "approve(address,uint256)") &&
		implementsAnyOf(hashes, "allowance(address,address)")

	isErc721 := implementsAnyOf(hashes, "balanceOf(address)") &&
		implementsAnyOf(hashes, "ownerOf(uint256)") &&
		implementsAnyOf(hashes, "transfer(address,uint256)", "transferFrom(address,address,uint256)") &&
		implementsAnyOf(hashes, "approve(address,uint256)")

	return &CsvContract{
		Address:           tron.Hex2TAddr(c.ContractAddress),
		Bytecode:          c.Bytecode,
		FunctionSighashes: strings.Join(hashes, ";"),
		IsErc20:           isErc20,
		IsErc721:          isErc721,
		BlockNumber:       0,

		// append
		ContractName:               c.Name,
		ConsumeUserResourcePercent: c.ConsumeUserResourcePercent,
		OriginAddress:              tron.Hex2TAddr(c.OriginAddress),
		OriginEnergyLimit:          c.OriginEnergyLimit,
	}
}

func implementsAnyOf(hashes []string, sigStrs ...string) bool {
	for i := range sigStrs {
		hash := hex.EncodeToString(keccakHasher.Sum([]byte(sigStrs[i])))
		for j := range hashes {
			if hashes[j] == hash {
				return true
			}
		}
	}

	return false
}

// CsvTokens is a standard EVM contract token
type CsvTokens struct {
	Address     string `csv:"address"`
	Symbol      string `csv:"symbol"`
	Name        string `csv:"name"`
	Decimals    uint64 `csv:"decimals"`
	TotalSupply uint64 `csv:"total_supply"`
	BlockNumber uint64 `csv:"block_number"`
}

func NewCsvTokens(cli *tron.TronClient, contract *tron.HTTPContract) *CsvTokens {
	contractAddr := contract.ContractAddress
	callerAddr := contract.OriginAddress

	symbolResult := cli.CallContract(contractAddr, callerAddr, 0, 1000,
		"symbol()",
	)
	symbol := ParseSymbol(symbolResult.ConstantResult)

	nameResult := cli.CallContract(contractAddr, callerAddr, 0, 1000,
		"name()",
	)
	name := ParseName(nameResult.ConstantResult)

	decimalsResult := cli.CallContract(contractAddr, callerAddr, 0, 1000,
		"decimals()",
	)
	decimals := ParseDecimals(decimalsResult.ConstantResult)

	totalSupplyResult := cli.CallContract(contractAddr, callerAddr, 0, 1000,
		"totalSupply()",
	)
	totalSupply := ParseTotalSupply(totalSupplyResult.ConstantResult)

	block := cli.GetJSONBlockByNumberWithTxIDs(nil)

	return &CsvTokens{
		Address:     contractAddr,
		Symbol:      symbol,
		Name:        name,
		Decimals:    *decimals,
		TotalSupply: *totalSupply,
		BlockNumber: uint64(*block.Number),
	}
}

func ParseSymbol(contractResults []string) string {
	if len(contractResults) == 0 {
		panic("failed to parse symbol")
	}

	result := contractResults[0]
	bigLlen, ok := new(big.Int).SetString(result[0:64], 16)
	if !ok {
		// TODO: warn log here
		return ""
	}
	l, ok := new(big.Int).SetString(result[64:64+bigLlen.Int64()*2], 16)
	if !ok {
		// TODO: warn log here
		return ""
	}
	hexStr := result[64+bigLlen.Int64()*2 : 64+bigLlen.Int64()*2+l.Int64()*2]
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		// TODO: err log here
		return ""
	}
	return string(decoded)
}

func ParseName(contractResults []string) string {
	if len(contractResults) == 0 {
		panic("failed to parse symbol")
	}

	result := contractResults[0]
	bigLlen, ok := new(big.Int).SetString(result[0:64], 16)
	if !ok {
		// TODO: warn log here
		return ""
	}
	l, ok := new(big.Int).SetString(result[64:64+bigLlen.Int64()*2], 16)
	if !ok {
		// TODO: warn log here
		return ""
	}
	hexStr := result[64+bigLlen.Int64()*2 : 64+bigLlen.Int64()*2+l.Int64()*2]
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		// TODO: err log here
		return ""
	}
	return string(decoded)
}

func ParseDecimals(contractResults []string) *uint64 {
	if len(contractResults) == 0 {
		panic("failed to parse symbol")
	}

	result, ok := new(big.Int).SetString(contractResults[0], 16)
	if !ok {
		// TODO: err log here
		return nil
	}

	rtn := result.Uint64()
	return &rtn
}

func ParseTotalSupply(contractResults []string) *uint64 {
	if len(contractResults) == 0 {
		panic("failed to parse symbol")
	}

	result, ok := new(big.Int).SetString(contractResults[0], 16)
	if !ok {
		// TODO: err log here
		return nil
	}

	rtn := result.Uint64()
	return &rtn
}
