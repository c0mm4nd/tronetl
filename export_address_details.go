package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"strings"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
	"github.com/jszwec/csvutil"
)

// ExportAddressDetailsOptions is the option for ExportAddressDetails func
type ExportAddressDetailsOptions struct {
	addrSource      io.Reader
	accountsOutput  io.Writer
	contractsOutput io.Writer

	Addresses []string

	ProviderURI string `json:"provider_uri,omitempty"`
}

func ExportAddressDetails(options *ExportAddressDetailsOptions) {
	// find all 34 length T-addr
	allAddrs := make([]string, 0, len(options.Addresses))
	if options.addrSource != nil {
		scanner := bufio.NewScanner(options.addrSource)
		for scanner.Scan() {
			line := scanner.Text()
			for _, sub := range strings.Split(line, ",") {
				if sub[0] == 'T' && len(sub) == 34 {
					// =Taddr
					allAddrs = append(allAddrs, Tstring2HexAddr(sub))
				}
			}
		}
	}

	for i := range options.Addresses {
		allAddrs = append(allAddrs, Tstring2HexAddr(options.Addresses[i]))
	}

	var accountsCsvEncoder, contractsEncoder *csvutil.Encoder
	if options.accountsOutput != nil {
		accountsCsvWriter := csv.NewWriter(options.accountsOutput)
		defer accountsCsvWriter.Flush()
		accountsCsvEncoder = csvutil.NewEncoder(accountsCsvWriter)
	}

	if options.contractsOutput != nil {
		contractsCsvWriter := csv.NewWriter(options.contractsOutput)
		defer contractsCsvWriter.Flush()
		contractsEncoder = csvutil.NewEncoder(contractsCsvWriter)
	}

	cli := tron.NewTronClient(options.ProviderURI)
	for _, addr := range allAddrs {
		acc := cli.GetAccount(addr)

		if options.accountsOutput != nil {
			accountsCsvEncoder.Encode(NewCsvAccount(acc))
		}

		if options.contractsOutput != nil && strings.ToLower(acc.AccountType) == "contract" {
			contract := cli.GetContract(addr)
			contractsEncoder.Encode(NewCsvContract(contract))
			// TODO: add token output
		}

		// TODO: support type == AssetIssue
	}
}
