package goldpinger

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/mohae/deepcopy"
	"gopkg.in/d4l3k/messagediff.v1"
	v1 "k8s.io/api/core/v1"
	watch "k8s.io/apimachinery/pkg/watch"
)

type model *[]host

type host struct {
	Source string `json:"source"`
	Pings  []ping `json:"pings"`
}

type ping struct {
	Target    string `json:"target"`
	Delay     int    `json:"delay"`
	Timestamp int64  `json:"timestamp"`
}

type Pinger struct {
	nextPing     *time.Ticker
	podsWatch    <-chan watch.Event
	nodesWatch   <-chan watch.Event
	pods         []*v1.Pod
	nodes        []*v1.Node
	model        *[]host
	synchronized chan func()
}

func NewPinger(p <-chan watch.Event, n <-chan watch.Event) *Pinger {
	c := make(chan func())
	return &Pinger{
		nextPing:     time.NewTicker(time.Second),
		podsWatch:    p,
		nodesWatch:   n,
		synchronized: c,
	}
}

func (p *Pinger) Mock() *[]host {
	return &[]host{
		host{
			Source: "host1",
			Pings: []ping{
				ping{Target: "host1", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host2", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host3", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host4", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
			},
		},
		host{
			Source: "host2",
			Pings: []ping{
				ping{Target: "host1", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host2", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host3", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host4", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
			},
		},
		host{
			Source: "host3",
			Pings: []ping{
				ping{Target: "host1", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host2", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host3", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host4", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
			},
		},
		host{
			Source: "host4",
			Pings: []ping{
				ping{Target: "host1", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host2", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host3", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host4", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
			},
		},
	}
}

func (p *Pinger) Start() {
	go func() {
		for {
			select {
			case <-p.nextPing.C:
				//t := selectRandomTarget(p.pods, p.nodes)
				//ping(t)
			case e := <-p.podsWatch:
				updatePods(p, e)
			case e := <-p.nodesWatch:
				updateNodes(p, e)
			case f := <-p.synchronized:
				f()
			}
		}
	}()
}

func (p *Pinger) Model() *[]host {
	r := make(chan *[]host)
	p.synchronized <- func() {
		c := deepcopy.Copy(p.model)
		r <- c.(*[]host)
		close(r)
	}
	return <-r
}

func updatePods(p *Pinger, e watch.Event) {
	switch e.Type {
	case watch.Added:
		pod, ok := e.Object.(*v1.Pod)
		if !ok {
			fmt.Printf("failed to cast %+v to a *v1.Pod\n", e.Object)
			return
		}
		addPod(p, pod)
	case watch.Modified:
		pod, ok := e.Object.(*v1.Pod)
		if !ok {
			fmt.Printf("failed to cast %+v to a *v1.Pod\n", e.Object)
			return
		}
		updatePod(p, pod)
	case watch.Deleted:
		pod, ok := e.Object.(*v1.Pod)
		if !ok {
			fmt.Printf("failed to cast %+v to a *v1.Pod\n", e.Object)
			return
		}
		removePod(p, pod)
	case watch.Error:
		fmt.Printf("%+v\n", e.Object)
	}
}

func addPod(p *Pinger, pod *v1.Pod) {
	log.Printf("add Pod %v", pod.Name)
	p.pods = append(p.pods, pod)
}

func updatePod(p *Pinger, pod *v1.Pod) {
	log.Printf("update Pod %v", pod.Name)
	for i, c := range p.pods {
		if pod.Name == c.Name {
			diff, _ := messagediff.PrettyDiff(pod, p.pods[i])
			log.Printf("found changes %s", diff)
			p.pods[i] = pod
			return
		}
	}
}

func removePod(p *Pinger, pod *v1.Pod) {
	log.Printf("remove Pod %v", pod.Name)
	for i, c := range p.pods {
		if pod.Name == c.Name {
			p.pods = append(p.pods[:i], p.pods[i+1:]...)
			return
		}
	}
}

func updateNodes(p *Pinger, e watch.Event) {
	switch e.Type {
	case watch.Added:
		node, ok := e.Object.(*v1.Node)
		if !ok {
			fmt.Printf("failed to cast %+v to a *v1.Node\n", e.Object)
			return
		}
		addNode(p, node)
	case watch.Modified:
		node, ok := e.Object.(*v1.Node)
		if !ok {
			fmt.Printf("failed to cast %+v to a *v1.Node\n", e.Object)
			return
		}
		updateNode(p, node)
	case watch.Deleted:
		node, ok := e.Object.(*v1.Node)
		if !ok {
			fmt.Printf("failed to cast %+v to a *v1.Node\n", e.Object)
			return
		}
		removeNode(p, node)
	case watch.Error:
		fmt.Printf("%+v\n", e.Object)
	}
}

func addNode(p *Pinger, n *v1.Node) {
	log.Printf("add Node %v", n.Name)
	p.nodes = append(p.nodes, n)
}

func updateNode(p *Pinger, n *v1.Node) {
	log.Printf("update Node %v", n.Name)
	for i, c := range p.nodes {
		if n.Name == c.Name {
			diff, _ := messagediff.PrettyDiff(n, p.nodes[i])
			log.Printf("found changes %s", diff)
			p.nodes[i] = n
			return
		}
	}
}

func removeNode(p *Pinger, n *v1.Node) {
	log.Printf("remove Node %v", n.Name)
	for i, c := range p.nodes {
		if n.Name == c.Name {
			p.nodes = append(p.nodes[:i], p.nodes[i+1:]...)
			return
		}
	}
}
