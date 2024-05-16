/*
Copyright © 2024 Ismail Abdelkefi abdelkefi.ismail@pm.me
*/
package forceWake

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/KronosOrg/kronos-cli/cmd/structs"
	"github.com/KronosOrg/kronos-cli/cmd/utils"
)

var (
	resourceNameOn      string
	resourceNamespaceOn string
)

// onCmd represents the on command
var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Activate Force Wake on your KronosApp",
	Long: `Use this command to activate Force Wake on a KronosApp:
	
Example:
$ kronos-cli forceWake on --name=my-kronosapp --namespace=my-namespace`,
	Run: func(cmd *cobra.Command, args []string) {
		name, namespace, err := utils.GetFlagNames(cmd)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Activating ForceWake on KronosApp: name=%s in namespace=%s \n", name, namespace)
		client := utils.InitializeClientConfig()
		sd := structs.KronosApp{}
		crdApi := utils.GetCrdApiUrl(name, namespace)
		sd = utils.GetKronosAppByName(client, crdApi)
		ok := utils.CheckForceSleep(&sd)
		if ok {
			fmt.Println("ForceSleep is ON! We cannot proceed with your request.")
			os.Exit(1)
		}
		ok = utils.CheckForceWake(&sd)
		if ok {
			fmt.Printf("\n*************************** WARNING *************************** \nKronosApp: %s is already on ForceWake!\n\n***************************************************************\n", name)
			os.Exit(1)
		}
		utils.ActivatingForceWake(client, &sd, crdApi, name)
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
	ForceWakeCmd.AddCommand(onCmd)
}
