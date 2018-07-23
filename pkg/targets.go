package goldpinger

import (
	"fmt"
	"log"

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
		n := &Node{
			HostIP:   pod.Status.HostIP,
			HostName: pod.Spec.NodeName,
			PodIP:    pod.Status.PodIP,
			PodName:  pod.Name,
		}
		s <- func(p *Pinger) {
			p.targets[pod.Spec.NodeName] = n
			src, ok := (*p.model)[pod.Spec.NodeName]
			if !ok {
				(*p.model)[pod.Spec.NodeName] = &Source{
					Node:         *n,
					Measurements: map[string]*Measurement{},
				}
				return
			}
			src.HostIP, src.HostName, src.PodName, src.PodIP = n.HostIP, n.HostName, n.PodName, n.PodIP
		}
	case watch.Error:
		fmt.Printf("%+v\n", e.Object)
	}
}
