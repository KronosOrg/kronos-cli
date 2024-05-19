/*
Copyright © 2024 Ismail Abdelkefi abdelkefi.ismail@pm.me
*/
package forceWake

import (
	"fmt"
	"os"
	"regexp"

	"github.com/KronosOrg/kronos-cli/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	resourceNameOff      string
	resourceNamespaceOff string
	matchRegexOff        string
)

// offCmd represents the off command
var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Deactivate Force Wake on your KronosApp",
	Long: `Use this command to deactivate Force Wake on a KronosApp:
	
Example:
$ kronos-cli forceWake off --name=my-kronosapp --namespace=my-namespace`,
	Run: func(cmd *cobra.Command, args []string) {
		spec := "ForceWake"
		action := "off"

		flags, err := utils.GetFlagNames(cmd)
		if err != nil {
			fmt.Println(err)
		}

		regexPattern := flags[0]
		namespace := flags[1]
		name := flags[2]

		err, client := utils.InitializeClientConfig()
		if err != nil {
			fmt.Println("ERROR ", err)
			os.Exit(1)
		}

		if regexPattern != "" {
			regex := regexp.MustCompile(regexPattern)
			err = utils.ApplyActionOnSpecByPattern(client, *regex, namespace, spec, action)
			if err != nil {
				fmt.Println("ERROR ", err)
				os.Exit(1)
			}
			os.Exit(0)
		} else {
			err = utils.ApplyActionOnSpecByName(client, name, namespace, spec, action)
			if err != nil {
				fmt.Println("ERROR ", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	},
}

func init() {
	offCmd.Flags().StringVarP(&resourceNameOff, "name", "n", "", "The KronosApp name you want to modify")
	offCmd.Flags().StringVarP(&resourceNamespaceOff, "namespace", "", "", "The KronosApp namespace you want to modify")
	offCmd.Flags().StringVarP(&matchRegexOff, "match-regex", "", "", "Pattern, applied on name, used to regroup KronosApps you want to modify")

	offCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if matchRegexOff != "" {
			if resourceNameOff != "" {
				return fmt.Errorf("the --match-regex flag cannot be used with the --name flag")
			}
		} else {
			if resourceNameOff == "" || resourceNamespaceOff == "" {
				return fmt.Errorf("both --name and --namespace flags are required unless --match-regex is used")
			}
		}
		return nil
	}

	ForceWakeCmd.AddCommand(offCmd)
}
