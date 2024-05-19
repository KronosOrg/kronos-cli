package utils

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/KronosOrg/kronos-cli/cmd/structs"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetSuccessMessage(spec, action, name string) string {
	return "\n*************************** SUCCESS *************************** \nKronosApp: " + name + " is now " + action + " " + spec + "!\n\n***************************************************************\n"
}

func GetWarningMessage(spec, action, name string) string {
	return "\n*************************** WARNING *************************** \nKronosApp: " + name + " is already " + action + " " + spec + "!\n\n***************************************************************\n"
}

func GetCrdApiUrl(name, namespace string) string {
	namespaceParam := ""
	if namespace != "" {
		namespaceParam = fmt.Sprintf("/namespaces/%s", namespace)
	}
	nameParam := ""
	if name != "" {
		nameParam = name
	}
	crdApiUrl := fmt.Sprintf("/apis/core.wecraft.tn/v1alpha1%s/kronosapps/%s", namespaceParam, nameParam)
	return crdApiUrl
}

func GetFlagNames(cmd *cobra.Command) ([]string, error) {
	flags := make([]string, 3)
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return flags, err
	}
	flags[1] = namespace
	matchRegex, err := cmd.Flags().GetString("match-regex")
	if err != nil {
		return flags, err
	}
	if matchRegex != "" {
		flags[0] = matchRegex
		return flags, nil
	}
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return flags, err
	}
	flags[2] = name

	return flags, nil
}
func InitializeClientConfig() (error, *kubernetes.Clientset) {
	var kubeconfig string

	kubeconfig = os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		fmt.Println("Kubeconfig path not found in KUBECONFIG environment variable.")
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		_, err := os.Stat(kubeconfig)
		if os.IsNotExist(err) {
			fmt.Println("Kubeconfig file not found in the default location ~/.kube/config.")
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter path to kubeconfig file: ")
			path, _ := reader.ReadString('\n')
			kubeconfig = strings.TrimSpace(path)
			fmt.Printf("Using kubeconfig file: %s\n", kubeconfig)
		}
	}

	// Build Kubernetes client configuration
	var config *rest.Config
	var err error
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err, nil
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err, nil
	}
	return nil, clientset
}
func GetKronosAppByName(clientset *kubernetes.Clientset, crdApi string) (error, structs.KronosApp) {
	sd := structs.KronosApp{}
	crd, err := clientset.RESTClient().Get().AbsPath(crdApi).DoRaw(context.TODO())
	if err != nil {
		return err, structs.KronosApp{}
	}
	if err := json.Unmarshal(crd, &sd); err != nil {
		return err, structs.KronosApp{}
	}
	return nil, sd
}

func PerformingActionOnSpec(clientset *kubernetes.Clientset, sd *structs.KronosApp, crdApi, spec, action string) error {
	switch spec {
	case "wake":
		{
			switch action {
			case "on":
				{
					sd.Spec.ForceWake = true
				}
			case "off":
				{
					sd.Spec.ForceWake = false
				}
			}
		}
	case "sleep":
		{
			switch action {
			case "on":
				{
					sd.Spec.ForceSleep = true
				}
			case "off":
				{
					sd.Spec.ForceSleep = false
				}
			}
		}
	}

	sdBytes, err := json.Marshal(sd)
	if err != nil {
		return err
	}
	_, err = clientset.RESTClient().Put().AbsPath(crdApi).Body(sdBytes).DoRaw(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

func GetAllKronosApps(clientset *kubernetes.Clientset, namespace string) (error, structs.KronosAppList) {
	kl := structs.KronosAppList{}
	apiUrl := GetCrdApiUrl("", namespace)
	kl_bytes, err := clientset.RESTClient().Get().AbsPath(apiUrl).DoRaw(context.TODO())
	if err != nil {
		return err, structs.KronosAppList{}
	}
	if err := json.Unmarshal(kl_bytes, &kl); err != nil {
		return err, structs.KronosAppList{}
	}
	return nil, kl
}

func GetKronosAppsByPattern(list structs.KronosAppList, regex regexp.Regexp, namespace string) (error, structs.KronosAppList) {
	var filteredList structs.KronosAppList
	for _, kronosapp := range list.Items {
		if regex.MatchString(kronosapp.Name) {
			filteredList.Items = append(filteredList.Items, kronosapp)
		}
	}
	return nil, filteredList
}

func GetKronosAppsNames(list structs.KronosAppList) []string {
	kronosappsNames := make([]string, len(list.Items))
	for index, kronosapp := range list.Items {
		kronosappsNames[index] = kronosapp.Name
	}
	return kronosappsNames
}

func DisplayAction(name, namespace string, kronosapps structs.KronosAppList) {
	kronosappsNames := GetKronosAppsNames(kronosapps)
	var namespaceInfo string
	if namespace != "" {
		namespaceInfo = " in namespace " + namespace
	}
	if len(kronosappsNames) != 1 {
		counter := 5
		fmt.Printf("Activating ForceWake on the following KronosApps%s: \n", namespaceInfo)
		if counter > len(kronosappsNames) {
			counter = len(kronosappsNames)
		}
		for i := 0; i < counter; i++ {
			fmt.Printf("- %s\n", kronosappsNames[i])
		}
		if counter < len(kronosappsNames) {
			fmt.Printf("...more(%d)\n", len(kronosappsNames)-counter)
		}
	} else {
		DisplayActionByName(kronosappsNames[0], namespace)
	}
}

func DisplayActionByName(name, namespace string) {
	fmt.Printf("Activating ForceWake on KronosApp %s in namespace %s \n", name, namespace)
}

func CheckIfListIsEmpty(list structs.KronosAppList, namespace string) {
	plural := ""
	if namespace == "" {
		namespace = "all"
		plural = "s"
	}
	if len(list.Items) == 0 {
		fmt.Printf("No resources found in %s namespace%s.\n", namespace, plural)
		os.Exit(0)
	}
}

func ApplyActionOnSpecByPattern(clientset *kubernetes.Clientset, regex regexp.Regexp, namespace, spec, action string) error {
	err, allKronosApps := GetAllKronosApps(clientset, namespace)
	if err != nil {
		return err
	}

	CheckIfListIsEmpty(allKronosApps, namespace)

	err, targetKronosApps := GetKronosAppsByPattern(allKronosApps, regex, namespace)
	if err != nil {
		return err
	}

	CheckIfListIsEmpty(targetKronosApps, namespace)

	DisplayAction("", namespace, targetKronosApps)

	var kronosappsUnderEffect []string
	var filteredKronosApps structs.KronosAppList
	for _, targetKronosApp := range targetKronosApps.Items {
		kronosAppUnderEffect, kronosApp := CheckIfActionEffectExist(targetKronosApp, spec, action)
		if kronosAppUnderEffect != "" {
			kronosappsUnderEffect = append(kronosappsUnderEffect, kronosAppUnderEffect)
		} else {
			filteredKronosApps.Items = append(filteredKronosApps.Items, *kronosApp)
		}
	}

	if len(kronosappsUnderEffect) == len(targetKronosApps.Items) {
		return fmt.Errorf("All targeted resources are already %s %s.", action, spec)
	}

	if len(kronosappsUnderEffect) != 0 {
		DisplayUnchangedKronosApps(kronosappsUnderEffect, spec, action)
	}

	apiUrl := GetCrdApiUrl("", namespace)
	var failedActions []string
	for _, kronosapp := range filteredKronosApps.Items {
		err = PerformingActionOnSpec(clientset, &kronosapp, apiUrl, spec, action)
		if err != nil {
			failedActions = append(failedActions, kronosapp.Name)
		}
	}
	if len(failedActions) != 0 {
		err = DisplayFailedActions(failedActions, filteredKronosApps, spec, action)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckIfActionEffectExist(kronosapp structs.KronosApp, spec, action string) (string, *structs.KronosApp) {
	switch spec {
	case "wake":
		{
			switch action {
			case "on":
				{
					if kronosapp.Spec.ForceWake {
						return kronosapp.Name, nil
					} else {
						return "", &kronosapp
					}
				}
			case "off":
				{
					if !kronosapp.Spec.ForceWake {
						return kronosapp.Name, nil
					} else {
						return "", &kronosapp
					}
				}
			}
		}
	case "sleep":
		{
			switch action {
			case "on":
				{
					if kronosapp.Spec.ForceSleep {
						return kronosapp.Name, nil
					} else {
						return "", &kronosapp
					}
				}
			case "off":
				{
					if !kronosapp.Spec.ForceSleep {
						return kronosapp.Name, nil
					} else {
						return "", &kronosapp
					}
				}
			}
		}
	}
	return "", nil
}

func DisplayUnchangedKronosApps(namesList []string, spec, action string) {
	fmt.Printf("The following resources are already %s %s.\n", action, spec)
	for _, itemName := range namesList {
		fmt.Printf("- %s\n", itemName)
	}
}

func DisplayUnchangedKronosApp(name, spec, action string) error {
	return fmt.Errorf("The following resource is already %s %s: KronosApp: %s\n", action, spec, name)
}

func DisplayActionError(spec, action, name string) error {
	var verb string
	if action == "on" {
		verb = "enabling"
	} else if action == "off" {
		verb = "disabling"
	}
	if name != "" {
		return fmt.Errorf("Error %s %s on targeted resource: KronosApp: %s.", verb, spec, name)
	}
	return fmt.Errorf("Error %s %s on targeted resources.", verb, spec)
}

func DisplayFailedActions(failedActions []string, filteredKronosApps structs.KronosAppList, spec, action string) error {
	if len(failedActions) == len(filteredKronosApps.Items) {
		return DisplayActionError(spec, action, "")
	}
	fmt.Println("Error applying changes on the following resources:")
	for _, failedAction := range failedActions {
		fmt.Printf("- %s\n", failedAction)
	}
	return nil
}

func ApplyActionOnSpecByName(clientset *kubernetes.Clientset, name, namespace, spec, action string) error {
	// var kronosappList structs.KronosAppList
	apiUrl := GetCrdApiUrl(name, namespace)
	err, kronosapp := GetKronosAppByName(clientset, apiUrl)
	if err != nil {
		return err
	}
	// kronosappList.Items = append(kronosappList.Items, kronosapp)
	DisplayActionByName(name, namespace)
	// failedActions, filteredKronosApps := CheckIfActionEffectExist(kronosappList, spec, action)
	kronosAppUnderEffect, kronosApp := CheckIfActionEffectExist(kronosapp, spec, action)
	if kronosAppUnderEffect != "" {
		return DisplayUnchangedKronosApp(name, spec, action)
	}
	err = PerformingActionOnSpec(clientset, kronosApp, apiUrl, spec, action)
	if err != nil {
		return DisplayActionError(spec, action, name)
	}
	return nil
}
