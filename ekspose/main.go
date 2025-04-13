package main

import (
	// "/lhome/kcharan/kubernetes_scaling/client-go/ekspose/customcontroller"
	"flag"
	"time"

	"github.com/kama001/client-go/ekspose/customcontroller"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/lhome/kcharan/.kube/config", "kubeconfig location")
	config, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	clientset, _ := kubernetes.NewForConfig(config)
	ch := make(chan struct{})
	ns := "default"
	informers := informers.NewSharedInformerFactoryWithOptions(clientset, 10*time.Minute, informers.WithNamespace(ns))
	c := customcontroller.NewController(*clientset, informers.Apps().V1().Deployments())
	informers.Start(ch)
	c.Run(ch)
}
