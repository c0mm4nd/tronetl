package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"

	"git.ngx.fi/c0mm4nd/tronetl/tron"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

func main() {
	defaults := pflag.NewFlagSet("defaults for all commands", pflag.ExitOnError)
	providerURI := defaults.String("provider-uri", "http://localhost", "visible for all commands")
	startBlock := defaults.Uint64("start-block", 0, "only visible for cmd A")
	endBlock := defaults.Uint64("end-block", 0, "only visible for cmd A")
	startTimestamp := defaults.Uint64("start-timestamp", 0, "only visible for cmd A")
	endTimestamp := defaults.Uint64("end-timestamp", 0, "only visible for cmd A")

	cmdBlocksAndTxs := pflag.NewFlagSet("export_blocks_and_transactions", pflag.ExitOnError)
	blksOutput := cmdBlocksAndTxs.String("blocks-output", "blocks.csv", "blocks output")
	txsOutput := cmdBlocksAndTxs.String("transactions-output", "transactions.csv", "transactions output")
	cmdBlocksAndTxs.AddFlagSet(defaults)

	cmdTokenTf := pflag.NewFlagSet("export_token_transfers", pflag.ExitOnError)
	tfOutput := cmdTokenTf.String("output", "token_transfer.csv", "transfer output")
	filterContracts := cmdTokenTf.StringArray("contracts", []string{}, "limit contracts")
	// withBlockOutput := cmdTokenTf.String("output", "blocks.csv", "with blocks output")
	// withTransactionOutput := cmdTokenTf.String("output", "txs.csv", "with transfer output")
	cmdTokenTf.AddFlagSet(defaults)
	// defaults.Parse(os.Args)

	if len(os.Args) == 1 {
		log.Fatal("no subcommand given")
	}

	switch os.Args[1] {
	case "export_blocks_and_transactions":
		cmdBlocksAndTxs.Parse(os.Args[2:])
		exportBlocksAndTransactions(*providerURI, *startBlock, *endBlock, *blksOutput, *txsOutput)
	case "export_token_transfers":
		cmdTokenTf.Parse(os.Args[2:])

		outFile, err := os.Create(*tfOutput)
		chk(err)

		// var blockOutput io.Writer
		// if *withBlockOutput != "" {
		// 	blockOutput, err = os.Create(*withBlockOutput)
		// 	chk(err)
		// }

		// var transactionOutput io.Writer
		// if *withTransactionOutput != "" {
		// 	transactionOutput, err = os.Create(*transactionOutput)
		// 	chk(err)
		// }

		// writer, err := rollingwriter.NewWriterFromConfig(&rollingwriter.Config{
		// 	LogPath:  "/mnt/hdd14t/tron_out/",
		// 	FileName: *tfOutput,
		// 	RollingVolumeSize:      "2G",
		// })

		exportTransfers(&ExportTransferOptions{
			output: outFile,
			// withBlockOutput:       blockOutput,
			// withTransactionOutput: transactionOutput,

			ProviderURI: *providerURI,
			StartBlock:  *startBlock,
			EndBlock:    *endBlock,

			StartTimestamp: *startTimestamp,
			EndTimestamp:   *endTimestamp,
			Contracts:      *filterContracts,
		})
	case "server":
		// run server by default
		cli := tron.NewTronClient("http://localhost")

		latestBlock := cli.GetJSONBlockByNumber(nil)
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
	}
}
