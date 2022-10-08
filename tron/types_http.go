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
