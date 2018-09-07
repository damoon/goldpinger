package goldpinger

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	watch "k8s.io/apimachinery/pkg/watch"
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
	podsWatch    <-chan watch.Event
	fetchHTTP    *time.Ticker
	gossip       *time.Ticker
	model        Model
}

func NewPinger(nodeName string, p <-chan watch.Event, r *rand.Rand) *Pinger {
	c := make(chan func(p *Pinger))
	return &Pinger{
		rand:         r,
		nodeName:     nodeName,
		synchronized: c,
		podsWatch:    p,
		fetchHTTP:    time.NewTicker(1 * time.Second),
		gossip:       time.NewTicker(2 * time.Second),
		model: Model{
			Nodes:        []*Node{},
			Measurements: map[string]map[string]*Measurement{},
		},
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
				go fetchHTTP(p.synchronized, p.model.Nodes, p.rand)
			case <-p.gossip.C:
				go gossip(p.synchronized, p.model.Nodes, p.rand)
			case event := <-p.podsWatch:
				go updateTargets(p.synchronized, event)
			}

			log.Printf("%d running go routines\n", runtime.NumGoroutine())
		}
	}()
}
