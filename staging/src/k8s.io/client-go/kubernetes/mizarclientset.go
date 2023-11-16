package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	kubeconfig                = "-kubeconfig"
	kubeconfigSecretNamespace = "default"
)

type MizarClientset struct {
	*Clientset
}

// NewMizarClientForClientset creates a new MizarClientset for the star-cluster Clientset.
func NewMizarClientForClientset(clientset *Clientset) *MizarClientset {
	return &MizarClientset{
		clientset,
	}
}

// NewMizarClientsetForConfig creates a new MizarClientset for the given config.
func NewMizarClientsetForConfig(config *rest.Config) (*MizarClientset, error) {
	clientset, err := NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &MizarClientset{
		clientset,
	}, nil
}

func (c MizarClientset) Cluster(context context.Context, clusterName string) *Clientset {
	// use start-cluster kubeconfig get planet-cluster kubeconfig by secret
	secret, err := c.coreV1.Secrets(kubeconfigSecretNamespace).Get(context, clusterName+kubeconfig, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
		return nil
	}

	// get planet-cluster RestConfig
	config, err := clientcmd.RESTConfigFromKubeConfig(secret.Data["value"])
	if err != nil {
		panic(err.Error())
		return nil
	}

	// get planet-cluster Clientset
	clientset, err := NewForConfig(config)
	if err != nil {
		panic(err.Error())
		return nil
	}
	return clientset
}
