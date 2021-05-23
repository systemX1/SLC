package cmd

import (
	"SLC/src/container/namespace"
	"fmt"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "container",
	Short: "container short",
	Long:  `container long`,
	Run: func(cmd *cobra.Command, args []string) {
		tty, err := cmd.Flags().GetBool("tty")
		if err != nil {
			fmt.Println(err)
		}

		var cmds []string
		for _, v := range args {
			cmds = append(cmds, v)
		}

		//container.NewInitProcess(tty)
		namespace.Init(cmds, tty)
	},
}
