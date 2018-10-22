package goldpinger

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	k8sClient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var Log = log.Printf

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
	rand          *rand.Rand
	nodeName      string
	synchronized  chan func(p *Pinger)
	kubeConfig    string
	kubeNamespace string
	fetchHTTP     *time.Ticker
	gossip        *time.Ticker
	model         Model
}

func StartNew(nodeName string, kubeConfig, namespace string, r *rand.Rand) *Pinger {
	c := make(chan func(p *Pinger))
	p := &Pinger{
		rand:          r,
		nodeName:      nodeName,
		synchronized:  c,
		kubeConfig:    kubeConfig,
		kubeNamespace: namespace,
		fetchHTTP:     time.NewTicker(1 * time.Second),
		gossip:        time.NewTicker(2 * time.Second),
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
			case <-p.gossip.C:
				go gossip(p.synchronized, p.model.Nodes, p.rand)
			case <-p.gossip.C:
				go fetchHTTP(p.synchronized, p.model.Nodes, p.rand)
			}
		}
	}()
	go p.updatePods()
}

func (p *Pinger) updatePods() {
	for {
		watch, err := p.podWatch()
		if err != nil {
			Log("failed to watch pods: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		Log("created new watch for kubernetes pods\n")
		for event := range watch.ResultChan() {
			go updateTargets(p.synchronized, event)
		}
		Log("pods watch channel got closed\n")
		time.Sleep(1 * time.Second)
	}
}

func (p *Pinger) podWatch() (watch.Interface, error) {

	config, err := clientcmd.BuildConfigFromFlags("", p.kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load config for kubernetes client: %v", err)
	}
	client, err := k8sClient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	return client.CoreV1().Pods(p.kubeNamespace).Watch(meta_v1.ListOptions{})
}
