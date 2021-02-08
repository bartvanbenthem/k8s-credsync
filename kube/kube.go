package kube

import (
	"context"
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
			log.Printf("The kubeconfig cannot be loaded: %v\n", err)
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
		log.Printf("\"%v\" Secret not found in \"%v\" Namespace\n",
			secretname, namespace)
		log.Printf("%v\n", err)
	}
	return sec
}

func (k *KubeCLient) GetSecretData(c *kubernetes.Clientset, namespace, secretname, datafield string) []byte {
	sec, err := c.CoreV1().Secrets(namespace).Get(context.TODO(), secretname, metav1.GetOptions{})
	if err != nil {
		log.Printf("Data-field from \"%v\" Secret not found in \"%v\" Namespace\n",
			secretname, namespace)
		log.Printf("%v\n", err)
	}
	return sec.Data[datafield]
}

func (k *KubeCLient) UpdateSecret(c *kubernetes.Clientset, namespace string, secret *v1.Secret) *v1.Secret {
	sec, err := c.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Error Updating \"%v\" Secret in \"%v\" Namespace \n",
			secret.Name, namespace)
		log.Printf("%v\n", err)
	}
	log.Printf("Updated \"%v\" Secret in \"%v\" Namespace\n",
		secret.Name, namespace)
	return sec
}

func (k *KubeCLient) CreateSecret(c *kubernetes.Clientset, namespace string, secret *v1.Secret) *v1.Secret {
	sec, err := c.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Error Creating \"%v\" Secret in \"%v\" Namespace \n",
			secret.Name, namespace)
		log.Printf("%v\n", err)
	}
	log.Printf("Created \"%v\" Secret in \"%v\" Namespace\n",
		secret.Name, namespace)
	return sec
}

func (k *KubeCLient) DeleteSecret(c *kubernetes.Clientset, namespace string, secret *v1.Secret) {
	deletePolicy := metav1.DeletePropagationForeground
	err := c.CoreV1().Secrets(namespace).Delete(context.TODO(), secret.Name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		log.Printf("Error Deleting \"%v\" in \"%v\" Namespace \n", secret.Name, namespace)
		log.Printf("%v\n", err)
	}
	log.Printf("Deleted \"%v\" Secret in \"%v\" Namespace\n",
		secret.Name, namespace)
}

func (k *KubeCLient) GetConfigmap(c *kubernetes.Clientset, namespace, cmname string) *v1.ConfigMap {
	cm, err := c.CoreV1().ConfigMaps(namespace).Get(context.TODO(), cmname, metav1.GetOptions{})
	if err != nil {
		log.Printf("\"%v\" Configmap not found in \"%v\" Namespace\n",
			cmname, namespace)
		log.Printf("%v\n", err)
	}
	return cm
}

func (k *KubeCLient) GetAllNamespaceNames(c *kubernetes.Clientset) []string {
	var namespaces []string
	ns, err := c.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error Getting Namespaces: %v\n", err)
	}
	for _, n := range ns.Items {
		namespaces = append(namespaces, n.Name)
	}
	return namespaces
}

func (k *KubeCLient) GetAllPodNames(c *kubernetes.Clientset, namespace string) []string {
	var pods []string
	ns, err := c.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error Getting Pods: %v\n", err)
	}
	for _, p := range ns.Items {
		pods = append(pods, p.Name)
	}
	return pods
}

func (k *KubeCLient) DeletePod(c *kubernetes.Clientset, namespace, podname string) {
	err := c.CoreV1().Pods(namespace).Delete(context.TODO(), podname, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Error Deleting \"%v\" Pod\n", err)
	}
	log.Printf("Deleted \"%v\" Pod in \"%v\" Namespace\n ",
		podname, namespace)
}
