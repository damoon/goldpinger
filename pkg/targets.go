package goldpinger

import (
	"fmt"
	"log"
	"sort"

	"k8s.io/api/core/v1"
	watch "k8s.io/apimachinery/pkg/watch"
)

func updateTargets(s chan<- func(p *Pinger), e watch.Event) {
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
		log.Printf("event %s for pod %s", e.Type, pod.Name)

		s <- func(p *Pinger) {
			p.model.Nodes = addNodeIfMissing(p.model.Nodes, &Node{
				HostIP:   pod.Status.HostIP,
				HostName: pod.Spec.NodeName,
				PodIP:    pod.Status.PodIP,
				PodName:  pod.Name,
			})
		}
	case watch.Error:
		fmt.Printf("%+v\n", e.Object)
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
