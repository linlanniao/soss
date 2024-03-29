package cmd

import (
	"os"

	"github.com/linlanniao/soss/internal/controller"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all objects of s3Service",
	Run: func(cmd *cobra.Command, _ []string) {
		opts := controller.ListOptions{
			Bucket: bucket,
			Prefix: listPrefix,
		}

		if err := ctrl.List(opts); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	},
}

var (
	listPrefix string
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listPrefix, "prefix", "p", "", `object prefix to list (default "")`)
}
