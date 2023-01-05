package tron

import (
	"encoding/json"
	"math/big"
)

// HTTPTxInfo is a TxInfo result from HTTP RESTful API
// the struct follows https://github.com/tronprotocol/protocol/blob/2351aa6c2d708bf5ef47baf70410b3bc87d65fa7/core/Tron.proto#L341
type HTTPTxInfo struct {
	ID                            string                     `json:"id,omitempty"`
	Fee                           int                        `json:"fee,omitempty"`
	BlockNumber                   int                        `json:"blockNumber,omitempty"`
	BlockTimeStamp                int64                      `json:"blockTimeStamp,omitempty"`
	ContractResult                []string                   `json:"contractResult,omitempty"`
	ContractAddress               string                     `json:"contract_address,omitempty"`
	Receipt                       *HTTPReceipt               `json:"receipt,omitempty"`
	Log                           []*HTTPTxInfoLog           `json:"log,omitempty"`
	Result                        any                        `json:"result,omitempty"` // enum code { SUCESS = 0; FAILED = 1; }
	ResMessage                    string                     `json:"resMessage,omitempty"`
	AssetIssueID                  string                     `json:"assetIssueID,omitempty"`
	WithdrawAmount                int64                      `json:"withdraw_amount,omitempty"`
	UnfreezeAmount                int64                      `json:"unfreeze_amount,omitempty"`
	InternalTransactions          []*HTTPInternalTransaction `json:"internal_transactions,omitempty"`
	ExchangeReceivedAmount        int64                      `json:"exchange_received_amount,omitempty"`
	ExchangeInjectAnotherAmount   int64                      `json:"exchange_inject_another_amount,omitempty"`
	ExchangeWithdrawAnotherAmount int64                      `json:"exchange_withdraw_another_amount,omitempty"`
	ExchangeID                    int64                      `json:"exchange_id,omitempty"`
	ShieldedTransactionFee        int64                      `json:"shielded_transaction_fee,omitempty"`
}

type HTTPReceipt struct {
	EnergyUsage       int64  `json:"energy_usage,omitempty"`
	EnergyFee         int64  `json:"energy_fee,omitempty"`
	OriginEnergyUsage int64  `json:"origin_energy_usage,omitempty"`
	EnergyUsageTotal  int64  `json:"energy_usage_total,omitempty"`
	NetUsage          int64  `json:"net_usage,omitempty"`
	NetFee            int64  `json:"net_fee,omitempty"`
	Result            string `json:"result,omitempty"`
}

// HTTPTxInfoLog is a Log result from HTTP RESTful API
type HTTPTxInfoLog struct {
	Address string   `json:"address,omitempty"`
	Topics  []string `json:"topics,omitempty"`
	Data    string   `json:"data,omitempty"`
}

// HTTPInternalTransaction is a internal transaction result from HTTP RESTful API
type HTTPInternalTransaction struct {
	InternalTransactionHash string                                  `json:"hash,omitempty"`
	CallerAddress           string                                  `json:"caller_address,omitempty"`
	TransferToAddress       string                                  `json:"transferTo_address,omitempty"`
	CallValueInfo           []*HTTPInternalTransactionCallValueInfo `json:"callValueInfo,omitempty"`
	Note                    string                                  `json:"note,omitempty"`
	Rejected                bool                                    `json:"rejected,omitempty"`
}

// HTTPInternalTransactionCallValueInfo is a field in HTTPInternalTransaction
// https://github.com/tronprotocol/protocol/blob/2351aa6c2d708bf5ef47baf70410b3bc87d65fa7/core/Tron.proto#L509
type HTTPInternalTransactionCallValueInfo struct {
	CallValue int64  `json:"callValue,omitempty"`
	TokenID   string `json:"tokenId,omitempty"`
}

// HTTPInternalTransaction is a Block result from HTTP RESTful API
// https://github.com/tronprotocol/protocol/blob/2351aa6c2d708bf5ef47baf70410b3bc87d65fa7/core/Tron.proto#L406
type HTTPBlock struct {
	BlockID      string            `json:"blockID,omitempty"`
	BlockHeader  *HTTPBlockHeader  `json:"block_header,omitempty"`
	Transactions []HTTPTransaction `json:"transactions,omitempty"`
}

// HTTPBlockHeader represents the block header from the Block result from HTTP RESTful API
type HTTPBlockHeader struct {
	RawData struct {
		Timestamp  int64  `json:"timestamp,omitempty"`
		TxTrieRoot string `json:"txTrieRoot,omitempty"`
		ParentHash string `json:"parentHash,omitempty"`

		Number           int    `json:"number,omitempty"`
		WitnessAddress   string `json:"witness_address,omitempty"`
		Version          int    `json:"version,omitempty"`
		AccountStateRoot string `json:"accountStateRoot,omitempty"`
	} `json:"raw_data,omitempty"`
	WitnessSignature string `json:"witness_signature,omitempty"`
}

// HTTPTransaction represents the Transaction result from HTTP RESTful API
type HTTPTransaction struct {
	Ret []struct {
		ContractRet string `json:"contractRet,omitempty"`
	} `json:"ret,omitempty"`
	Signature []string `json:"signature,omitempty"`
	TxID      string   `json:"txID,omitempty"`
	RawData   struct {
		Data          string          `json:"data,omitempty"`
		Contract      []*ContractCall `json:"contract,omitempty"`
		RefBlockBytes string          `json:"ref_block_bytes,omitempty"`
		// - The height of the transaction reference block, using the 6th to 8th (exclusive) bytes of the reference block height, a total of 2 bytes.
		// The reference block is used in the TRON TAPOS mechanism, which can prevent a replay of a transaction on forks that do not include the referenced block.
		// Generally the latest solidified block is used as the reference block.
		RefBlockHash string `json:"ref_block_hash,omitempty"`
		// - The hash of the transaction reference block, using the 8th to 16th (exclusive) bytes of the reference block hash, a total of 8 bytes.
		// The reference block is used in the TRON TAPOS mechanism, which can prevent a replay of a transaction on forks that do not include the referenced block.
		// Generally the latest solidified block is used as the reference block.
		Expiration uint64 `json:"expiration,omitempty"`
		Timestamp  uint64 `json:"timestamp,omitempty"`
		FeeLimit   uint64 `json:"fee_limit,omitempty"`
	} `json:"raw_data,omitempty"`
	RawDataHex string `json:"raw_data_hex,omitempty"`
}

// ContractCall represents a tron native contract call inside the Transaction
// Details of Parameter in https://tronprotocol.github.io/documentation-en/mechanism-algorithm/system-contracts/
type ContractCall struct {
	ContractType string `json:"type,omitempty"`
	Parameter    struct {
		Value   json.RawMessage `json:"value,omitempty"`
		TypeURL string          `json:"type_url,omitempty"`
	} `json:"parameter,omitempty"` // google.any decode with ContractType
	Provider     string `json:"provider,omitempty"`
	PermissionID int32  `json:"Permission_id,omitempty"`
}

// TRC10TransferParams can be the params of the following calls:
// - TransferAssetContract
// - TransferContract
type TRC10TransferParams struct {
	AssetName    string   `json:"asset_name,omitempty"`
	Amount       *big.Int `json:"amount,omitempty"`
	OwnerAddress string   `json:"owner_address,omitempty"`
	ToAddress    string   `json:"to_address,omitempty"`
}

type HTTPVote struct {
	VoteAddress string `json:"vote_address,omitempty"`
	VoteCount   int64  `json:"vote_count,omitempty"`
}

type HTTPFrozen struct {
	FrozenBalance int64 `json:"frozen_balance,omitempty"`
	ExpireTime    int64 `json:"expire_time,omitempty"`
}

type HTTPPermision struct {
	PermType       string `json:"type,omitempty"`
	ID             int32  `json:"id,omitempty"`
	PermissionName string `json:"permission_name,omitempty"`
	Threshold      int64  `json:"threshold,omitempty"`
	Parent_id      int32  `json:"parent_id,omitempty"`
	Operations     string `json:"operations,omitempty"`
	Keys           []*struct {
		Address string `json:"address,omitempty"`
		Weight  int64  `json:"weight,omitempty"`
	} `json:"keys,omitempty"`
}

type HTTPAccount struct {
	AccountName string      `json:"account_name,omitempty"`
	AccountType string      `json:"type,omitempty"` // Normal = 0; AssetIssue = 1; Contract = 2;
	Address     string      `json:"address,omitempty"`
	Balance     int64       `json:"balance,omitempty"`
	Votes       []*HTTPVote `json:"votes,omitempty"`
	Asset       []struct {
		Key   string `json:"key,omitempty"`
		Value int    `json:"value,omitempty"`
	} `json:"asset,omitempty"`
	AssetV2 []struct {
		Key   string `json:"key,omitempty"`
		Value int    `json:"value,omitempty"`
	} `json:"assetV2,omitempty"`
	Frozen                                     []*HTTPFrozen `json:"frozen,omitempty"`
	NetUsage                                   int64         `json:"net_usage,omitempty"`
	AcquiredDelegatedFrozenBalanceForBandwidth int64         `json:"acquired_delegated_frozen_balance_for_bandwidth,omitempty"`
	DelegatedFrozenBalanceForBandwidth         int64         `json:"delegated_frozen_balance_for_bandwidth,omitempty"`
	CreateTime                                 int64         `json:"create_time,omitempty"`
	LatestOprationTime                         int64         `json:"latest_opration_time,omitempty"`
	Allowance                                  int64         `json:"allowance,omitempty"`
	// last withdraw time
	LatestWithdrawTime int64 `json:"latest_withdraw_time,omitempty"`
	// not used so far
	Code        string `json:"code,omitempty"`
	IsWitness   bool   `json:"is_witness,omitempty"`
	IsCommittee bool   `json:"is_committee,omitempty"`
	// frozen asset(for asset issuer)
	FrozenSupply []*HTTPFrozen `json:"frozen_supply,omitempty"`
	// asset_issued_name
	AssetIssuedName          string `json:"asset_issued_name,omitempty"`
	AssetIssuedID            string `json:"asset_issued_ID,omitempty"`
	LatestAssetOperationTime []struct {
		Key   string `json:"key,omitempty"`
		Value int    `json:"value,omitempty"`
	} `json:"latest_asset_operation_time,omitempty"`
	LatestAssetOperationTimeV2 []struct {
		Key   string `json:"key,omitempty"`
		Value int    `json:"value,omitempty"`
	} `json:"latest_asset_operation_timeV2,omitempty"`
	FreeNetUsage      int64 `json:"free_net_usage,omitempty"`
	FreeAssetNetUsage []struct {
		Key   string `json:"key,omitempty"`
		Value int    `json:"value,omitempty"`
	} `json:"free_asset_net_usage,omitempty"`
	FreeAssetNetUsageV2 []struct {
		Key   string `json:"key,omitempty"`
		Value int    `json:"value,omitempty"`
	} `json:"free_asset_net_usageV2,omitempty"`
	LatestConsumeTime     int64  `json:"latest_consume_time,omitempty"`
	LatestConsumeFreeTime int64  `json:"latest_consume_free_time,omitempty"`
	AccountID             string `json:"account_id,omitempty"`
	AccountResource       *struct {
		// energy resource, get from frozen
		EnergyUsage int64 `json:"energy_usage,omitempty"`
		// the frozen balance for energy
		FrozenBalanceForEnergy     *HTTPFrozen `json:"frozen_balance_for_energy,omitempty"`
		LatestConsumeTimeForEnergy int64       `json:"latest_consume_time_for_energy,omitempty"`

		// Frozen balance provided by other accounts to this account
		AcquiredDelegatedFrozenBalanceForEnergy int64 `json:"acquired_delegated_frozen_balance_for_energy,omitempty"`
		// Frozen balances provided to other accounts
		DelegatedFrozenBalanceForEnergy int64 `json:"delegated_frozen_balance_for_energy,omitempty"`

		// storage resource, get from market
		StorageLimit              int64 `json:"storage_limit,omitempty"`
		StorageUsage              int64 `json:"storage_usage,omitempty"`
		LatestExchangeStorageTime int64 `json:"latest_exchange_storage_time,omitempty"`
	} `json:"account_resource,omitempty"`
	CodeHash          string           `json:"codeHash,omitempty"`
	OwnerPermission   *HTTPPermision   `json:"owner_permission,omitempty"`
	WitnessPermission *HTTPPermision   `json:"witness_permission,omitempty"`
	ActivePermission  []*HTTPPermision `json:"active_permission,omitempty"`
}

// WARN: smart contract is not in core/tron proto, but
// https://github.com/tronprotocol/protocol/blob/master/core/contract/smart_contract.proto
type HTTPContract struct {
	OriginAddress   string `json:"origin_address,omitempty"`
	ContractAddress string `json:"contract_address,omitempty"`
	Abi             struct {
		Entrys []*HTTPABIEntry `json:"entrys,omitempty"`
	} `json:"abi,omitempty"`
	Bytecode                   string `json:"bytecode,omitempty"`
	CallValue                  int64  `json:"call_value,omitempty"`
	ConsumeUserResourcePercent int    `json:"consume_user_resource_percent,omitempty"`
	Name                       string `json:"name,omitempty"`
	OriginEnergyLimit          int64  `json:"origin_energy_limit,omitempty"`
	CodeHash                   string `json:"code_hash,omitempty"`
	TRXHash                    string `json:"trx_hash,omitempty"`
}

type HTTPABIEntry struct {
	Anonymous       bool             `json:"anonymous,omitempty"`
	Constant        bool             `json:"constant,omitempty"`
	Name            string           `json:"name,omitempty"`
	Inputs          []*HTTPABIParams `json:"inputs,omitempty"`
	Outputs         []*HTTPABIParams `json:"outputs,omitempty"`
	Type            string           `json:"type,omitempty"`
	Paybale         bool             `json:"payable,omitempty"`
	StateMutability string           `json:"stateMutability,omitempty"`
}

type HTTPABIParams struct {
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Indexed bool   `json:"indexed,omitempty"`
}
