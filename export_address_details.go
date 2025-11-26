package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"strings"
	"sync"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
	"github.com/jszwec/csvutil"
)

// ExportAddressDetailsOptions is the option for ExportAddressDetails func
type ExportAddressDetailsOptions struct {
	addrSource      io.Reader
	accountsOutput  io.Writer
	contractsOutput io.Writer
	tokensOutput    io.Writer

	Addresses []string

	ProviderURI string `json:"provider_uri,omitempty"`
}

func ExportAddressDetails(options *ExportAddressDetailsOptions) {
	allAddrs := collectAllAddrs(options)

	var accountsCsvEncoder, contractsEncoder, tokensEncoder *csvutil.Encoder
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

	if options.tokensOutput != nil {
		tokensCsvWriter := csv.NewWriter(options.tokensOutput)
		defer tokensCsvWriter.Flush()
		tokensEncoder = csvutil.NewEncoder(tokensCsvWriter)
	}

	cli := tron.NewTronClient(options.ProviderURI)
	for _, addr := range allAddrs {
		acc := cli.GetAccount(addr)
		if acc == nil {
			log.Printf("nil account for %s; skipping", addr)
			continue
		}
		if acc.Address == "" {
			acc.Address = addr // fallback to requested addr when API returns empty
		}

		if options.accountsOutput != nil {
			accountsCsvEncoder.Encode(NewCsvAccount(acc))
		}

		if options.contractsOutput != nil && strings.ToLower(acc.AccountType) == "contract" {
			contract := cli.GetContract(addr)
			csvContract := NewCsvContract(contract)
			contractsEncoder.Encode(csvContract)

			if options.tokensOutput != nil && (csvContract.IsErc20 || csvContract.IsErc721) {
				tokensEncoder.Encode(NewCsvTokens(cli, contract))
			}
		}

		// TODO: support type == AssetIssue
	}
}

// ExportAddressDetailsWithWorkers runs ExportAddressDetails logic concurrently when workers > 0.
func ExportAddressDetailsWithWorkers(options *ExportAddressDetailsOptions, workers uint) {
	allAddrs := collectAllAddrs(options)

	var receiverWG sync.WaitGroup

	var accountsEncCh, contractsEncCh, tokensEncCh chan any
	if options.accountsOutput != nil {
		accountsCsvWriter := csv.NewWriter(options.accountsOutput)
		defer accountsCsvWriter.Flush()
		accountsEncoder := csvutil.NewEncoder(accountsCsvWriter)
		accountsEncCh = createCSVEncodeCh(&receiverWG, accountsEncoder, workers)
	}

	if options.contractsOutput != nil {
		contractsCsvWriter := csv.NewWriter(options.contractsOutput)
		defer contractsCsvWriter.Flush()
		contractsEncoder := csvutil.NewEncoder(contractsCsvWriter)
		contractsEncCh = createCSVEncodeCh(&receiverWG, contractsEncoder, workers)
	}

	if options.tokensOutput != nil {
		tokensCsvWriter := csv.NewWriter(options.tokensOutput)
		defer tokensCsvWriter.Flush()
		tokensEncoder := csvutil.NewEncoder(tokensCsvWriter)
		tokensEncCh = createCSVEncodeCh(&receiverWG, tokensEncoder, workers)
	}

	cli := tron.NewTronClient(options.ProviderURI)

	exportWork := func(wg *sync.WaitGroup, workerID uint) {
		for idx := workerID; idx < uint(len(allAddrs)); idx += workers {
			addr := allAddrs[idx]
			acc := cli.GetAccount(addr)
			if acc == nil {
				log.Printf("nil account for %s; skipping", addr)
				continue
			}
			if acc.Address == "" {
				acc.Address = addr
			}

			if accountsEncCh != nil {
				accountsEncCh <- NewCsvAccount(acc)
			}

			if contractsEncCh != nil && strings.ToLower(acc.AccountType) == "contract" {
				contract := cli.GetContract(addr)
				csvContract := NewCsvContract(contract)
				contractsEncCh <- csvContract

				if tokensEncCh != nil && (csvContract.IsErc20 || csvContract.IsErc721) {
					tokensEncCh <- NewCsvTokens(cli, contract)
				}
			}
		}
		wg.Done()
	}

	var senderWG sync.WaitGroup
	for workerID := uint(0); workerID < workers; workerID++ {
		senderWG.Add(1)
		go exportWork(&senderWG, workerID)
	}

	senderWG.Wait()
	if accountsEncCh != nil {
		close(accountsEncCh)
	}
	if contractsEncCh != nil {
		close(contractsEncCh)
	}
	if tokensEncCh != nil {
		close(tokensEncCh)
	}
	receiverWG.Wait()

	log.Printf("exported %d addresses with %d workers", len(allAddrs), workers)
}

func collectAllAddrs(options *ExportAddressDetailsOptions) []string {
	// find all 34 length T-addr
	allAddrs := make([]string, 0, len(options.Addresses))
	if options.addrSource != nil {
		scanner := bufio.NewScanner(options.addrSource)
		for scanner.Scan() {
			line := scanner.Text()
			for _, sub := range strings.Split(line, ",") {
				if len(sub) > 0 && sub[0] == 'T' && len(sub) == 34 {
					// =Taddr
					allAddrs = append(allAddrs, tron.EnsureHexAddr(sub))
				}
			}
		}
	}

	for i := range options.Addresses {
		allAddrs = append(allAddrs, tron.EnsureHexAddr(options.Addresses[i]))
	}

	return allAddrs
}
