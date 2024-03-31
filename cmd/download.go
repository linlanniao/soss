package cmd

import (
	"os"

	"github.com/linlanniao/soss/internal/controller"
	"github.com/linlanniao/soss/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	downloadDecryptKey string
	downloadOutputDir  string

	// downloadCmd represents the download command
	downloadCmd = &cobra.Command{
		Use:     "download files [files ...]",
		Short:   "Download files from s3Service and decrypt it",
		Aliases: []string{"down", "d"},
		Run: func(cmd *cobra.Command, keys []string) {
			cType := controller.S3ClientType(s3ClientType)
			if err := cType.Validate(); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			opts := controller.DownloadOptions{
				S3ClientType: cType,
				Endpoint:     endpoint,
				Bucket:       bucket,
				OutputDir:    downloadOutputDir,
				DecryptKey:   downloadDecryptKey,
				S3keys:       utils.RemoveDuplicates(keys),
			}

			if err := ctrl.Download(opts); err != nil {
				//logger.Error(err.Error())
				os.Exit(1)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&downloadDecryptKey, "decrypt_key", "k", "", "decryption key (required)")
	_ = downloadCmd.MarkFlagRequired("decrypt_key")
	downloadCmd.Flags().StringVarP(&downloadOutputDir, "output_dir", "o", "./download", `output directory`)
}
