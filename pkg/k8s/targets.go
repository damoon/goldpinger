package k8s

import (
	"fmt"
	"time"

	"github.com/damoon/goldpinger/pkg"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	k8sClient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func PodListSyncing(kubeconfig, namespace string, ch goldpinger.ModelAgent) {
	for {
		watch, err := newPodWatch(kubeconfig, namespace)
		if err != nil {
			goldpinger.Log("failed to watch pods: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		goldpinger.Log("created new watch for kubernetes pods\n")
		for event := range watch.ResultChan() {
			go updateTargets(ch, event)
		}
		goldpinger.Log("pods watch channel got closed\n")
		time.Sleep(1 * time.Second)
	}
}

func newPodWatch(kubeconfig, namespace string) (watch.Interface, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load config for kubernetes client: %v", err)
	}
	client, err := k8sClient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	return client.CoreV1().Pods(namespace).Watch(meta_v1.ListOptions{})
}

func updateTargets(s goldpinger.ModelAgent, e watch.Event) {
	switch e.Type {
	case watch.Added:
		fallthrough
	case watch.Modified:
		fallthrough
	case watch.Deleted:
		pod, ok := e.Object.(*v1.Pod)
		if !ok {
			fmt.Printf("failed to cast %+v to a *v1.Pod\n", e.Object)
			return
		}

		s <- func(m *goldpinger.Model) {
			node := &goldpinger.Node{
				HostIP:   pod.Status.HostIP,
				HostName: pod.Spec.NodeName,
				PodIP:    pod.Status.PodIP,
				PodName:  pod.Name,
			}
			m.Nodes = goldpinger.Add(m.Nodes, node)
		}
	case watch.Error:
		fmt.Printf("%+v\n", e.Object)
	default:
		goldpinger.Log("unknown event: %+v\n", e)
	}
}
