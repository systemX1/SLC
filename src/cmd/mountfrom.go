package cmd

import (
	"SLC/src/reexec"
	"github.com/spf13/cobra"
)

var mountfromCmd = &cobra.Command{
	Use:   "mountfrom",
	Short: "mountfrom short",
	Long:  `mountfrom long`,
	Run: func(cmd *cobra.Command, args []string) {
		reexec.Init("slc-mountfrom")
	},
}
