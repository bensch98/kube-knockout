package main

import (
	"fmt"
	"os"

	"github.com/bensch98/kube-knockout/internal/knockout"
	"github.com/spf13/pflag"
)

func main() {
	var namespace string
	var kubeconfig string

	pflag.StringVarP(&namespace, "namespace", "n", "", "Namespace to finish")
	pflag.StringVarP(&kubeconfig, "kubeconfig", "k", "", "Path to the kubeconfig file")

	pflag.Parse()

	if namespace == "" {
		fmt.Println("Error: Please specify a namespace useing the '--namespace' flag.")
		os.Exit(1)
	}

	err := knockout.TerminateNamespace(namespace, kubeconfig)
	if err != nil {
		fmt.Printf("Error finishing namespace '%s': %v\n", namespace, err)
		os.Exit(1)
	}

	fmt.Printf("Namespace '%s' finished successfully.\n", namespace)
}
