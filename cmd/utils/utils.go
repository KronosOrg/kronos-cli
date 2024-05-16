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
func GetKronosAppByName(clientset *kubernetes.Clientset, crdApi string) structs.KronosApp {
	sd := structs.KronosApp{}
	crd, err := clientset.RESTClient().Get().AbsPath(crdApi).DoRaw(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	if err := json.Unmarshal(crd, &sd); err != nil {
		panic(err)
	}
	return sd
}

func CheckForceWake(sd *structs.KronosApp) bool {
	return sd.Spec.ForceWake
}

func CheckForceSleep(sd *structs.KronosApp) bool {
	return sd.Spec.ForceSleep
}

func ActivatingForceWake(clientset *kubernetes.Clientset, sd *structs.KronosApp, crdApi string, name string) {
	sd.Spec.ForceWake = true
	sdBytes, err := json.Marshal(sd)
	if err != nil {
		panic(err)
	}
	_, err = clientset.RESTClient().Put().AbsPath(crdApi).Body(sdBytes).DoRaw(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n*************************** SUCCESS *************************** \nKronosApp: %s is now on ForceWake!\n\n***************************************************************\n", name)
}

func DeactivatingForceWake(clientset *kubernetes.Clientset, sd *structs.KronosApp, crdApi string, name string) {
	sd.Spec.ForceWake = false
	sdBytes, err := json.Marshal(sd)
	if err != nil {
		panic(err)
	}
	_, err = clientset.RESTClient().Put().AbsPath(crdApi).Body(sdBytes).DoRaw(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n*************************** SUCCESS *************************** \nKronosApp: %s is now off ForceWake!\n\n***************************************************************\n", name)
}

func ActivatingForceSleep(clientset *kubernetes.Clientset, sd *structs.KronosApp, crdApi string, name string) {
	sd.Spec.ForceSleep = true
	sdBytes, err := json.Marshal(sd)
	if err != nil {
		panic(err)
	}
	_, err = clientset.RESTClient().Put().AbsPath(crdApi).Body(sdBytes).DoRaw(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n*************************** SUCCESS *************************** \nKronosApp: %s is now on ForceSleep!\n\n***************************************************************\n", name)
}

func DeactivatingForceSleep(clientset *kubernetes.Clientset, sd *structs.KronosApp, crdApi string, name string) {
	sd.Spec.ForceSleep = false
	sdBytes, err := json.Marshal(sd)
	if err != nil {
		panic(err)
	}
	_, err = clientset.RESTClient().Put().AbsPath(crdApi).Body(sdBytes).DoRaw(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n*************************** SUCCESS *************************** \nKronosApp: %s is now off ForceSleep!\n\n***************************************************************\n", name)
}
