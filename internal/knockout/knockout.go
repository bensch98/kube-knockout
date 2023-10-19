package knockout

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// TerminateNamespace finalizes a Kubernetes namespace by removing finalizers. The namespace will terminate then even if its status is 'Terminating'.
func TerminateNamespace(namespace, kubeconfig string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	ctx := context.TODO()

	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		ns, err := clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
		if err != nil {
			return err
		}

		ns.Finalizers = nil
		_, err = clientset.CoreV1().Namespaces().Update(ctx, ns, metav1.UpdateOptions{})
		return err
	})
	if err != nil {
		return err
	}

	err = clientset.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}
