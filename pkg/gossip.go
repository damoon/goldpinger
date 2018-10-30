package goldpinger

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		*m = MergeModel(*fetchedModel, *m)
	}
}
