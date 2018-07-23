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

func fetchHTTP(s chan<- func(p *Pinger), targets map[string]*Node, r *rand.Rand) {
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

		for _, source := range *p.model {
			if source.HostName == p.nodeName {
				for _, messurement := range source.Measurements {
					if messurement.Target == t.HostName {
						messurement.Delay, messurement.Error, messurement.Timestamp = d, errMsg, time.Now().UnixNano()
						return
					}
				}
				source.Measurements = append(source.Measurements, &Measurement{
					Target:    t.HostName,
					Delay:     d,
					Error:     errMsg,
					Timestamp: time.Now().UnixNano(),
				})
				return
			}
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

func randTarget(m map[string]*Node, r *rand.Rand) (*Node, error) {
	i := r.Intn(len(m))
	for k := range m {
		if i == 0 {
			return m[k], nil
		}
		i--
	}
	return nil, fmt.Errorf("can not select from empty target list")
}
