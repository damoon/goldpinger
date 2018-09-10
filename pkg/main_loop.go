package goldpinger

import (
	"log"
	"math/rand"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
)

type Model struct {
	Nodes        []*Node                            `json:"nodes"`
	Measurements map[string]map[string]*Measurement `json:"measurements"`
}

type Node struct {
	HostName string `json:"hostName"`
	HostIP   string `json:"hostIP"`
	PodName  string `json:"podName"`
	PodIP    string `json:"podIP"`
}

type Measurement struct {
	Timestamp int64  `json:"timestamp"`
	Delay     int64  `json:"delay"`
	Error     string `json:"error"`
}

type Pinger struct {
	rand         *rand.Rand
	nodeName     string
	synchronized chan func(p *Pinger)
	pods         v1.PodInterface
	fetchHTTP    *time.Ticker
	gossip       *time.Ticker
	model        Model
}

func StartNew(nodeName string, pods v1.PodInterface, r *rand.Rand) *Pinger {
	c := make(chan func(p *Pinger))
	p := &Pinger{
		rand:         r,
		nodeName:     nodeName,
		synchronized: c,
		pods:         pods,
		fetchHTTP:    time.NewTicker(1 * time.Second),
		gossip:       time.NewTicker(2 * time.Second),
		model: Model{
			Nodes:        []*Node{},
			Measurements: map[string]map[string]*Measurement{},
		},
	}
	p.start()
	return p
}

func (p *Pinger) start() {
	go func() {
		for {
			select {
			case f := <-p.synchronized:
				f(p)
			}
		}
	}()
	go func() {
		for {
			watch, err := p.pods.Watch(meta_v1.ListOptions{})
			if err != nil {
				log.Fatalf("failed to watch pods: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			for {
				select {
				case <-p.fetchHTTP.C:
					go fetchHTTP(p.synchronized, p.model.Nodes, p.rand)
				case <-p.gossip.C:
					go gossip(p.synchronized, p.model.Nodes, p.rand)
				case event, ok := <-watch.ResultChan():
					if !ok {
						time.Sleep(1 * time.Second)
						break
					}
					go updateTargets(p.synchronized, event)
				}
			}
		}
	}()
}
