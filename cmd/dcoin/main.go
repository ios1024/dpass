package main

import (
	"fmt"
	"os"

	"github.com/ios1024/dpass/internal/dcoin/bitcoin"
	"github.com/ios1024/dpass/internal/dcoin/dogecoin"
	"github.com/ios1024/dpass/internal/dcoin/ethereum"
	"github.com/ios1024/dpass/internal/dcoin/mnemonic"
	"github.com/ios1024/dpass/internal/dcoin/solana"
	"github.com/ios1024/dpass/internal/dcoin/sui"
	"github.com/ios1024/dpass/internal/dcoin/tron"
	"github.com/spf13/cobra"
)

func main() {
	if err := newCmd().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "dcoin", Args: cobra.NoArgs}
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.AddCommand(
		mnemonic.NewCmd(),
		bitcoin.NewCmd(),
		ethereum.NewCmd(),
		tron.NewCmd(),
		solana.NewCmd(),
		dogecoin.NewCmd(),
		sui.NewCmd(),
	)
	return cmd
}
