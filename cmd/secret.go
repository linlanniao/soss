package cmd

import (
	"os"

	"github.com/linlanniao/soss/internal/secret"
	"github.com/spf13/cobra"
)

var (
	secretReplace bool
	// secretCmd represents the secret command
	secretCmd = &cobra.Command{
		Use:     "secret",
		Short:   "generate secret file",
		Aliases: []string{"sec"},
		Run: func(cmd *cobra.Command, _ []string) {
			s := secret.GenerateSecret()
			var err error
			var newerSec *secret.Secret
			if secretReplace {
				newerSec, err = s.Replace()
			} else {
				err = s.Save()
			}

			if err != nil {
				logger.Error("error generating secret", "err", err.Error())
				os.Exit(1)
			}

			if newerSec != nil && newerSec.Path() != s.Path() {
				logger.Info("secret file has been replaced",
					"newer", newerSec.Path(),
					"backup", s.Path())
			} else {
				logger.Info("secret generated", "filepath", s.Path())
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(secretCmd)
	secretCmd.Flags().BoolVarP(&secretReplace, "replace", "r", false, "replace secret file")
}
