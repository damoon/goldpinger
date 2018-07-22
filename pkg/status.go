package goldpinger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mohae/deepcopy"
)

// Status send the http latencies json encoded to the http request
func Status(w http.ResponseWriter, r *http.Request, p *Pinger) {

	m := p.mock()

	json, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("failed to marshal model to json: %v", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, string(json))
}

func (p *Pinger) Model() *Model {
	r := make(chan *Model)
	p.synchronized <- func(p *Pinger) {
		c := deepcopy.Copy(p.model)
		r <- c.(*Model)
		close(r)
	}
	return <-r
}
