/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package forceSleep

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/KronosOrg/kronos-cli/cmd/structs"
	"github.com/KronosOrg/kronos-cli/cmd/utils"
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
			fmt.Println(err)
		}
		fmt.Printf("Deactivating ForceSleep on KronosApp: name=%s in namespace=%s \n", name, namespace)
		client := utils.InitializeClientConfig()
		sd := structs.KronosApp{}
		crdApi := utils.GetCrdApiUrl(name, namespace)
		sd = utils.GetKronosAppByName(client, crdApi)
		ok := utils.CheckForceSleep(&sd)
		if !ok {
			fmt.Printf("\n*************************** WARNING *************************** \nKronosApp: %s is already off ForceSleep!\n\n***************************************************************\n", name)
			os.Exit(1)
		}
		utils.DeactivatingForceSleep(client, &sd, crdApi, name)
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
