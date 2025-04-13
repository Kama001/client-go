package main

// return values = rv
// input values = iv
import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/lhome/kcharan/.kube/config", "location to kube Config")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig) // iv: masterUrl, kubeConfigPath, rv:  *restclient.Config, error
	if err != nil {

	}
	// fmt.Println(config) // prints kube config file
	clientset, err := kubernetes.NewForConfig(config) // iv: *restclient.Config, rv: *Clientset, err
	if err != nil {

	}
	// fmt.Printf("%+v", clientset) // prints something aabout client set
	context := context.Background()
	// corev1 := https://pkg.go.dev/k8s.io/client-go/kubernetes#Clientset.CoreV1
	// Pods := https://pkg.go.dev/k8s.io/client-go@v0.32.3/kubernetes/typed/core/v1#PodsGetter
	// list := https://pkg.go.dev/k8s.io/client-go@v0.32.3/kubernetes/typed/core/v1#PodInterface
	// List options := https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#ListOptions
	pods, err := clientset.CoreV1().Pods("default").List(context, metav1.ListOptions{})
	// fmt.Printf("%+v", pods)
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}
	if err != nil {

	}
	deployments, err := clientset.AppsV1().Deployments("default").List(context, metav1.ListOptions{})
	for _, deployment := range deployments.Items {
		fmt.Println(deployment.Name)
	}

}
