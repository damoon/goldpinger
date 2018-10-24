package goldpinger

import (
	"fmt"
	"sort"
	"time"

	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	k8sClient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func PodListSyncing(kubeconfig, namespace string, ch ModelAgent) {
	for {
		watch, err := newPodWatch(kubeconfig, namespace)
		if err != nil {
			Log("failed to watch pods: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		Log("created new watch for kubernetes pods\n")
		for event := range watch.ResultChan() {
			go updateTargets(ch, event)
		}
		Log("pods watch channel got closed\n")
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

func updateTargets(s ModelAgent, e watch.Event) {
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

		s <- func(m *Model) {
			m.Nodes = addNodeIfMissing(m.Nodes, &Node{
				HostIP:   pod.Status.HostIP,
				HostName: pod.Spec.NodeName,
				PodIP:    pod.Status.PodIP,
				PodName:  pod.Name,
			})
		}
	case watch.Error:
		fmt.Printf("%+v\n", e.Object)
	default:
		Log("unknown event: %+v\n", e)
	}
}

func addNodeIfMissing(nodes []*Node, node *Node) []*Node {
	for _, n := range nodes {
		if n.HostName == node.HostName {
			n.HostIP, n.PodName, n.PodIP = node.HostIP, node.PodName, node.PodIP
			return nodes
		}
	}

	list := append(nodes, node)
	sort.Sort(byHostname(list))
	return list
}

type byHostname []*Node

func (a byHostname) Len() int           { return len(a) }
func (a byHostname) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byHostname) Less(i, j int) bool { return a[i].HostName < a[j].HostName }
