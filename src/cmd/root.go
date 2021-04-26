package cmd

import (
	"SLC/src/test"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "slc",
	Short: "SLC(systemxrs linux container) is a experiment tool designed to make it easier to create, deploy, and runAction applications by using containers.",
	Long: `SLC(systemxrs linux container) is a experiment tool designed to make it easier to create, deploy, and runAction applications by using containers. Containers allow a developer to package up an application with all of the parts it needs, such as libraries and other dependencies, and deploy it as one package.`,
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("args: ", strings.Join(args, " "))
		isTest, _ := cmd.Flags().GetBool("test")
		if isTest == true { test.Run() }
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(InitConfig)

	rootCmd.Flags().BoolP("test", "t", false, "test")


	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("tty", "i", false, "runAction frontGround")

	rootCmd.AddCommand(daemonCmd)
	daemonCmd.Flags().BoolP("tty", "i", false, "runAction frontGround")
}



