package goldpinger

import (
	"math/rand"
	"time"

	watch "k8s.io/apimachinery/pkg/watch"
)

type Model map[string]*Source

type Source struct {
	Node
	Measurements map[string]*Measurement `json:"measurements"`
}

type Measurement struct {
	Timestamp int64  `json:"timestamp"`
	Delay     int64  `json:"delay"`
	Error     string `json:"error"`
}

type Node struct {
	HostName string `json:"hostName"`
	HostIP   string `json:"hostIP"`
	PodName  string `json:"podName"`
	PodIP    string `json:"podIP"`
}

type Pinger struct {
	rand         *rand.Rand
	nodeName     string
	synchronized chan func(p *Pinger)
	podsWatch    <-chan watch.Event
	targets      map[string]*Node
	fetchHTTP    *time.Ticker
	gossip       *time.Ticker
	model        *Model
}

func NewPinger(nodeName string, p <-chan watch.Event, r *rand.Rand) *Pinger {
	c := make(chan func(p *Pinger))
	return &Pinger{
		rand:         r,
		nodeName:     nodeName,
		synchronized: c,
		podsWatch:    p,
		targets:      map[string]*Node{},
		fetchHTTP:    time.NewTicker(2 * time.Second),
		gossip:       time.NewTicker(4 * time.Second),
		model:        &Model{},
	}
}

func (p *Pinger) Start() {
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
			select {
			case <-p.fetchHTTP.C:
				go fetchHTTP(p.synchronized, p.targets, p.rand)
			case <-p.gossip.C:
				//				go gossip(p.targets, p.rand)
			case event := <-p.podsWatch:
				go updateTargets(p.synchronized, event)
			}
		}
	}()
}
