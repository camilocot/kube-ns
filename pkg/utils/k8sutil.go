package utils

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	api_v1 "k8s.io/api/core/v1"
	net_v1 "k8s.io/api/networking/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClient returns a k8s clientset to the request from inside of cluster
func GetClient() kubernetes.Interface {
	config, err := rest.InClusterConfig()
	if err != nil {
		logrus.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Can not create kubernetes client: %v", err)
	}

	return clientset
}

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

// GetClientOutOfCluster returns a k8s clientset to the request from outside of cluster
func GetClientOutOfCluster() kubernetes.Interface {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		logrus.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Can not get kubernetes config: %v", err)
	}

	return clientset
}

// GetObjectMetaData returns metadata of a given k8s object
func GetObjectMetaData(obj interface{}) (objectMeta meta_v1.ObjectMeta, err error) {

	switch object := obj.(type) {
	case *api_v1.Namespace:
		objectMeta = object.ObjectMeta
	case *net_v1.NetworkPolicy:
		objectMeta = object.ObjectMeta
	default:
		err = fmt.Errorf("don't know about type %T", object)
	}

	return objectMeta, err
}
