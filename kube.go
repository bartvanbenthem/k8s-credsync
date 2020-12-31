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

func (k *KubeCLient) GetSecret(c *kubernetes.Clientset, namespace, secretname string) *v1.Secret {
	sec, err := c.CoreV1().Secrets(namespace).Get(context.TODO(), secretname, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("\nUnable to retreive secret from: %v\n", namespace)
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

func (k *KubeCLient) UpdateSecret(c *kubernetes.Clientset, namespace string, secret *v1.Secret) *v1.Secret {
	sec, err := c.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("\nError Updating \"%v\" in \"%v\" Namespace \n", secret.Name, namespace)
		fmt.Printf("%v\n", err)
	}
	return sec
}

func (k *KubeCLient) CreateSecret(c *kubernetes.Clientset, namespace string, secret *v1.Secret) *v1.Secret {
	sec, err := c.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("\nError Creating \"%v\" in \"%v\" Namespace \n", secret.Name, namespace)
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("\nCreated \"%v\" Secret.\n", secret.Name)
	return sec
}

func (k *KubeCLient) DeleteSecret(c *kubernetes.Clientset, namespace string, secret *v1.Secret) {
	fmt.Println("Deleting Secret...")
	deletePolicy := metav1.DeletePropagationForeground
	err := c.CoreV1().Secrets(namespace).Delete(context.TODO(), secret.Name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		fmt.Printf("\nError Deleting \"%v\" in \"%v\" Namespace \n", secret.Name, namespace)
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("\nDeleted \"%v\" Secret.\n", secret.Name)
}

func (k *KubeCLient) GetAllNamespaces(c *kubernetes.Clientset) []string {
	var namespaces []string
	ns, err := c.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("\nError Listing Namespaces: %v\n", err)
	}
	for _, n := range ns.Items {
		namespaces = append(namespaces, n.Name)
	}
	return namespaces
}

func (k *KubeCLient) GetAllPods(c *kubernetes.Clientset, namespace string) []string {
	var pods []string
	ns, err := c.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("\nError Listing Pods: %v\n", err)
	}
	for _, p := range ns.Items {
		pods = append(pods, p.Name)
	}
	return pods
}

func (k *KubeCLient) DeletePod(c *kubernetes.Clientset, namespace, podname string) {
	err := c.CoreV1().Pods(namespace).Delete(context.TODO(), podname, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("\nError Deleting Pod: %v\n", err)
	}
	fmt.Println()
	log.Printf("\nPod \"%v\" Deleted\n", podname)
}

func int32Ptr(i int32) *int32 { return &i }
