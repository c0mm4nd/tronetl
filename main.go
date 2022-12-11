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
	defaults := pflag.NewFlagSet("defaults for all commands", pflag.ExitOnError)
	providerURI := defaults.String("provider-uri", "http://localhost", "the base uri of the tron fullnode (without port)")
	startBlock := defaults.Uint64("start-block", 0, "the starting block number")
	endBlock := defaults.Uint64("end-block", 0, "the ending block number")
	startTimestamp := defaults.Uint64("start-timestamp", 0, "the starting block's timestamp (in UTC)")
	endTimestamp := defaults.Uint64("end-timestamp", 0, "the ending block's timestamp (in UTC)")

	cmdBlocksAndTxs := pflag.NewFlagSet("export_blocks_and_transactions", pflag.ExitOnError)
	blksOutput := cmdBlocksAndTxs.String("blocks-output", "blocks.csv", "the CSV file for block outputs, use - to omit")
	txsOutput := cmdBlocksAndTxs.String("transactions-output", "transactions.csv", "the CSV file for transaction outputs, use - to omit")
	trc10Output := cmdBlocksAndTxs.String("trc10-output", "trc10.csv", "the CSV file for trc10 outputs, use - to omit")
	cmdBlocksAndTxs.AddFlagSet(defaults)

	cmdTokenTf := pflag.NewFlagSet("export_token_transfers", pflag.ExitOnError)
	tfOutput := cmdTokenTf.String("transfers-output", "token_transfers.csv", "the CSV file for token transfer outputs, use - to omit")
	logOutput := cmdTokenTf.String("logs-output", "logs.csv", "the CSV file for transaction log outputs, use - to omit")
	internalTxOutput := cmdTokenTf.String("internal-tx-output", "internal_transactions.csv", "the CSV file for internal transaction outputs, use - to omit")
	filterContracts := cmdTokenTf.StringArray("contracts", []string{}, "just output selected contracts' transfers")
	cmdTokenTf.AddFlagSet(defaults)

	exportBlocksAndTransactionsCmd := &cobra.Command{
		Use:   "export_blocks_and_transactions",
		Short: "export blocks, with the blocks' trx and trc10 transactions",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var blksOut, txsOut, trc10Out *os.File
			if *blksOutput != "-" {
				blksOut, err = os.Create(*blksOutput)
				chk(err)
			}

			if *txsOutput != "-" {
				txsOut, err = os.Create(*txsOutput)
				chk(err)
			}

			if *trc10Output != "-" {
				trc10Out, err = os.Create(*trc10Output)
				chk(err)
			}

			ExportBlocksAndTransactions(&ExportBlocksAndTransactionsOptions{
				blksOutput:  blksOut,
				txsOutput:   txsOut,
				trc10Output: trc10Out,

				ProviderURI: *providerURI,

				StartBlock: *startBlock,
				EndBlock:   *endBlock,

				StartTimestamp: *startTimestamp,
				EndTimestamp:   *endTimestamp,

				WithTRXTransactions: txsOut != nil,
				WithTRC10Transfers:  trc10Out != nil,
			})
		},
	}
	exportBlocksAndTransactionsCmd.Flags().AddFlagSet(cmdBlocksAndTxs)

	exportTokenTransfersCmd := &cobra.Command{
		Use:   "export_token_transfers",
		Short: "export smart contract token's transfers",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var tfOut *os.File
			if *txsOutput != "-" {
				tfOut, err = os.Create(*tfOutput)
				chk(err)
			}

			var logOut *os.File
			if *txsOutput != "-" {
				logOut, err = os.Create(*logOutput)
				chk(err)
			}

			var internalTxOut *os.File
			if *txsOutput != "-" {
				internalTxOut, err = os.Create(*internalTxOutput)
				chk(err)
			}

			ExportTransfers(&ExportTransferOptions{
				tfOutput:         tfOut,
				logOutput:        logOut,
				internalTxOutput: internalTxOut,

				ProviderURI: *providerURI,
				StartBlock:  *startBlock,
				EndBlock:    *endBlock,

				StartTimestamp: *startTimestamp,
				EndTimestamp:   *endTimestamp,
				Contracts:      *filterContracts,
			})
		},
	}
	exportTokenTransfersCmd.Flags().AddFlagSet(cmdTokenTf)

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

					WithTRXTransactions: txsOut != nil,   // TODO
					WithTRC10Transfers:  trc10Out != nil, // TODO
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

	rootCmd.AddCommand(exportBlocksAndTransactionsCmd)
	rootCmd.AddCommand(exportTokenTransfersCmd)
	rootCmd.AddCommand(serverCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
