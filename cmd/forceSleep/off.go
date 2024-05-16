/*
Copyright © 2024 Ismail Abdelkefi abdelkefi.ismail@pm.me
*/
package forceSleep

import (
	"fmt"
	"github.com/KronosOrg/kronos-cli/cmd/utils"
	"github.com/spf13/cobra"
	"os"
)

var (
	resourceNameOff      string
	resourceNamespaceOff string
)

// offCmd represents the off command
var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Deactivate Force Sleep on your KronosApp",
	Long: `Use this command to deactivate Force Sleep on a KronosApp:

Example:
$ kronos-cli forceSleep off --name=my-kronosapp --namespace=my-namespace`,
	Run: func(cmd *cobra.Command, args []string) {
		name, namespace, err := utils.GetFlagNames(cmd)
		if err != nil {
			fmt.Println("ERROR", err)
		}
		fmt.Printf("Deactivating ForceSleep on KronosApp: name=%s in namespace=%s \n", name, namespace)
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
		if !sd.Spec.ForceSleep {
			fmt.Println(utils.GetWarningMessage("ForceSleep", "off", name))
			os.Exit(0)
		}
		err = utils.PerformingActionOnSpec(client, &sd, crdApi, "sleep", "off")
		if err != nil {
			fmt.Println("ERROR ", err)
			os.Exit(1)
		}
		fmt.Println(utils.GetSuccessMessage("ForceSleep", "off", name))
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

	ForceSleepCmd.AddCommand(offCmd)
}
