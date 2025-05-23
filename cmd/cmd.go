package cmd

import (
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:     "mxcat",
	Version: "1.0.0",
	Short:   "Audit tool for SMTP stack",
}

func Execute() error {
	return cmd.Execute()
}
