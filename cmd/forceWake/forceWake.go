/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package forceWake

import (
	// "fmt"

	"github.com/spf13/cobra"
)

// forceWakeCmd represents the forceWake command
var ForceWakeCmd = &cobra.Command{
	Use:   "forceWake",
	Short: "Package forceWake provides the functionality to activate/deactivate Force Wake on a KronosApp resource.",
	Long:  `Package forceWake provides the functionality to activate/deactivate Force Wake on a KronosApp resource.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {

}
