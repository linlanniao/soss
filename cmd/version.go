/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/linlanniao/soss/pkg/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "show version information",
	Aliases: []string{"version", "v"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Print("soss"))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
