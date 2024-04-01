package cmd

import (
	"os"

	"github.com/linlanniao/soss/internal/secret"
	"github.com/spf13/cobra"
)

var (
	secretForce bool
	// secretCmd represents the secret command
	secretCmd = &cobra.Command{
		Use:     "secret",
		Short:   "generate secret file",
		Aliases: []string{"sec"},
		Run: func(cmd *cobra.Command, _ []string) {
			s := secret.GenerateSecret()
			err := s.Save(secretForce)
			if err != nil {
				logger.Error("error generating secret", "err", err.Error())
				os.Exit(1)
			}
			logger.Info("secret generated", "filepath", s.Path())
		},
	}
)

func init() {
	rootCmd.AddCommand(secretCmd)
	secretCmd.Flags().BoolVarP(&secretForce, "force", "f", false, "force overwrite secret file")
}
