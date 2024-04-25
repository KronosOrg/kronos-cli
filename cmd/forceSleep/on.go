/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package forceSleep

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.infra.wecraft.tn/wecraft/automation/ifra/kronos-cli/cmd/structs"
	"gitlab.infra.wecraft.tn/wecraft/automation/ifra/kronos-cli/cmd/utils"
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
			fmt.Println(err)
		}
		fmt.Printf("Activating ForceSleep on KronosApp: name=%s in namespace=%s \n", name, namespace)
		client := utils.InitializeClientConfig()
		sd := structs.KronosApp{}
		crdApi := utils.GetCrdApiUrl(name, namespace)
		sd = utils.GetKronosAppByName(client, crdApi)
		ok := utils.CheckForceWake(&sd)
		if ok {
			fmt.Println("ForceWake is ON! We cannot proceed with your request.")
			os.Exit(1)
		}
		ok = utils.CheckForceSleep(&sd)
		if ok {
			fmt.Printf("\n*************************** WARNING *************************** \nKronosApp: %s is already on ForceSleep!\n\n***************************************************************\n", name)
			os.Exit(1)
		}
		utils.ActivatingForceSleep(client, &sd, crdApi, name)
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
