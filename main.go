package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "tronetl",
	Short: "tronetl",
	Long:  `tronetl is a CLI tool for parsing blockchain data from tron network`,
}

func main() {
	defaults := pflag.NewFlagSet("defaults for all commands", pflag.ExitOnError)
	providerURI := defaults.String("provider-uri", "http://localhost", "visible for all commands")
	startBlock := defaults.Uint64("start-block", 0, "only visible for cmd A")
	endBlock := defaults.Uint64("end-block", 0, "only visible for cmd A")
	startTimestamp := defaults.Uint64("start-timestamp", 0, "only visible for cmd A")
	endTimestamp := defaults.Uint64("end-timestamp", 0, "only visible for cmd A")

	cmdBlocksAndTxs := pflag.NewFlagSet("export_blocks", pflag.ExitOnError)
	blksOutput := cmdBlocksAndTxs.String("blocks-output", "blocks.csv", "blocks output")
	txsOutput := cmdBlocksAndTxs.String("transactions-output", "transactions.csv", "transactions output")
	trc10Output := cmdBlocksAndTxs.String("trc10-output", "trc10.csv", "trc10 output")
	cmdBlocksAndTxs.AddFlagSet(defaults)

	cmdTokenTf := pflag.NewFlagSet("export_token_transfers", pflag.ExitOnError)
	tfOutput := cmdTokenTf.String("output", "token_transfer.csv", "transfer output")
	filterContracts := cmdTokenTf.StringArray("contracts", []string{}, "limit contracts")
	cmdTokenTf.AddFlagSet(defaults)

	exportBlocksCmd := &cobra.Command{
		Use: "export_blocks",
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

			exportBlocksAndTransactions(&ExportBlocksOptions{
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
	exportBlocksCmd.Flags().AddFlagSet(cmdBlocksAndTxs)

	exportTokenTransfersCmd := &cobra.Command{
		Use: "export_token_transfers",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var tfOut *os.File
			if *txsOutput != "-" {
				tfOut, err = os.Create(*tfOutput)
				chk(err)
			}

			exportTransfers(&ExportTransferOptions{
				output: tfOut,
				// withBlockOutput:       blockOutput,
				// withTransactionOutput: transactionOutput,

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
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			cli := tron.NewTronClient("http://localhost")

			latestBlock := cli.GetJSONBlockByNumber(nil, false)
			log.Printf("latest block: %d", *latestBlock.Number)

			tryStr2Uint := func(str string) uint64 {
				u, err := strconv.ParseUint(str, 10, 64)
				if err != nil {
					return 0
				}
				return u
			}

			r := gin.Default()
			r.GET("/export_token_transfers", func(ctx *gin.Context) {
				// c, err := websocket.Accept(ctx.Writer, ctx.Request, &websocket.AcceptOptions{
				// 	InsecureSkipVerify: true,
				// 	OriginPatterns:     []string{"172.24.1.1:54173"},
				// })
				// chk(err)
				// defer c.Close(websocket.StatusInternalError, "the sky is falling")

				// writer, err := c.Writer(ctx, websocket.MessageText)
				// chk(err)
				writer := new(bytes.Buffer)
				options := &ExportTransferOptions{
					output:         writer,
					ProviderURI:    "http://localhost",
					StartBlock:     tryStr2Uint(ctx.Query("start-block")),
					EndBlock:       tryStr2Uint(ctx.Query("end-block")),
					StartTimestamp: tryStr2Uint(ctx.Query("start-timestamp")),
					EndTimestamp:   tryStr2Uint(ctx.Query("end-timestamp")),
					Contracts:      ctx.QueryArray("contracts"),
				}
				exportTransfers(options)

				// writer.Close()
				// c.Close(websocket.StatusNormalClosure, "")
				ctx.Header("Content-Disposition", "attachment;filename=token_transfer.csv")
				ctx.Data(http.StatusOK, "text/csv", writer.Bytes())
			})
			r.Run(":54173")

		},
	}

	rootCmd.AddCommand(exportBlocksCmd)
	rootCmd.AddCommand(exportTokenTransfersCmd)
	rootCmd.AddCommand(serverCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
