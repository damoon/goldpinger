package goldpinger

import (
	"fmt"
	"log"
	"math/rand"
	"sort"

	"github.com/mohae/deepcopy"
)

var Log = log.Printf

const maxHistoryLength = 10

type Model struct {
	Participants []*Node                       `json:"nodes"`
	Worldview    map[string]map[string]History `json:"measurements"`
}

type ModelAgent chan<- func(m *Model)

type Node struct {
	HostName string `json:"hostName"`
	HostIP   string `json:"hostIP"`
	PodName  string `json:"podName"`
	PodIP    string `json:"podIP"`
}

type History []*Measurement

type Measurement struct {
	Timestamp int64  `json:"timestamp"`
	Delay     int64  `json:"delay"`
	Error     string `json:"error"`
}

func StartNewModel() ModelAgent {
	c := make(chan func(m *Model))
	m := &Model{
		Participants: []*Node{},
		Worldview:    map[string]map[string]History{},
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
		l := len(m.Participants)
		if l == 0 {
			return nil, fmt.Errorf("can not select from empty target list")
		}
		return m.Participants[r.Intn(l)], nil
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

func MergeModel(right, left Model) Model {
	return Model{
		Participants: mergeParticipants(right.Participants, left.Participants),
		Worldview:    mergeWorldview(right.Worldview, left.Worldview),
	}
}

func mergeParticipants(right, left []*Node) []*Node {
	for _, r := range right {
		if !participantExist(r, left) {
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

func participantExist(node *Node, nodes []*Node) bool {
	for _, n := range nodes {
		if n.HostName == node.HostName {
			return true
		}
	}
	return false
}

func mergeWorldview(right, left map[string]map[string]History) map[string]map[string]History {
	for k, v := range right {
		l, ok := left[k]
		if ok {
			v = mergeParticipantsView(v, l)
		}
		left[k] = v
	}
	return left
}

func mergeParticipantsView(right, left map[string]History) map[string]History {
	for k, v := range right {
		l, ok := left[k]
		if ok {
			v = mergeHistories(v, l)
		}
		left[k] = v
	}
	return left
}

func mergeHistories(right, left History) History {
	h := append(right, left...)
	sort.Sort(byTimestamp(h))
	size := maxHistoryLength
	if size > len(h) {
		size = len(h)
	}
	return h[:size]
}

type byTimestamp History

func (a byTimestamp) Len() int           { return len(a) }
func (a byTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTimestamp) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }
