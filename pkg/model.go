package goldpinger

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/mohae/deepcopy"
)

var Log = log.Printf

type Model struct {
	Nodes        []*Node                            `json:"nodes"`
	Measurements map[string]map[string]*Measurement `json:"measurements"`
}

type ModelAgent chan<- func(m *Model)

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

func StartNewModel() ModelAgent {
	c := make(chan func(m *Model))
	m := &Model{
		Nodes:        []*Node{},
		Measurements: map[string]map[string]*Measurement{},
	}
	go func() {
		for f := range c {
			f(m)
		}
	}()
	return c
}

func model(ch ModelAgent) Model {
	r := make(chan Model)
	ch <- func(m *Model) {
		c := deepcopy.Copy(*m)
		r <- c.(Model)
		close(r)
	}
	return <-r
}

type RandomNode func(ModelAgent) (*Node, error)

func NewRandomNode(r *rand.Rand) RandomNode {
	return func(ch ModelAgent) (*Node, error) {
		m := model(ch)
		l := len(m.Nodes)
		if l == 0 {
			return nil, fmt.Errorf("can not select from empty target list")
		}
		return m.Nodes[r.Intn(l)], nil
	}
}
