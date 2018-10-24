package goldpinger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

// Status send the http latencies json encoded to the http request
func Status(w http.ResponseWriter, r *http.Request, ch ModelAgent) {
	json, err := json.Marshal(model(ch))
	if err != nil {
		Log("failed to marshal model to json: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func Gossiping(ch ModelAgent, r RandomNode) {
	t := time.NewTicker(2 * time.Second)
	for range t.C {
		trg, err := r(ch)
		if err != nil {
			Log("failed to gossip: %v", err)
			return
		}
		go gossip(ch, trg)
	}
}

func gossip(ch ModelAgent, t *Node) {
	resp, err := netClient.Get(fmt.Sprintf("http://%s/status.json", t.PodIP))
	if err != nil {
		Log("failed to fetch http: %s", err)
		return
	}
	defer resp.Body.Close()

	fetchedModel := &Model{}
	err = json.NewDecoder(resp.Body).Decode(fetchedModel)
	if err != nil {
		Log("failed to decode json model: %s", err)
		return
	}

	ch <- func(m *Model) {
		*m = Merge(*fetchedModel, *m)
	}
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
	sort.Sort(ByHostname(left))
	return left
}

type ByHostname []*Node

func (a ByHostname) Len() int           { return len(a) }
func (a ByHostname) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHostname) Less(i, j int) bool { return a[i].HostName < a[j].HostName }

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
