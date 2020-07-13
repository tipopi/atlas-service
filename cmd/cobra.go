package cmd


import (
	"github.com/spf13/cobra"
	"os"
	"atlus-service/cmd/api"
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
	//todo:设置自杀机制
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
