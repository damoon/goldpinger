package goldpinger

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mohae/deepcopy"
)

// Status send the http latencies json encoded to the http request
func (p *Pinger) Status(w http.ResponseWriter, r *http.Request) {
	status(w, r, p.synchronized)
}

func status(w http.ResponseWriter, r *http.Request, sync chan<- func(p *Pinger)) {
	json, err := json.Marshal(model(sync))
	if err != nil {
		Log("failed to marshal model to json: %v", err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, string(json))
}

func model(sync chan<- func(p *Pinger)) Model {
	r := make(chan Model)
	sync <- func(p *Pinger) {
		c := deepcopy.Copy(p.model)
		r <- c.(Model)
		close(r)
	}
	return <-r
}
