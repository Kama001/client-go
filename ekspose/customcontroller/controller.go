package customcontroller

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type controller struct {
	clientset      kubernetes.Clientset
	depCacheSynced cache.InformerSynced
}

func NewController(clientset kubernetes.Clientset, depInformer appsinformers.DeploymentInformer) *controller {
	c := &controller{
		clientset:      clientset,
		depCacheSynced: depInformer.Informer().HasSynced,
	}
	depInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerDetailedFuncs{
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
		},
	)
	return c
}

func (c *controller) Run(ch chan struct{}) {
	fmt.Println("starting controller")
	if !cache.WaitForCacheSync(ch, c.depCacheSynced) {
		fmt.Print("waiting for cache to be synced\n")
	}

	// go wait.Until(c.worker, 1*time.Second, ch)

	<-ch
}

func getDetails(obj interface{}) (string, string) {
	deploy, ok := obj.(*appsv1.Deployment)
	if !ok {
		fmt.Println("Received unexpected type")
		return "", ""
	}
	return deploy.Namespace, deploy.Name
}

func (c *controller) handleAdd(obj interface{}, isInIntialList bool) {
	fmt.Println("add was called")
	ns, name := getDetails(obj)
	if ns != "" && name != "" {
		fmt.Printf("🆕 Deployment Added: %s/%s\n", ns, name)
	}
}

func (c *controller) handleDel(obj interface{}) {
	fmt.Println("del was called")
	ns, name := getDetails(obj)
	if ns != "" && name != "" {
		fmt.Printf("❌ Deployment Deleted:%s/%s\n", ns, name)
	}
}
