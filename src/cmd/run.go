package cmd

import (
	"SLC/src/daemon"
	"fmt"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run short",
	Long: `run long`,
	Run: func(cmd *cobra.Command, args []string) {
		tty, err := cmd.Flags().GetBool("tty")
		if err != nil {
			fmt.Println(err)
		}

		var cmds []string
		cmds = append(cmds, "daemon")
		for _, v := range args {
			cmds = append(cmds, v)
		}

		run(cmds, tty)
	},
}

func run(cmds []string, tty bool) {
	//daemon.RunParentProcess(cmds, tty)
	daemon.Run()
}




