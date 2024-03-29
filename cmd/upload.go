package cmd

import (
	"os"

	"github.com/linlanniao/soss/internal/controller"
	"github.com/linlanniao/soss/pkg/utils"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var (
	uploadEncryptKey string
	uploadPrefix     string
	uploadCmd        = &cobra.Command{
		Use:   "upload files [files ...]",
		Short: "Encrypt and upload files to s3service",
		Run: func(cmd *cobra.Command, paths []string) {
			opts := controller.UploadOptions{
				Endpoint:   endpoint,
				Bucket:     bucket,
				Prefix:     uploadPrefix,
				EncryptKey: uploadEncryptKey,
				Paths:      utils.RemoveDuplicates(paths),
			}

			if err := ctrl.Upload(opts); err != nil {
				//logger.Error(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&uploadEncryptKey, "encrypt_key", "k", "", "encryption key (required)")
	_ = uploadCmd.MarkFlagRequired("encrypt_key")
	uploadCmd.Flags().StringVarP(&uploadPrefix, "prefix", "p", "", `prefix path to add to the file key (default "")`)
}
