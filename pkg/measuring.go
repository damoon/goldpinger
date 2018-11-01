package goldpinger

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// OK confirms a http connection was created
func OK(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	if err != nil {
		Log("failed to send response: %v", err)
	}
}

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

func Measuring(ch ModelAgent, r RandomNode, hostname string) {
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		trg, err := r(ch)
		if err != nil {
			Log("failed to ping: %v", err)
			return
		}
		url := fmt.Sprintf("http://%s/ok", trg.PodIP)
		go fetchHTTP(ch, trg.HostName, hostname, url)
	}
}

func fetchHTTP(ch ModelAgent, target, source, url string) {
	d, err := measureHTTP(url)
	t := time.Now().UnixNano()
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	ch <- func(m *Model) {
		_, ok := m.Worldview[source]
		if !ok {
			m.Worldview[source] = map[string]History{}
		}
		h := History{Measurement{
			Delay:     d,
			Error:     errMsg,
			Timestamp: t,
		}}
		m.Worldview[source][target] = mergeHistories(m.Worldview[source][target], h)
	}
}

func addMeasurement(table map[string]map[string]History, source, target string, m Measurement) {
	_, ok := table[source]
	if !ok {
		table[source] = map[string]History{}
	}
	table[source][target] = mergeHistories(table[source][target], History{m})
}

func measureHTTP(url string) (int64, error) {
	before := time.Now().UnixNano()
	resp, err := netClient.Get(url)
	d := time.Now().UnixNano() - before
	if err != nil {
		return d, fmt.Errorf("failed to fetch http: %s", err)
	}
	defer resp.Body.Close()
	// https://husobee.github.io/golang/memory/leak/2016/02/11/go-mem-leak.html
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return d, fmt.Errorf("failed to read http result: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return d, fmt.Errorf("failed to fetch http: returned status code %d from %s", resp.StatusCode, url)
	}
	return d, nil
}
