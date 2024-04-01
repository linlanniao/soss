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
		Use:     "upload files [files ...]",
		Short:   "Encrypt and upload files to s3service",
		Aliases: []string{"up", "u"},
		Run: func(cmd *cobra.Command, paths []string) {
			cType := controller.S3ClientType(s3ClientType)
			if err := cType.Validate(); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			initSecretKey()
			var k string
			if useSecretFile && len(secretKey) > 0 {
				k = secretKey
			} else {
				if len(uploadEncryptKey) == 0 {
					logger.Error("encrypt_key is required")
					os.Exit(1)
				}
				k = uploadEncryptKey
			}

			opts := controller.UploadOptions{
				S3ClientType: cType,
				Endpoint:     endpoint,
				Bucket:       bucket,
				Prefix:       uploadPrefix,
				EncryptKey:   k,
				Paths:        utils.RemoveDuplicates(paths),
			}

			if err := ctrl.Upload(opts); err != nil {
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&uploadEncryptKey, "encrypt_key", "k", "", "encryption key (required)")
	uploadCmd.Flags().StringVarP(&uploadPrefix, "prefix", "p", "", `prefix path to add to the file key (default "")`)
}
