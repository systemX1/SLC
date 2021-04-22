package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"helloDB/src/daemon"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run short",
	Long: `run long`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			fmt.Println(v)
		}

		tty, err := cmd.Flags().GetBool("tty")
		if err != nil {
			fmt.Println(err)
		}

		run(tty)
	},
}

func run(tty bool) {
	daemon.RunParentProcess(tty)

}




