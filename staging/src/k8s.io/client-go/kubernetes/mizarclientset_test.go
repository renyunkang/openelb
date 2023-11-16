package kubernetes

import (
	"context"
	"fmt"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestMizarClientSet(t *testing.T) {
	var err error
	var config *rest.Config
	var kubeconfig = "/root/.kube/config"

	if config, err = rest.InClusterConfig(); err != nil {
		if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
			panic(err.Error())
		}
	}

	mizarClientset, err := NewMizarClientsetForConfig(config)
	list, err := mizarClientset.Cluster(context.TODO(), "planet-cluster-0").coreV1.Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
		return
	}

	items := list.Items
	for i := range items {
		fmt.Printf("%v \n", items[i])
	}
}
