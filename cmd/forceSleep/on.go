/*
Copyright © 2024 Ismail Abdelkefi abdelkefi.ismail@pm.me
*/
package forceSleep

import (
	"fmt"
	"os"
	"regexp"

	"github.com/KronosOrg/kronos-cli/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	resourceNameOn      string
	resourceNamespaceOn string
	matchRegexOn        string
)

// onCmd represents the on command
var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Activate Force Sleep on your KronosApp",
	Long: `Use this command to activate Force Sleep on a KronosApp:

Example:
$ kronos-cli forceSleep on --name=my-kronosap --namespace=my-namespace`,

	Run: func(cmd *cobra.Command, args []string) {
		spec := "ForceSleep"
		action := "on"

		flags, err := utils.GetFlagNames(cmd)
		if err != nil {
			fmt.Println(err)
		}

		regexPattern := flags[0]
		namespace := flags[1]
		name := flags[2]

		err, client := utils.InitializeClientConfig()
		if err != nil {
			fmt.Println("ERROR", err)
			os.Exit(1)
		}

		if regexPattern != "" {
			regex := regexp.MustCompile(regexPattern)
			err = utils.ApplyActionOnSpecByPattern(client, *regex, namespace, spec, action)
			if err != nil {
				fmt.Println("ERROR", err)
				os.Exit(1)
			}
			os.Exit(0)
		} else {
			err = utils.ApplyActionOnSpecByName(client, name, namespace, spec, action)
			if err != nil {
				fmt.Println("ERROR", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	},
}

func init() {
	onCmd.Flags().StringVarP(&resourceNameOn, "name", "n", "", "The KronosApp name you want to modify")
	onCmd.Flags().StringVarP(&resourceNamespaceOn, "namespace", "", "", "The KronosApp namespace you want to modify")
	onCmd.Flags().StringVarP(&matchRegexOn, "match-regex", "", "", "Pattern, applied on name, used to regroup KronosApps you want to modify")

	onCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if matchRegexOn != "" {
			if resourceNameOn != "" {
				return fmt.Errorf("the --match-regex flag cannot be used with the --name flag")
			}
		} else {
			if resourceNameOn == "" || resourceNamespaceOn == "" {
				return fmt.Errorf("both --name and --namespace flags are required unless --match-regex is used")
			}
		}
		return nil
	}
	ForceSleepCmd.AddCommand(onCmd)

}
