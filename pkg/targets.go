package goldpinger

import (
	"fmt"
	"math/rand"

	"k8s.io/api/core/v1"
	watch "k8s.io/apimachinery/pkg/watch"
)

func updateTargets(p *Pinger, e watch.Event) {
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
		t := &Target{
			HostIP:   pod.Status.HostIP,
			HostName: pod.Spec.NodeName,
			PodIP:    pod.Status.PodIP,
			PodName:  pod.Name,
		}
		p.targets[pod.Spec.NodeName] = t
		src, ok := p.model[pod.Spec.NodeName]
		if !ok {
			p.model[pod.Spec.NodeName] = &Source{Target: *t}
			return
		}
		src.HostIP, src.HostName, src.PodName, src.PodIP = t.HostIP, t.HostName, t.PodName, t.PodIP
	case watch.Error:
		fmt.Printf("%+v\n", e.Object)
	}
}

func randTarget(m map[string]*Target, r *rand.Rand) (*Target, error) {
	i := r.Intn(len(m))
	for k := range m {
		if i == 0 {
			return m[k], nil
		}
		i--
	}
	return nil, fmt.Errorf("can not select from empty target list")
}
