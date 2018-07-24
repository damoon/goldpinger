package goldpinger

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
)

func gossip(s chan<- func(p *Pinger), targets []*Node, r *rand.Rand) {
	t, err := randTarget(targets, r)
	if err != nil {
		log.Printf("failed to gossip: %s", err)
		return
	}

	resp, err := netClient.Get(fmt.Sprintf("http://%s/status.json", t.PodIP))
	if err != nil {
		log.Printf("failed to fetch http: %s", err)
		return
	}
	defer resp.Body.Close()

	fetchedModel := &Model{}
	err = json.NewDecoder(resp.Body).Decode(fetchedModel)

	s <- func(p *Pinger) {
		p.model = merge(*fetchedModel, p.model)
	}
}

func merge(right, left Model) Model {
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
