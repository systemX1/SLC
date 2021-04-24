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
		var cmds []string
		for _, v := range args {
			cmds = append(cmds, v)
		}

		tty, err := cmd.Flags().GetBool("tty")
		if err != nil {
			fmt.Println(err)
		}

		run(cmds, tty)
	},
}

func run(cmds []string, tty bool) {
	daemon.RunParentProcess(cmds, tty)

}




