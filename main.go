package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/damoon/goldpinger/pkg"
)

func main() {

	var hostName = flag.String("hostName", "", "name of node the pod is running on")
	var seed = flag.Int64("seed", time.Now().UnixNano(), "seed to use for random number generators")
	var addr = flag.String("addr", ":80", "address to listen on")
	var kubeconfig = flag.String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	var namespace = flag.String("namespace", "goldpinger", "namespace to ping pods in")
	flag.Parse()

	log.Printf("hostName: %v\n", *hostName)
	log.Printf("seed: %d\n", *seed)
	log.Printf("addr: %v\n", *addr)
	log.Printf("kubeconfig: %v\n", *kubeconfig)
	log.Printf("namespace: %v\n", *namespace)

	if *hostName == "" {
		log.Fatalf("hostName is not set\n")
	}

	r := rand.New(rand.NewSource(*seed))
	log.Printf("starting goldpinger")
	pinger := goldpinger.StartNew(*hostName, *kubeconfig, *namespace, r, log.Printf)

	m := http.NewServeMux()
	m.HandleFunc("/ok", goldpinger.OK)
	m.HandleFunc("/status.json", pinger.Status)
	m.HandleFunc("/", http.FileServer(http.Dir("./static/")).ServeHTTP)
	log.Printf("start to listen on %v", *addr)
	server := &http.Server{Addr: *addr, Handler: m}
	log.Fatalln(server.ListenAndServe())
}
