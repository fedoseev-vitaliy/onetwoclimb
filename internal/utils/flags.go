package utils

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// todo add test
func BindEnv(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		envVar := strings.ToUpper(f.Name)

		if val := os.Getenv(envVar); val != "" {
			if err := cmd.Flags().Set(f.Name, val); err != nil {
				logrus.WithError(err).Error("failed to set flag")
			}
		}
	})
}
