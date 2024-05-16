/*
Copyright © 2024 Ismail Abdelkefi abdelkefi.ismail@pm.me
*/
package forceSleep

import (
	// "fmt"

	"github.com/spf13/cobra"
)

// forceSleepCmd represents the forceSleep command
var ForceSleepCmd = &cobra.Command{
	Use:   "forceSleep",
	Short: "Package forceSleep provides the functionality to activate/deactivate Force Sleep on a KronosApp resource.",
	Long:  `Package forceSleep provides the functionality to activate/deactivate Force Sleep on a KronosApp resource.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {

}
