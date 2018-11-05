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
	Status Status
	random *rand.Rand
}

type Status struct {
	Participants *[]Node                       `json:"nodes"`
	Worldview    map[string]map[string]History `json:"measurements"`
}

type ModelAccess chan<- func(m Model)

type Node struct {
	HostName string `json:"hostName"`
	HostIP   string `json:"hostIP"`
	PodName  string `json:"podName"`
	PodIP    string `json:"podIP"`
}

type History []Measurement

type Measurement struct {
	Timestamp int64  `json:"timestamp"`
	Delay     int64  `json:"delay"`
	Error     string `json:"error"`
}

func StartNewModel(r *rand.Rand) ModelAccess {
	c := make(chan func(m Model))
	m := Model{
		Status: Status{
			Participants: &[]Node{},
			Worldview:    map[string]map[string]History{},
		},
		random: r,
	}
	go func() {
		for f := range c {
			f(m)
		}
	}()
	return c
}

func (ch ModelAccess) randomNode() (Node, error) {
	type response struct {
		node Node
		err  error
	}
	c := make(chan response)
	ch <- func(m Model) {
		defer close(c)
		l := len(*m.Status.Participants)
		if l == 0 {
			c <- response{Node{}, fmt.Errorf("can not select from empty target list")}
			return
		}
		i := m.random.Intn(l)
		n := deepcopy.Copy((*m.Status.Participants)[i])
		c <- response{n.(Node), nil}
	}
	n := <-c
	return n.node, n.err
}

func (ch ModelAccess) Add(node Node) {
	ch <- func(m Model) {
		for _, n := range *m.Status.Participants {
			if n.HostName == node.HostName {
				n.HostIP, n.PodName, n.PodIP = node.HostIP, node.PodName, node.PodIP
				return
			}
		}
		*m.Status.Participants = append(*m.Status.Participants, node)
		sort.Sort(byHostname(*m.Status.Participants))
	}
}

func MergeStatus(right, left Status) Status {
	return Status{
		Participants: mergeParticipants(right.Participants, left.Participants),
		Worldview:    mergeWorldview(right.Worldview, left.Worldview),
	}
}

func mergeParticipants(pp ...*[]Node) *[]Node {
	resp := []Node{}
	for _, p := range pp {
		for _, n := range *p {
			if participantMissing(resp, n) {
				resp = append(resp, n)
			}
		}
	}
	sort.Sort(byHostname(resp))
	return &resp
}

type byHostname []Node

func (a byHostname) Len() int           { return len(a) }
func (a byHostname) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byHostname) Less(i, j int) bool { return a[i].HostName < a[j].HostName }

func participantMissing(nodes []Node, node Node) bool {
	for _, n := range nodes {
		if n.HostName == node.HostName {
			return false
		}
	}
	return true
}

func mergeWorldview(right, left map[string]map[string]History) map[string]map[string]History {
	for k, v := range right {
		l, ok := left[k]
		if ok {
			v = mergeParticipantview(v, l)
		}
		left[k] = v
	}
	return left
}

func mergeParticipantview(right, left map[string]History) map[string]History {
	for k, r := range right {
		l, ok := left[k]
		if ok {
			r = mergeHistories(r, l)
		}
		left[k] = r
	}
	return left
}

func mergeHistories(right, left History) History {

	i, j := 0, 0
	slice := History{}

	for len(slice) < maxHistoryLength && (i < len(right) || j < len(left)) {

		if i == len(right) {
			slice = append(slice, left[j])
			j++
			continue
		}
		if j == len(left) {
			slice = append(slice, right[i])
			i++
			continue
		}

		if right[i].Timestamp == left[j].Timestamp {
			slice = append(slice, right[i])
			i++
			j++
			continue
		}
		if right[i].Timestamp > left[j].Timestamp {
			slice = append(slice, right[i])
			i++
			continue
		}
		if right[i].Timestamp < left[j].Timestamp {
			slice = append(slice, left[j])
			j++
			continue
		}
	}

	return slice

}

type byTimestamp History

func (a byTimestamp) Len() int           { return len(a) }
func (a byTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTimestamp) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }
