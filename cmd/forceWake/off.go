/*
Copyright © 2024 Ismail Abdelkefi abdelkefi.ismail@pm.me
*/
package forceWake

import (
	"fmt"
	"os"

	"github.com/KronosOrg/kronos-cli/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	resourceNameOff      string
	resourceNamespaceOff string
)

// offCmd represents the off command
var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Deactivate Force Wake on your KronosApp",
	Long: `Use this command to deactivate Force Wake on a KronosApp:
	
Example:
$ kronos-cli forceWake off --name=my-kronosapp --namespace=my-namespace`,
	Run: func(cmd *cobra.Command, args []string) {
		name, namespace, err := utils.GetFlagNames(cmd)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Deactivating ForceWake on KronosApp: name=%s in namespace=%s \n", name, namespace)
		client := utils.InitializeClientConfig()
		crdApi := utils.GetCrdApiUrl(name, namespace)
		err, sd := utils.GetKronosAppByName(client, crdApi)
		if err != nil {
			fmt.Println(err)
		}
		if !sd.Spec.ForceWake {
			fmt.Printf(utils.GetWarningMessage("ForceWake", "off", name))
			os.Exit(1)
		}
		err = utils.PerformingActionOnSpec(client, &sd, crdApi, "wake", "off")
		if err != nil {
			fmt.Println("ERROR ", err)
			os.Exit(1)
		}
		fmt.Println(utils.GetSuccessMessage("ForceWake", "off", name))
	},
}

func init() {
	offCmd.Flags().StringVarP(&resourceNameOff, "name", "n", "", "The KronosApp name you want to modify")
	offCmd.Flags().StringVarP(&resourceNamespaceOff, "namespace", "", "", "The KronosApp namespace you want to modify")

	if err := offCmd.MarkFlagRequired("name"); err != nil {
		fmt.Println(err)
	}
	if err := offCmd.MarkFlagRequired("namespace"); err != nil {
		fmt.Println(err)
	}

	ForceWakeCmd.AddCommand(offCmd)
}
