package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func createService(clientset *kubernetes.Clientset, namespace string) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-service",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
			Selector: map[string]string{
				"app": "myapp",
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	svc, err := clientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {

	}
	fmt.Printf("service %s created\n", svc.Name)
}
func main() {
	ns := "default"
	kubeConfig := flag.String("kubeconfig", "/lhome/kcharan/.kube/config", "kubeconfig location")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {

	}
	// clientset is of type *kubernetes.Clientset,
	// which provides typed clients for each API group/version.
	//  For example:
	// 	clientset.CoreV1()       // For core resources like Pods, Services, ConfigMaps
	// clientset.AppsV1()       // For Deployments, StatefulSets, DaemonSets
	// clientset.BatchV1()      // For Jobs, CronJobs
	// clientset.NetworkingV1() // For Ingress, NetworkPolicy

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {

	}
	ctx := context.Background()
	// for {
	pods, _ := clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
	for _, pod := range pods.Items {
		for k, v := range pod.Labels {
			fmt.Println(k, v)
		}
	}
	createService(clientset, ns)
	time.Sleep(10 * time.Second)
	// }
	// fmt.Printf("%+v", clientset)
}

// create a deployment using cmd
// kubectl create deployment nginx -n default --image nginx
