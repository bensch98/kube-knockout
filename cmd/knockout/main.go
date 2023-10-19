package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/bensch98/kube-knockout/internal/knockout"
	"github.com/spf13/pflag"
)

func main() {
	var namespace string
	var kubeconfig string

	pflag.StringVarP(&namespace, "namespace", "n", "", "Namespace to finish")
	pflag.StringVarP(&kubeconfig, "kubeconfig", "k", "", "Path to the kubeconfig file")

	pflag.Parse()

	args := pflag.Args()
	if len(args) < 2 {
		fmt.Println("Error: Please specify the resource type and resource name as arguments.")
		os.Exit(1)
	}

	resourceType := args[0]
	resourceName := args[1]

	// Determine the default kubeconfig file path
	if kubeconfig == "" {
		kubeconfig = getDefaultKubeconfigPath()
	}

	err := knockout.DeleteFinalizers(resourceType, resourceName, namespace, kubeconfig)
	if err != nil {
		fmt.Printf("Error terminating resource %s: %v\n", namespace, err)
		os.Exit(1)
	}

	fmt.Printf("%s %s finished successfully.\n", resourceType, resourceName)
}

// Returns the default kubeconfig file path.
func getDefaultKubeconfigPath() string {
	usr, err := user.Current()
	if err != nil {
		return filepath.Join("~", ".kube", "config")
	}
	return filepath.Join(usr.HomeDir, ".kube", "config")
}
