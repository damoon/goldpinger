package main

import (
	"log"
	"net/http"

	"github.com/damoon/goldpinger/pkg"
)

func main() {
	http.HandleFunc("/", goldpinger.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
