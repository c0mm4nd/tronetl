# TRONETL

[![release](https://github.com/c0mm4nd/tronetl/actions/workflows/release.yml/badge.svg?branch=dev)](https://github.com/c0mm4nd/tronetl/actions/workflows/release.yml)

TRONETL is a CLI tool for parsing blockchain data from the TRON network to CSV format files.

[中文版](./README.CHS.md)

## Prerequisites

1.  Run a TRON full node. You can use [TronDeploy](https://github.com/c0mm4nd/trondeploy) to set one up.
2.  Make sure the full node is running and accessible.
3.  Install `tronetl`.

## Installation

git clone https://github.com/c0mm4nd/tronetl && cd tronetl
docker build -t tronetl .
docker run -it tronetl -h
```

### Using Go

```bash
git clone https://github.com/c0mm4nd/tronetl && cd tronetl
go install .
tronetl -h
```

## Usage

```
tronetl is a CLI tool for parsing blockchain data from tron network to CSV format files

Usage:
  tronetl [command]

Available Commands:
  completion                     Generate the autocompletion script for the specified shell
  export_blocks_and_transactions export blocks, with the blocks' trx and trc10 transactions
  export_token_transfers         export smart contract token's transfers
  export_address_details         export the addresses' type and smart contract related details
  help                           Help about any command
  server                         run a server for servings the export tasks

Flags:
  -h, --help   help for tronetl

Use "tronetl [command] --help" for more information about a command.
```

ETL results' schema are written [in this doc](./SCHEMA.md)

### export_blocks_and_transactions

Exports blocks, transactions, and TRC10 transfers to CSV files.

```
Usage:
  tronetl export_blocks_and_transactions [flags]

Flags:
      --blocks-output string         the CSV file for block outputs, use - to omit (default "blocks.csv")
      --end-block uint               the ending block number
      --end-timestamp uint           the ending block's timestamp (in UTC)
  -h, --help                         help for export_blocks_and_transactions
      --provider-uri string          the base uri of the tron fullnode (without port) (default "http://localhost")
      --start-block uint             the starting block number
      --start-timestamp uint         the starting block's timestamp (in UTC)
      --transactions-output string   the CSV file for transaction outputs, use - to omit (default "transactions.csv")
      --trc10-output string          the CSV file for trc10 outputs, use - to omit (default "trc10.csv")
      --workers uint                 the count of the workers in parallel
```

### export_token_transfers

Exports TRC20 token transfers, logs, internal transactions, receipts, and newly created contracts to CSV files.

```
Usage:
  tronetl export_token_transfers [flags]

Flags:
      --contracts stringArray         just output selected contracts' transfers
      --end-block uint                the ending block number
      --end-timestamp uint            the ending block's timestamp (in UTC)
  -h, --help                          help for export_token_transfers
      --internal-tx-output string     the CSV file for internal transaction outputs, use - to omit (default "internal_transactions.csv")
      --logs-output string            the CSV file for transaction log outputs, use - to omit (default "logs.csv")
      --new-contracts-output string   the CSV file for newly created contract outputs, use - to omit (default "new_contracts.csv")
      --provider-uri string           the base uri of the tron fullnode (without port) (default "http://localhost")
      --receipts-output string        the CSV file for transaction receipt outputs, use - to omit (default "receipts.csv")
      --start-block uint              the starting block number
      --start-timestamp uint          the starting block's timestamp (in UTC)
      --transfers-output string       the CSV file for token transfer outputs, use - to omit (default "token_transfers.csv")
      --workers uint                  the count of the workers in parallel
```

### export_address_details

Exports address details, including account info, contract details, and token information.

```
Usage:
  tronetl export_address_details [flags]

Flags:
      --accounts-output string    the CSV file for all account info outputs, use - to omit (default "accounts.csv")
      --addrs stringArray         a supplementary address list to load details
      --addrs-source string       the CSV file for block outputs, or use address (default "-")
      --contracts-output string   the CSV file for contract account detail outputs, use - to omit (default "contract.csv")
  -h, --help                      help for export_address_details
      --provider-uri string       the base uri of the tron fullnode (without port) (default "http://localhost")
      --tokens-output string      the CSV file for token contract detail outputs, use - to omit (default "tokens.csv")
```

### server

Runs a web server that exposes the export functionalities through a REST API.

The server listens on port `54173`.

**Endpoints:**

*   `/export_blocks_and_transactions`:  Exports blocks and transactions.
    *   Query Parameters: `start-block`, `end-block`, `start-timestamp`, `end-timestamp`
*   `/export_token_transfers`: Exports token transfers.
    *   Query Parameters: `start-block`, `end-block`, `start-timestamp`, `end-timestamp`, `contracts` (can be repeated)

The output is a zip file containing the corresponding CSV files.