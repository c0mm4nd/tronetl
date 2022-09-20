package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	defaults := pflag.NewFlagSet("defaults for all commands", pflag.ExitOnError)
	providerURL := defaults.String("provider-uri", "http://localhost:8090", "visible for all commands")
	startBlock := defaults.Uint64("start-block", 0, "only visible for cmd A")
	endBlock := defaults.Uint64("end-block", 0, "only visible for cmd A")

	// cmdBlocksAndTxs := pflag.NewFlagSet("export_blocks_and_transactions", pflag.ExitOnError)
	// cmdBlocksAndTxs.AddFlagSet(defaults)

	cmdTokenTf := pflag.NewFlagSet("export_token_transfers", pflag.ExitOnError)
	tfOutput := cmdTokenTf.String("output", "token_transfer.csv", "transfer output")
	cmdTokenTf.AddFlagSet(defaults)
	// defaults.Parse(os.Args)

	if len(os.Args) == 1 {
		log.Fatal("no subcommand given")
	}

	switch os.Args[1] {
	// case "export_blocks_and_transactions":
	// 	cmdBlocksAndTxs.Parse(os.Args[2:])
	// 	fmt.Println(*cmdAspecific)
	// 	fmt.Println(*forAll)
	case "export_token_transfers":
		cmdTokenTf.Parse(os.Args[2:])
		exportTransfers(*providerURL, *startBlock, *endBlock, *tfOutput)
	default:
		fmt.Printf("%q is no valid subcommand.\n", os.Args[1])
		os.Exit(2)
	}
}
