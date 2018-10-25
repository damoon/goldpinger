package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"syscall/js"
	"time"

	"github.com/damoon/goldpinger/pkg"
)

func main() {
	ch := startNewModel()
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("recovered: %v\n", err)
				}
			}()
			err := updateJson(ch)
			displayError(err)
			t := time.NewTicker(1 * time.Second)
			for range t.C {
				err := updateJson(ch)
				displayError(err)
			}
		}()
	}
}

func displayError(err error) {
	el := js.Global().Get("document").Call("getElementById", "errors")
	if err == nil {
		el.Set("innerHTML", "")
	}
	el.Set("innerHTML", err.Error())
}

func updateJson(ch ModelAgent) error {
	resp, err := http.Get("./status.json")
	if err != nil {
		return fmt.Errorf("updatejson: %v", err)
	}
	defer resp.Body.Close()
	//	b, err := ioutil.ReadAll(resp.Body)
	//	if err != nil {
	//		return fmt.Errorf("updatejson: %v", err)
	//	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("updatejson: bad http status code %v", resp.StatusCode)
	}

	fetchedModel := &goldpinger.Model{}
	err = json.NewDecoder(resp.Body).Decode(fetchedModel)
	if err != nil {
		return fmt.Errorf("failed to decode json model: %s", err)
	}

	ch <- func(m *Model) {
		m.Model = goldpinger.Merge(*fetchedModel, m.Model)
	}

	return nil
}
