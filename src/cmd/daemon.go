package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"helloDB/src/daemon"
)



var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "daemon short",
	Long: `daemon long`,
	Run: func(cmd *cobra.Command, args []string) {
		tty, err := cmd.Flags().GetBool("tty")
		if err != nil {
			fmt.Println(err)
		}
		daemon.NewInitProcess(tty)
	},
}


