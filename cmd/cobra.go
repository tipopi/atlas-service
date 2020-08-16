package cmd

import (
	"atlas-service/cmd/api"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:               "atlas",
	Short:             "RTX ON",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Long:              `Start atlas`,
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
}

func init() {
	//嵌套指令
	rootCmd.AddCommand(api.StartCmd)
}

//Execute : run commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
