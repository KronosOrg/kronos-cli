/*
Copyright © 2024 Ismail Abdelkefi abdelkefi.ismail@pm.me
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func printVersion() {
	cliVersion := "1.0.0"
	goVersion := "1.22.0"
	fmt.Print(`    __ __ ____  ____  _   ______  _____
   / //_// __ \/ __ \/ | / / __ \/ ___/
  / ,<  / /_/ / / / /  |/ / / / /\__ \ 
 / /| |/ _, _/ /_/ / /|  / /_/ /___/ / 
/_/ |_/_/ |_|\____/_/ |_/\____//____/  

CLI Version:` + " " + cliVersion + "\n" + "GO version: " + goVersion + "\n")
}

// versionCmd represents the version command
var customVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints out the version of kronos-cli.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {

}
