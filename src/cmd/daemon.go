package cmd

import (
	"SLC/src/daemon"
	"fmt"
	"github.com/spf13/cobra"
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

		var cmds []string
		for _, v := range args {
			cmds = append(cmds, v)
		}

		//daemon.NewInitProcess(tty)
		daemon.Init(cmds, tty)
	},
}


