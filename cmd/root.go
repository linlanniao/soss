/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
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
	//	Long: `A longer description that spans multiple lines and likely contains
	//examples and usage of using your application. For example:
	//
	//Cobra is a CLI library for Go that empowers applications.
	//This application is a tool to generate the needed files
	//to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	config   *controller.Config
	endpoint string
	bucket   string
	ctrl     *controller.Controller
	logger   *slog.Logger
)

func init() {
	initLogger()
	initConfig()
	initController()

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", config.Endpoint, "endpoint")
	rootCmd.PersistentFlags().StringVarP(&bucket, "bucket", "b", config.Bucket, "bucket")
}

func initConfig() {
	var err error
	config, err = controller.TryGetConfig()
	if err != nil {
		logger.Error(
			"failed to load the config file", "err", err.Error())
		os.Exit(1)
	}
}

func initLogger() {
	logger = log.DefaultConsoleLogger()
}

func initController() {
	s3Client := ossclient.NewClient(config.Endpoint, config.AccessKey, config.SecretKey)
	fileHandler := filehandler.NewFileHandler()
	ctrl = controller.NewController(config, s3Client, fileHandler, controller.WithLogger(logger))
}
