package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"syscall/js"
	"time"
)

func main() {
	for {
		updateJsonLoop()
	}
}

func updateJsonLoop() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("recovered: %v\n", err)
		}
	}()
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		updateJson()
	}
}

func updateJson() {
	resp, err := http.Get("./status.json")
	if err != nil {
		fmt.Printf("updatejson: %v\n", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("updatejson: %v\n", err)
	}

	el := js.Global().Get("document").Call("getElementById", "thing")
	el.Set("innerHTML", string(b))
}
