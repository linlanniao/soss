package cmd

import (
	"log/slog"
	"os"

	"github.com/linlanniao/soss/internal/controller"
	"github.com/linlanniao/soss/internal/filehandler"
	"github.com/linlanniao/soss/internal/s3clients/ossclient"
	"github.com/linlanniao/soss/pkg/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "soss",
	Short: "SOSS: Secure Object Storage Service",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	config       *controller.Config
	endpoint     string
	bucket       string
	ctrl         *controller.Controller
	logger       *slog.Logger
	s3ClientType string
)

const s3ClientTypeDefault = "oss"

func init() {
	initLogger()
	initConfig()
	initController()

	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", config.Endpoint, "endpoint")
	rootCmd.PersistentFlags().StringVarP(&bucket, "bucket", "b", config.Bucket, "bucket")
	if config.ClientType == "" {
		config.ClientType = s3ClientTypeDefault
	}
	rootCmd.PersistentFlags().StringVarP(&s3ClientType, "s3_type", "s", config.ClientType, "s3 service type")
}

func initConfig() {
	var err error
	config, err = controller.NewConfig()
	if err != nil {
		logger.Error(
			"failed to init config", "err", err.Error())
		os.Exit(1)
	}

	// set default clientType
	if config.ClientType == "" {
		config.ClientType = s3ClientType
	}
}

func initLogger() {
	logger = log.DefaultConsoleLogger()
}

func initController() {
	ossClient := ossclient.NewClient(config.Endpoint, config.AccessKey, config.SecretKey)
	fileHandler := filehandler.NewFileHandler()
	ctrl = controller.NewController(
		controller.WithBucket(config.Bucket),
		controller.WithEndpoint(config.Endpoint),
		controller.WithS3Client(controller.S3ClientTypeOSS, ossClient),
		controller.WithFileHandler(fileHandler),
		controller.WithLogger(logger),
	)
}
