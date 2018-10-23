package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"syscall/js"
	"time"
)

func main() {

	go updateJson()

	select {}
}

func updateJson() {
	t := time.NewTicker(1 * time.Second)
	for range t.C {
		//
		resp, err := http.Get("./status.json")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		//	enc := base64.StdEncoding.EncodeToString(b)

		el := js.Global().Get("document").Call("getElementById", "thing")
		el.Set("innerHTML", string(b))
	}
}
