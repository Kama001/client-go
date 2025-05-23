package customctrlwithqueue

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type controller struct {
	clientset      kubernetes.Clientset
	depCacheSynced cache.InformerSynced // InformerSynced = func() bool
	// queue          workqueue.TypedRateLimitingInterface[string]
	queue     workqueue.RateLimitingInterface
	depLister appslisters.DeploymentLister
}

func NewController(clientset kubernetes.Clientset, depInformer appsinformers.DeploymentInformer) *controller {
	c := &controller{
		clientset: clientset,
		// redfer /lhome/kcharan/go_bacics/interface for examples
		depCacheSynced: depInformer.Informer().HasSynced, // HasSynced means func() bool, HasSynced() means bool
		// queue:          workqueue.NewTypedRateLimitingQueue[string](workqueue.NewTypedItemExponentialFailureRateLimiter[string](1*time.Second, 30*time.Second)),
		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(),
			"ekspose"),
		depLister: depInformer.Lister(),
	}
	// depInformer.Informer().SharedInformer.AddEventhandler
	// AddEventHandler(handler ResourceEventhandler)
	// type ResourceEventHandler interface {
	// 	OnAdd(obj interface{}, isInInitialList bool)
	// 	OnUpdate(oldObj, newObj interface{})
	// 	OnDelete(obj interface{})
	// }

	depInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerDetailedFuncs{
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
		},
	)
	// func (r ResourceEventHandlerDetailedFuncs) OnAdd(obj interface{}) {
	// 	if r.AddFunc != nil {
	// 		r.AddFunc(obj) // ← 💥 YOUR function gets called here
	// 	}
	// }
	return c
}

func (c *controller) Run(ch chan struct{}) {
	fmt.Println("⏳ starting controller")
	// WaitForCacheSync waits for caches to populate.
	// It returns true if it was successful,
	// false if the controller should shutdown callers should prefer WaitForNamedCacheSync()
	// c.depCacheSynced = func() bool
	if !cache.WaitForCacheSync(ch, c.depCacheSynced) {
		fmt.Print("waiting for cache to be synced\n")
	}

	go wait.Until(c.worker, 1*time.Second, ch)

	<-ch
}

func (c *controller) worker() {
	for c.processItem() {

	}
}

func (c *controller) processItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Forget(item)
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Printf("getting key from cahce %s\n", err.Error())
	}
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Printf("splitting key into namespace and name %s\n", err.Error())
		return false
	}
	ctx := context.Background()
	_, err = c.clientset.AppsV1().Deployments(ns).Get(ctx, name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {

		err := c.clientset.CoreV1().Services(ns).Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("deleting service %s\n", name)
			return false
		}
		fmt.Printf("❌ deployment %s deleted\n", name)
		return true
	}
	err = c.syncDeployment(ns, name)
	if err != nil {
		fmt.Printf("syncing deployment %s\n", err.Error())
		return false
	}
	return true
}

func (c *controller) syncDeployment(ns, name string) error {
	ctx := context.Background()
	dep, err := c.depLister.Deployments(ns).Get(name)
	if err != nil {
		fmt.Printf("getting deployment from lister %s\n", err.Error())
		return err
	}
	port := c.getContainerPorts(dep)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: dep.Name,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: dep.Name,
					Port: port,
				},
			},
			Selector: dep.Labels,
			Type:     corev1.ServiceTypeClusterIP,
		},
	}
	svc, err := c.clientset.CoreV1().Services(ns).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {

	}
	fmt.Printf("🆕 service %s created\n", svc.Name)
	return nil
}
func (c *controller) getContainerPorts(dep *appsv1.Deployment) int32 {
	for _, container := range dep.Spec.Template.Spec.Containers {
		for _, port := range container.Ports {
			return port.ContainerPort
		}
	}
	return 80
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
	// fmt.Println("add was called")
	// key, _ := cache.MetaNamespaceKeyFunc(obj)
	// fmt.Println(cache.SplitMetaNamespaceKey(key))
	// ns, name := getDetails(obj)
	// if ns != "" && name != "" {
	// 	fmt.Printf("🆕 Deployment Added: %s/%s\n", ns, name)
	// }
	c.queue.Add(obj)
}

func (c *controller) handleDel(obj interface{}) {
	// fmt.Println("del was called")
	// ns, name := getDetails(obj)
	// if ns != "" && name != "" {
	// 	fmt.Printf("❌ Deployment Deleted:%s/%s\n", ns, name)
	// }
	c.queue.Add(obj)
}
