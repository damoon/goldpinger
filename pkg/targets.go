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

		s <- func(p *Pinger) {
			p.targets[pod.Spec.NodeName] = &Node{
				HostIP:   pod.Status.HostIP,
				HostName: pod.Spec.NodeName,
				PodIP:    pod.Status.PodIP,
				PodName:  pod.Name,
			}

			for _, s := range *p.model {
				if s.HostName == pod.Spec.NodeName {
					s.HostIP, s.PodName, s.PodIP = pod.Status.HostIP, pod.Name, pod.Status.PodIP
					return
				}
			}
			*p.model = append((*p.model), &Source{
				Node:         *p.targets[pod.Spec.NodeName],
				Measurements: []*Measurement{},
			})
		}
	case watch.Error:
		fmt.Printf("%+v\n", e.Object)
	}
}
