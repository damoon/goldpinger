package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/damoon/goldpinger/pkg"
	"github.com/damoon/goldpinger/pkg/k8s"
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

	//goldpinger.Log = log.Printf

	log.Printf("starting goldpinger")
	r := rand.New(rand.NewSource(*seed))
	ch := goldpinger.StartNewModel()
	nodeSelector := goldpinger.NewRandomNode(r)

	go k8s.PodListSyncing(*kubeconfig, *namespace, ch)
	go goldpinger.Gossiping(ch, nodeSelector)
	go goldpinger.Measuring(ch, nodeSelector, *hostName)

	m := http.NewServeMux()
	m.HandleFunc("/ok", goldpinger.OK)
	m.HandleFunc("/status.json", func(w http.ResponseWriter, r *http.Request) {
		goldpinger.Status(w, r, ch)
	})
	m.HandleFunc("/", http.FileServer(http.Dir("./public/")).ServeHTTP)
	log.Printf("start to listen on %v", *addr)
	server := &http.Server{Addr: *addr, Handler: m}
	log.Fatalln(server.ListenAndServe())
}
