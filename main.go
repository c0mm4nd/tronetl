package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"archive/zip"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "tronetl",
	Short: "tronetl",
	Long:  `tronetl is a CLI tool for parsing blockchain data from tron network to CSV format files`,
}

func main() {
	nodeConfigs := pflag.NewFlagSet("node config", pflag.ExitOnError)
	providerURI := nodeConfigs.String("provider-uri", "http://localhost", "the base uri of the tron fullnode (without port)")

	defaults := pflag.NewFlagSet("defaults for all commands", pflag.ExitOnError)
	startBlock := defaults.Uint64("start-block", 0, "the starting block number")
	endBlock := defaults.Uint64("end-block", 0, "the ending block number")
	startTimestamp := defaults.Uint64("start-timestamp", 0, "the starting block's timestamp (in UTC)")
	endTimestamp := defaults.Uint64("end-timestamp", 0, "the ending block's timestamp (in UTC)")
	defaults.AddFlagSet(nodeConfigs)

	cmdBlocksAndTxs := pflag.NewFlagSet("export_blocks_and_transactions", pflag.ExitOnError)
	blksOutput := cmdBlocksAndTxs.String("blocks-output", "blocks.csv", "the CSV file for block outputs, use - to omit")
	txsOutput := cmdBlocksAndTxs.String("transactions-output", "transactions.csv", "the CSV file for transaction outputs, use - to omit")
	trc10Output := cmdBlocksAndTxs.String("trc10-output", "trc10.csv", "the CSV file for trc10 outputs, use - to omit")
	cmdBlocksAndTxs.AddFlagSet(defaults)

	cmdTokenTf := pflag.NewFlagSet("export_token_transfers", pflag.ExitOnError)
	tfOutput := cmdTokenTf.String("transfers-output", "token_transfers.csv", "the CSV file for token transfer outputs, use - to omit")
	logOutput := cmdTokenTf.String("logs-output", "logs.csv", "the CSV file for transaction log outputs, use - to omit")
	internalTxOutput := cmdTokenTf.String("internal-tx-output", "internal_transactions.csv", "the CSV file for internal transaction outputs, use - to omit")
	receiptOutput := cmdTokenTf.String("receipts-output", "receipts.csv", "the CSV file for transaction receipt outputs, use - to omit")
	filterContracts := cmdTokenTf.StringArray("contracts", []string{}, "just output selected contracts' transfers")
	cmdTokenTf.AddFlagSet(defaults)

	cmdAddrDetails := pflag.NewFlagSet("export_address_details", pflag.ExitOnError)
	addrs := cmdAddrDetails.StringArray("addrs", []string{}, "a supplementary address list to load details")
	addrsSource := cmdAddrDetails.String("addrs-source", "-", "the CSV file for block outputs, or use address")
	accountsOutput := cmdAddrDetails.String("accounts-output", "accounts.csv", "the CSV file for all account info outputs, use - to omit")
	contractsOutput := cmdAddrDetails.String("contracts-output", "contract.csv", "the CSV file for contract account detail outputs, use - to omit")
	tokensOutput := cmdAddrDetails.String("tokens-output", "tokens.csv", "the CSV file for token contract detail outputs, use - to omit")
	cmdAddrDetails.AddFlagSet(nodeConfigs)

	exportBlocksAndTransactionsCmd := &cobra.Command{
		Use:   "export_blocks_and_transactions",
		Short: "export blocks, with the blocks' trx and trc10 transactions",
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			options := &ExportBlocksAndTransactionsOptions{
				ProviderURI: *providerURI,

				StartBlock: *startBlock,
				EndBlock:   *endBlock,

				StartTimestamp: *startTimestamp,
				EndTimestamp:   *endTimestamp,
			}

			if *blksOutput != "-" {
				options.blksOutput, err = os.Create(*blksOutput)
				chk(err)
			}

			if *txsOutput != "-" {
				options.txsOutput, err = os.Create(*txsOutput)
				chk(err)
			}

			if *trc10Output != "-" {
				options.trc10Output, err = os.Create(*trc10Output)
				chk(err)
			}

			ExportBlocksAndTransactions(options)
		},
	}
	exportBlocksAndTransactionsCmd.Flags().AddFlagSet(cmdBlocksAndTxs)

	exportTokenTransfersCmd := &cobra.Command{
		Use:   "export_token_transfers",
		Short: "export smart contract token's transfers",
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			options := &ExportTransferOptions{
				ProviderURI: *providerURI,
				StartBlock:  *startBlock,
				EndBlock:    *endBlock,

				StartTimestamp: *startTimestamp,
				EndTimestamp:   *endTimestamp,
				Contracts:      *filterContracts,
			}

			if *tfOutput != "-" {
				options.tfOutput, err = os.Create(*tfOutput)
				chk(err)
			}

			if *logOutput != "-" {
				options.logOutput, err = os.Create(*logOutput)
				chk(err)
			}

			if *internalTxOutput != "-" {
				options.internalTxOutput, err = os.Create(*internalTxOutput)
				chk(err)
			}

			if *receiptOutput != "-" {
				options.receiptOutput, err = os.Create(*receiptOutput)
				chk(err)
			}

			ExportTransfers(options)
		},
	}
	exportTokenTransfersCmd.Flags().AddFlagSet(cmdTokenTf)

	exportAddressDetailsCmd := &cobra.Command{
		Use:   "export_address_details",
		Short: "export the addresses' type and smart contract related details (require T-address format)",
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			options := &ExportAddressDetailsOptions{
				ProviderURI: *providerURI,

				Addresses: *addrs,
			}

			if *addrsSource != "" && *addrsSource != "-" {
				options.addrSource, err = os.Open(*addrsSource)
				chk(err)
			}

			if *accountsOutput != "-" {
				options.accountsOutput, err = os.Create(*accountsOutput)
				chk(err)
			}

			if *contractsOutput != "-" {
				options.contractsOutput, err = os.Create(*contractsOutput)
				chk(err)
			}

			if *tokensOutput != "-" {
				options.tokensOutput, err = os.Create(*tokensOutput)
				chk(err)
			}

			ExportAddressDetails(options)
		},
	}
	exportAddressDetailsCmd.Flags().AddFlagSet(cmdAddrDetails)

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "run a server for servings the export tasks",
		Run: func(cmd *cobra.Command, args []string) {
			cli := tron.NewTronClient("http://localhost")

			latestBlock := cli.GetJSONBlockByNumberWithTxIDs(nil)
			log.Printf("latest block: %d", *latestBlock.Number)

			tryStr2Uint := func(str string) uint64 {
				u, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					return 0
				}
				return u
			}

			r := gin.Default()
			r.GET("/export_blocks_and_transactions", func(ctx *gin.Context) {
				var zipBuffer *bytes.Buffer = new(bytes.Buffer)
				var zipWriter *zip.Writer = zip.NewWriter(zipBuffer)
				blksOut, err := zipWriter.Create("blocks.csv")
				chk(err)
				txsOut, err := zipWriter.Create("transactions.csv")
				chk(err)
				trc10Out, err := zipWriter.Create("trc10.csv")
				chk(err)

				options := &ExportBlocksAndTransactionsOptions{
					blksOutput:  blksOut,
					txsOutput:   txsOut,
					trc10Output: trc10Out,

					ProviderURI:    *providerURI,
					StartBlock:     tryStr2Uint(ctx.Query("start-block")),
					EndBlock:       tryStr2Uint(ctx.Query("end-block")),
					StartTimestamp: tryStr2Uint(ctx.Query("start-timestamp")),
					EndTimestamp:   tryStr2Uint(ctx.Query("end-timestamp")),
				}
				ExportBlocksAndTransactions(options)

				ctx.Header("Content-Disposition", "attachment;filename=export.zip")
				ctx.Data(http.StatusOK, "application/zip", zipBuffer.Bytes())
			}).GET("/export_token_transfers", func(ctx *gin.Context) {
				var zipBuffer *bytes.Buffer = new(bytes.Buffer)
				var zipWriter *zip.Writer = zip.NewWriter(zipBuffer)
				tfOut, err := zipWriter.Create("token_transfers.csv")
				chk(err)
				logOut, err := zipWriter.Create("logs.csv")
				chk(err)
				internalTxOut, err := zipWriter.Create("internal_transactions.csv")
				chk(err)

				options := &ExportTransferOptions{
					tfOutput:         tfOut,
					logOutput:        logOut,
					internalTxOutput: internalTxOut,
					ProviderURI:      *providerURI,
					StartBlock:       tryStr2Uint(ctx.Query("start-block")),
					EndBlock:         tryStr2Uint(ctx.Query("end-block")),
					StartTimestamp:   tryStr2Uint(ctx.Query("start-timestamp")),
					EndTimestamp:     tryStr2Uint(ctx.Query("end-timestamp")),
					Contracts:        ctx.QueryArray("contracts"),
				}
				ExportTransfers(options)

				ctx.Header("Content-Disposition", "attachment;filename=export.zip")
				ctx.Data(http.StatusOK, "application/zip", zipBuffer.Bytes())
			})
			r.Run(":54173")

		},
	}

	rootCmd.AddCommand(
		exportBlocksAndTransactionsCmd,
		exportTokenTransfersCmd,
		exportAddressDetailsCmd,
		serverCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
