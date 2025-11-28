package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"strings"
	"sync"
	"sync/atomic"

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
			acc.Address = tron.EnsureTAddr(addr) // ensure downstream uses T-addr
		}

		if options.accountsOutput != nil {
			accountsCsvEncoder.Encode(NewCsvAccount(acc))
		}

		if options.contractsOutput != nil {
			contract := cli.GetContract(addr)
			if contract != nil && contract.ContractAddress != "" {
				csvContract := NewCsvContract(contract)
				contractsEncoder.Encode(csvContract)

				if options.tokensOutput != nil && (csvContract.IsErc20 || csvContract.IsErc721) {
					if tokens := NewCsvTokens(cli, contract); tokens != nil {
						tokensEncoder.Encode(tokens)
					} else {
						log.Printf("skip writing tokens for %s: unable to parse token metadata", addr)
					}
				}
			}
		}

		// TODO: support type == AssetIssue
	}
}

// ExportAddressDetailsWithWorkers runs ExportAddressDetails logic concurrently when workers > 0.
func ExportAddressDetailsWithWorkers(options *ExportAddressDetailsOptions, workers uint) {
	allAddrs := collectAllAddrs(options)
	var processed atomic.Uint64
	log.Printf("start exporting %d addresses with %d workers", len(allAddrs), workers)

	var receiverWG sync.WaitGroup

	var accountsEncCh, contractsEncCh, tokensEncCh chan any
	if options.accountsOutput != nil {
		accountsCsvWriter := csv.NewWriter(options.accountsOutput)
		defer accountsCsvWriter.Flush()
		accountsEncoder := csvutil.NewEncoder(accountsCsvWriter)
		accountsEncCh = createCSVEncodeCh(&receiverWG, accountsEncoder, accountsCsvWriter, workers)
	}

	if options.contractsOutput != nil {
		contractsCsvWriter := csv.NewWriter(options.contractsOutput)
		defer contractsCsvWriter.Flush()
		contractsEncoder := csvutil.NewEncoder(contractsCsvWriter)
		contractsEncCh = createCSVEncodeCh(&receiverWG, contractsEncoder, contractsCsvWriter, workers)
	}

	if options.tokensOutput != nil {
		tokensCsvWriter := csv.NewWriter(options.tokensOutput)
		defer tokensCsvWriter.Flush()
		tokensEncoder := csvutil.NewEncoder(tokensCsvWriter)
		tokensEncCh = createCSVEncodeCh(&receiverWG, tokensEncoder, tokensCsvWriter, workers)
	}

	cli := tron.NewTronClient(options.ProviderURI)

	exportWork := func(wg *sync.WaitGroup, workerID uint) {
		for idx := workerID; idx < uint(len(allAddrs)); idx += workers {
			addr := allAddrs[idx]
			acc := cli.GetAccount(addr)
			if acc == nil {
				log.Printf("nil account for %s; skipping", addr)
				if n := processed.Add(1); n <= 1000 && n%100 == 0 || n%10000 == 0 {
					log.Printf("progress: %d / %d addresses processed", n, len(allAddrs))
				}
				continue
			}
			if acc.Address == "" {
				acc.Address = tron.EnsureTAddr(addr)
			}

			if accountsEncCh != nil {
				accountsEncCh <- NewCsvAccount(acc)
			}

			if contractsEncCh != nil {
				contract := cli.GetContract(addr)
				if contract != nil && contract.ContractAddress != "" {
					csvContract := NewCsvContract(contract)
					contractsEncCh <- csvContract

					if tokensEncCh != nil && (csvContract.IsErc20 || csvContract.IsErc721) {
						if tokens := NewCsvTokens(cli, contract); tokens != nil {
							tokensEncCh <- tokens
						} else {
							log.Printf("skip writing tokens for %s: unable to parse token metadata", addr)
						}
					}
				}
			}

			if n := processed.Add(1); n%10000 == 0 {
				log.Printf("progress: %d / %d addresses processed", n, len(allAddrs))
			} else if n <= 1000 && n%100 == 0 { // finer-grain progress for the first batch
				log.Printf("progress: %d / %d addresses processed", n, len(allAddrs))
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
	// find all 34 length T-addr (or convert hex to T-addr)
	allAddrs := make([]string, 0, len(options.Addresses))
	if options.addrSource != nil {
		scanner := bufio.NewScanner(options.addrSource)
		for scanner.Scan() {
			line := scanner.Text()
			for _, sub := range strings.Split(line, ",") {
				if len(sub) == 0 {
					continue
				}
				switch {
				case sub[0] == 'T' && len(sub) == 34:
					allAddrs = append(allAddrs, sub) // already T-addr
				case len(sub) >= 40 && len(sub) <= 64:
					allAddrs = append(allAddrs, tron.EnsureTAddr(sub)) // likely hex, convert
				default:
					// skip header or malformed rows
					continue
				}
			}
		}
	}

	for i := range options.Addresses {
		addr := options.Addresses[i]
		switch {
		case len(addr) > 0 && addr[0] == 'T' && len(addr) == 34:
			allAddrs = append(allAddrs, addr)
		case len(addr) >= 40 && len(addr) <= 64:
			allAddrs = append(allAddrs, tron.EnsureTAddr(addr))
		default:
			// skip invalid address
		}
	}

	return allAddrs
}
