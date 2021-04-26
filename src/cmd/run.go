package cmd

import (
	"SLC/src/daemon"
	"fmt"
	"github.com/spf13/cobra"
	"os"
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
		for _, v := range os.Args[2:] {
			cmds = append(cmds, v)
		}

		runAction(cmds, tty)
	},
}

// cmds, whether run front
func runAction(cmds []string, tty bool) {
	//daemon.RunParentProcess(cmds, tty)
	daemon.Init(cmds, tty)
}




