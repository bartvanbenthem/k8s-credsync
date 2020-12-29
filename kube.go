package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeCLient struct{}

func (k *KubeCLient) CreateClientSet() *kubernetes.Clientset {
	// When running the binary inside of a pod in a cluster,
	// the kubelet will automatically mount a service account into the container at:
	// /var/run/secrets/kubernetes.io/serviceaccount.
	// It replaces the kubeconfig file and is turned into a rest.Config via the rest.InClusterConfig() method
	config, err := rest.InClusterConfig()
	if err != nil {
		// fallback to kubeconfig
		kubeconfig := filepath.Join("~", ".kube", "config")
		if envvar := os.Getenv("K8S_KUBECONFIG"); len(envvar) > 0 {
			kubeconfig = envvar
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Printf("The kubeconfig cannot be loaded: %v\n", err)
			os.Exit(1)
		}
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	return clientset
}

func (k *KubeCLient) UpdateSecret(c *kubernetes.Clientset, namespace string, secret *v1.Secret) *v1.Secret {
	sec, err := c.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("\nUnable to retreive secret data field from: %v\n", namespace)
		fmt.Printf("%v\n", err)
	}
	return sec
}

func (k *KubeCLient) GetSecret(c *kubernetes.Clientset, namespace, secretname string) *v1.Secret {
	sec, err := c.CoreV1().Secrets(namespace).Get(context.TODO(), secretname, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("\nUnable to retreive secret data field from: %v\n", namespace)
		fmt.Printf("%v\n", err)
	}
	return sec
}

func (k *KubeCLient) GetSecretData(c *kubernetes.Clientset, namespace, secretname, datafield string) []byte {
	sec, err := c.CoreV1().Secrets(namespace).Get(context.TODO(), secretname, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("\nUnable to retreive secret data field from: %v\n", namespace)
		fmt.Printf("%v\n", err)
	}
	return sec.Data[datafield]
}

func (k *KubeCLient) GetAllNamespaces(c *kubernetes.Clientset) []string {
	var namespaces []string
	ns, err := c.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("\nError listing namespaces: %v\n", err)
	}
	for _, n := range ns.Items {
		namespaces = append(namespaces, n.Name)
	}
	return namespaces
}
