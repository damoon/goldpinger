package goldpinger

import (
	"fmt"
	"log"
	"math/rand"
	"sort"

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

func Add(nodes []*Node, node *Node) []*Node {
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

func Merge(right, left Model) Model {
	return Model{
		Nodes:        mergeNodes(right.Nodes, left.Nodes),
		Measurements: mergeMeasurementRows(right.Measurements, left.Measurements),
	}
}

func mergeNodes(right, left []*Node) []*Node {
	for _, r := range right {
		if !nodeExist(r, left) {
			left = append(left, r)
		}
	}
	sort.Sort(byHostname(left))
	return left
}

type byHostname []*Node

func (a byHostname) Len() int           { return len(a) }
func (a byHostname) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byHostname) Less(i, j int) bool { return a[i].HostName < a[j].HostName }

func nodeExist(node *Node, nodes []*Node) bool {
	for _, n := range nodes {
		if n.HostName == node.HostName {
			return true
		}
	}
	return false
}

func mergeMeasurementRows(right, left map[string]map[string]*Measurement) map[string]map[string]*Measurement {
	for k, v := range right {
		l, ok := left[k]
		if ok {
			v = mergeMeasurements(v, l)
		}
		left[k] = v
	}
	return left
}

func mergeMeasurements(right, left map[string]*Measurement) map[string]*Measurement {
	for k, v := range right {
		l, ok := left[k]
		if ok {
			v = newestMeasurements(v, l)
		}
		left[k] = v
	}
	return left
}

func newestMeasurements(right, left *Measurement) *Measurement {
	if right.Timestamp > left.Timestamp {
		return right
	}
	return left
}
