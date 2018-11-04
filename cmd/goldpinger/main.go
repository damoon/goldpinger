package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/damoon/goldpinger/pkg"
	"github.com/damoon/goldpinger/pkg/k8s"
	"github.com/lpar/gzipped"
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
	random := rand.New(rand.NewSource(*seed))
	ch := goldpinger.StartNewModel(random)

	go k8s.PodListSyncing(*kubeconfig, *namespace, ch)
	go goldpinger.Gossiping(ch)
	go goldpinger.Measuring(ch, *hostName)

	m := http.NewServeMux()
	m.HandleFunc("/ok", goldpinger.OK)
	m.HandleFunc("/status.json", func(w http.ResponseWriter, r *http.Request) {
		goldpinger.PublishStatus(w, r, ch)
	})
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Location", "/index.html")
			w.WriteHeader(http.StatusMovedPermanently)
			return
		}
		gzipped.FileServer(http.Dir("./public/")).ServeHTTP(w, r)
	})
	log.Printf("start to listen on %v", *addr)
	server := &http.Server{Addr: *addr, Handler: m}
	log.Fatalln(server.ListenAndServe())
}
