# TRONETL

TRONETL is a ETL helper to TRON protocol

[中文版](./README.CHS.md)

## Prerequisites

1. RUN A FULLNODE WITH [TronDeploy](https://git.ngx.fi/c0mm4nd/trondeploy)!

2. Make sure the output-directory was downloaded and node started

3. install the tronetl
```bash
git clone https://git.ngx.fi/c0mm4nd/tronetl && cd tronetl

# if using docker (recommend)
docker build -t tronetl .
docker run -it tronetl -h

# else using latest golang
go install .
tronetl -h
```

## Usage

```bash
tronetl is a CLI tool for parsing blockchain data from tron network to CSV format files

Usage:
  tronetl [command]

Available Commands:
  completion                     Generate the autocompletion script for the specified shell
  export_blocks_and_transactions export blocks, with the blocks' trx and trc10 transactions
  export_token_transfers         export smart contract token's transfers
  help                           Help about any command
  server                         run a server for servings the export tasks

Flags:
  -h, --help   help for tronetl

Use "tronetl [command] --help" for more information about a command.
```

### export_blocks_and_transactions

```bash
export blocks, with the blocks' trx and trc10 transactions

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
```

### export_token_transfers

```bash
export smart contract token's transfers

Usage:
  tronetl export_token_transfers [flags]

Flags:
      --contracts stringArray   just output selected contracts' transfers
      --end-block uint          the ending block number
      --end-timestamp uint      the ending block's timestamp (in UTC)
  -h, --help                    help for export_token_transfers
      --output string           the CSV file for token transfer outputs, use - to omit (default "token_transfer.csv")
      --provider-uri string     the base uri of the tron fullnode (without port) (default "http://localhost")
      --start-block uint        the starting block number
      --start-timestamp uint    the starting block's timestamp (in UTC)
```
