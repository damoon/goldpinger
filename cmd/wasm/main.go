package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/damoon/goldpinger/pkg"
)

func main() {
	ch := startNewModel()
	t := time.NewTicker(1 * time.Second)
	for ; true; <-t.C {
		err := updateModel(ch)
		if err != nil {
			ch <- func(m *Model) {
				m.FetchError = err.Error()
			}
			fmt.Printf("update model: %v\n", err)
		}
	}
}

func updateModel(ch ModelAgent) error {
	resp, err := http.Get("./status.json")
	if err != nil {
		return fmt.Errorf("fetch failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("bad http status code %v", resp.StatusCode)
	}

	update := &goldpinger.Status{}
	err = json.NewDecoder(resp.Body).Decode(update)
	if err != nil {
		return fmt.Errorf("failed to decode json model: %s", err)
	}

	ch <- func(m *Model) {
		m.FetchError = ""
		m.Status = goldpinger.MergeStatus(*update, m.Status)
	}

	return nil
}
