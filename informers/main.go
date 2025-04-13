package main

import (
	"flag"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/lhome/kcharan/.kube/config", "kubeconfig location")
	config, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	clientset, _ := kubernetes.NewForConfig(config)
	ns := "default"
	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		time.Minute,
		informers.WithNamespace(ns),
		informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
			opts.FieldSelector = fields.Everything().String()
		}),
	)
	podInformer := factory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Println("🆕 Pod Added:", pod.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod := newObj.(*corev1.Pod)
			fmt.Println("✏️ Pod Updated:", pod.Name)
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Println("❌ Pod Deleted:", pod.Name)
		},
	})
	stopCh := make(chan struct{})
	defer close(stopCh)
	fmt.Println("🔄 Starting Pod Informer...")
	factory.Start(stopCh)

	// Step 6: Wait for cache sync
	if ok := cache.WaitForCacheSync(stopCh, podInformer.HasSynced); !ok {
		fmt.Println("❌ Failed to sync")
		return
	}

	fmt.Println("✅ Cache synced. Listening for Pod events...")

	select {}
}
