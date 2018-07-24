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

func fetchHTTP(s chan<- func(p *Pinger), targets []*Node, r *rand.Rand) {
	t, err := randTarget(targets, r)
	if err != nil {
		log.Printf("failed to fetch http: %v", err)
		return
	}

	d, err := measureHTTP(fmt.Sprintf("http://%s/ok", t.PodIP))
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	s <- func(p *Pinger) {
		addMessurement(p.model.Measurements, p.nodeName, t.HostName, &Measurement{
			Delay:     d,
			Error:     errMsg,
			Timestamp: time.Now().UnixNano(),
		})
	}
}

func addMessurement(table map[string]map[string]*Measurement, source, target string, m *Measurement) {
	_, ok := table[source]
	if !ok {
		table[source] = map[string]*Measurement{}
	}
	table[source][target] = m
}

func measureHTTP(url string) (int64, error) {
	before := time.Now().UnixNano()
	resp, err := netClient.Get(url)
	d := time.Now().UnixNano() - before
	if err != nil {
		return d, fmt.Errorf("failed to fetch http: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return d, fmt.Errorf("failed to fetch http: returned status code %d from %s", resp.StatusCode, url)
	}
	return d, nil
}

func randTarget(m []*Node, r *rand.Rand) (*Node, error) {
	l := len(m)
	if l == 0 {
		return nil, fmt.Errorf("can not select from empty target list")
	}
	return m[r.Intn(l)], nil
}
