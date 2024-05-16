/*
Copyright © 2024 Ismail Abdelkefi abdelkefi.ismail@pm.me
*/
package forceSleep

import (
	"fmt"
	"os"

	"github.com/KronosOrg/kronos-cli/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	resourceNameOn      string
	resourceNamespaceOn string
)

// onCmd represents the on command
var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Activate Force Sleep on your KronosApp",
	Long: `Use this command to activate Force Sleep on a KronosApp:

Example:
$ kronos-cli forceSleep on --name=my-kronosap --namespace=my-namespace`,

	Run: func(cmd *cobra.Command, args []string) {
		name, namespace, err := utils.GetFlagNames(cmd)
		if err != nil {
			fmt.Println("ERROR", err)
		}
		fmt.Printf("Activating ForceSleep on KronosApp: name=%s in namespace=%s \n", name, namespace)
		err, client := utils.InitializeClientConfig()
		if err != nil {
			fmt.Println("ERROR", err)
			os.Exit(1)
		}
		crdApi := utils.GetCrdApiUrl(name, namespace)
		err, sd := utils.GetKronosAppByName(client, crdApi)
		if err != nil {
			fmt.Println("ERROR", err)
			os.Exit(1)
		}
		if sd.Spec.ForceSleep {
			fmt.Println(utils.GetWarningMessage("ForceSleep", "on", name))
			os.Exit(1)
		}
		err = utils.PerformingActionOnSpec(client, &sd, crdApi, "sleep", "on")
		if err != nil {
			fmt.Println("ERROR ", err)
			os.Exit(1)
		}
		fmt.Println(utils.GetSuccessMessage("ForceSleep", "on", name))
	},
}

func init() {
	onCmd.Flags().StringVarP(&resourceNameOn, "name", "n", "", "The KronosApp name you want to modify")
	onCmd.Flags().StringVarP(&resourceNamespaceOn, "namespace", "", "", "The KronosApp namespace you want to modify")

	if err := onCmd.MarkFlagRequired("name"); err != nil {
		fmt.Println(err)
	}
	if err := onCmd.MarkFlagRequired("namespace"); err != nil {
		fmt.Println(err)
	}
	ForceSleepCmd.AddCommand(onCmd)

}
