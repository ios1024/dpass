package main

import (
	"fmt"
	"os"

	"github.com/rbee3u/ios1024/internal/ios1024/aes256"
	"github.com/rbee3u/ios1024/internal/ios1024/qrcode"
	"github.com/rbee3u/ios1024/internal/ios1024/shamir"
	"github.com/spf13/cobra"
)

func main() {
	if err := newCmd().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "ios1024", Args: cobra.NoArgs}
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.AddCommand(
		aes256.NewCmdEncrypt(),
		aes256.NewCmdDecrypt(),
		shamir.NewCmdSplit(),
		shamir.NewCmdCombine(),
		qrcode.NewCmd(),
	)
	return cmd
}
