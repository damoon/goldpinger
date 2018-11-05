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

func Measuring(ch ModelAccess, netClient *http.Client, hostname string) {
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		trg, err := ch.randomNode()
		if err != nil {
			Log("failed to ping: %v", err)
			return
		}
		url := fmt.Sprintf("http://%s/ok", trg.PodIP)
		go fetchHTTP(ch, netClient, trg.HostName, hostname, url)
	}
}

func fetchHTTP(ch ModelAccess, netClient *http.Client, target, source, url string) {
	d, err := measureHTTP(netClient, url)
	t := time.Now().UnixNano()
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	ch <- func(m Model) {
		_, ok := m.Status.Worldview[source]
		if !ok {
			m.Status.Worldview[source] = map[string]History{}
		}
		h := History{Measurement{
			Delay:     d,
			Error:     errMsg,
			Timestamp: t,
		}}
		m.Status.Worldview[source][target] = mergeHistories(m.Status.Worldview[source][target], h)
	}
}

func measureHTTP(netClient *http.Client, url string) (int64, error) {
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
