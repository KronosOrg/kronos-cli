package utils

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/KronosOrg/kronos-cli/cmd/structs"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"strings"
)

func GetSuccessMessage(spec, action, name string) string {
	return "\n*************************** SUCCESS *************************** \nKronosApp: " + name + " is now " + action + " " + spec + "!\n\n***************************************************************\n"
}

func GetWarningMessage(spec, action, name string) string {
	return "\n*************************** WARNING *************************** \nKronosApp: " + name + " is already " + action + " "  + spec + "!\n\n***************************************************************\n"
}


func GetCrdApiUrl(name, namespace string) string {
	crdApi := fmt.Sprintf("/apis/core.wecraft.tn/v1alpha1/namespaces/%s/kronosapps/%s", namespace, name)
	return crdApi
}

func GetFlagNames(cmd *cobra.Command) (string, string, error) {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return "", "", err
	}
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return "", "", err
	}
	return name, namespace, nil
}
func InitializeClientConfig() *kubernetes.Clientset {
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
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	return clientset
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
		case "wake": {
			switch action {
				case "on": {
					sd.Spec.ForceWake = true
				}
				case "off": {
					sd.Spec.ForceWake = false
				}
			}
		}
		case "sleep": {
			switch action {
				case "on": {
					sd.Spec.ForceSleep = true
				}
				case "off": {
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