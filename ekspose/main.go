package main

import (
	// "/lhome/kcharan/kubernetes_scaling/client-go/ekspose/customcontroller"
	"flag"
	"fmt"
	"time"

	"github.com/kama001/client-go/ekspose/customctrlwithqueue"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/lhome/kcharan/.kube/config", "kubeconfig location")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("error building from flags %s", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error getting incluster config %s", err.Error())
		}
	}
	clientset, _ := kubernetes.NewForConfig(config)
	ch := make(chan struct{})
	ns := "default"
	informers := informers.NewSharedInformerFactoryWithOptions(clientset, 10*time.Minute, informers.WithNamespace(ns))
	c := customctrlwithqueue.NewController(*clientset, informers.Apps().V1().Deployments())
	informers.Start(ch)
	c.Run(ch)
}
