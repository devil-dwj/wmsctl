package main

import (
	"github.com/devil-dwj/wmsctl/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wmsctl",
	Short: "Wmsctl: toolkit for wms",
	Long:  "Wmsctl: toolkit for wms",
}

func init() {
	rootCmd.AddCommand(cmd.Init)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
