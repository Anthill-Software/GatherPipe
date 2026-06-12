package cmd

import (
	"fmt"

	"github.com/Anthill-Software/GatherPipe/core"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Afficher la version de GatherPipe",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GatherPipe v%s\n", core.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
