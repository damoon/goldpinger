package goldpinger

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// OK confirms a http connection was created
func OK(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Printf("failed to send response: %v", err)
	}
}

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

func fetchHTTP(s chan<- func(p *Pinger), targets map[string]*Target, r *rand.Rand) {
	t, err := randTarget(targets, r)
	if err != nil {
		log.Printf("failed to fetch http: %v", err)
		return
	}

	d, err := measureHTTP(fmt.Sprintf("http://%s/ok", t.PodIP))
	s <- func(p *Pinger) {
		_, ok := p.model[p.nodeName]
		if !ok {
			log.Printf("failed to fetch http: source is not set up yet", err)
			return
		}
		p.model[p.nodeName].Measurements[t.HostName] = &Measurement{
			Delay:     d,
			Error:     err.Error(),
			Timestamp: time.Now().UnixNano(),
		}
	}
}

func measureHTTP(url string) (int64, error) {
	before := time.Now().UnixNano()
	response, err := netClient.Get(url)
	d := time.Now().UnixNano() - before
	if err != nil {
		return d, fmt.Errorf("failed to fetch http: %s", err)
	}
	defer response.Body.Close()
	return d, nil
}
