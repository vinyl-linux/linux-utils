/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show the version of vinyl linux utils installed",
	Long:  "show the version of vinyl linux utils installed",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("linux-utils version\n---\nVersion: %s\nBuild User: %s\nBuilt On: %s\n",
			Ref, BuildUser, BuiltOn,
		)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
