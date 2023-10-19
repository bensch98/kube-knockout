package knockout

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

// TerminateNamespace finalizes a Kubernetes namespace by removing finalizers. The namespace will terminate then even if its status is 'Terminating'.
func DeleteFinalizers(resourceType, resourceName, namespace, kubeconfig string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	ctx := context.TODO()

	switch resourceType {
	case "ns", "namespace":
		return terminateNamespace(clientset, ctx, resourceName)
	case "pvc", "pv":
		return terminatePersistentVolumeClaim(clientset, ctx, resourceName)
	default:
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

// Terminates a whole namespace.
func terminateNamespace(clientset *kubernetes.Clientset, ctx context.Context, namespace string) error {
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
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

// Terminates a PersistentVolumeClaim.
func terminatePersistentVolumeClaim(clientset *kubernetes.Clientset, ctx context.Context, resourceName string) error {
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		fmt.Println(resourceName)
		pv, err := clientset.CoreV1().PersistentVolumes().Get(ctx, resourceName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		pv.Finalizers = nil
		_, err = clientset.CoreV1().PersistentVolumes().Update(ctx, pv, metav1.UpdateOptions{})
		return err
	})
	if err != nil {
		return err
	}

	err = clientset.CoreV1().PersistentVolumes().Delete(ctx, resourceName, metav1.DeleteOptions{})
	if err != nil {
		return nil
	}

	return nil
}
