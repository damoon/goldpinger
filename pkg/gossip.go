package goldpinger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mohae/deepcopy"
)

// PublishStatus sends the latencies as json encoded to via http.
func PublishStatus(w http.ResponseWriter, r *http.Request, ch ModelAgent) {

	c := make(chan Status)
	ch <- func(m *Model) {
		defer close(c)
		cp := deepcopy.Copy((*m).Status)
		c <- cp.(Status)
	}
	status := <-c

	json, err := json.Marshal(status)
	if err != nil {
		Log("failed to marshal model to json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500 - internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func Gossiping(ch ModelAgent) {
	t := time.NewTicker(2 * time.Second)
	for range t.C {
		trg, err := ch.randomNode()
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

	fetchedStatus := &Status{}
	err = json.NewDecoder(resp.Body).Decode(fetchedStatus)
	if err != nil {
		Log("failed to decode json model: %s", err)
		return
	}

	ch <- func(m *Model) {
		m.Status = MergeStatus(*fetchedStatus, m.Status)
	}
}
