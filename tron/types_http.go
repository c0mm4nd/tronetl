package tron

type TxInfo struct {
	Log []struct {
		Address string   `json:"address"`
		Data    string   `json:"data"`
		Topics  []string `json:"topics"`
	} `json:"log,omitempty"`
	Fee            int      `json:"fee,omitempty"`
	BlockNumber    int      `json:"blockNumber"`
	ContractResult []string `json:"contractResult"`
	BlockTimeStamp int64    `json:"blockTimeStamp"`
	Receipt        struct {
		Result            string `json:"result"`
		NetFee            int    `json:"net_fee"`
		EnergyUsageTotal  int    `json:"energy_usage_total"`
		OriginEnergyUsage int    `json:"origin_energy_usage"`
	} `json:"receipt"`
	ID                   string `json:"id"`
	ContractAddress      string `json:"contract_address,omitempty"`
	InternalTransactions []struct {
		CallerAddress     string `json:"caller_address"`
		Note              string `json:"note"`
		TransferToAddress string `json:"transferTo_address"`
		CallValueInfo     []struct {
		} `json:"callValueInfo"`
		Hash string `json:"hash"`
	} `json:"internal_transactions,omitempty"`
}

type HTTPBlock struct {
	BlockID     string `json:"blockID"`
	BlockHeader struct {
		RawData struct {
			Number         int    `json:"number"`
			TxTrieRoot     string `json:"txTrieRoot"`
			WitnessAddress string `json:"witness_address"`
			ParentHash     string `json:"parentHash"`
			Version        int    `json:"version"`
			Timestamp      int64  `json:"timestamp"`
		} `json:"raw_data"`
		WitnessSignature string `json:"witness_signature"`
	} `json:"block_header"`
	Transactions []HTTPTransaction `json:"transactions"`
}

// Values: https://tronprotocol.github.io/documentation-en/mechanism-algorithm/system-contracts/
//TransferAssetContract
//TriggerSmartContract
//TransferContract

type HTTPTransaction struct {
	Ret []struct {
		ContractRet string `json:"contractRet"`
	} `json:"ret"`
	Signature []string `json:"signature"`
	TxID      string   `json:"txID"`
	RawData   struct {
		Data     string `json:"data"`
		Contract []struct {
			Parameter struct {
				Value struct {
					AssetName    string `json:"asset_name"`
					Amount       uint64 `json:"amount"`
					OwnerAddress string `json:"owner_address"`
					ToAddress    string `json:"to_address"`
				} `json:"value"`
				TypeURL string `json:"type_url"`
			} `json:"parameter"`
			Type string `json:"type"`
		} `json:"contract"`
		RefBlockBytes string `json:"ref_block_bytes"`
		RefBlockHash  string `json:"ref_block_hash"`
		Expiration    uint64 `json:"expiration"`
		Timestamp     uint64 `json:"timestamp"`
		FeeLimit      uint64 `json:"fee_limit"`
	} `json:"raw_data"`
	RawDataHex string `json:"raw_data_hex"`
}
