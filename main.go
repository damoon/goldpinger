package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/damoon/goldpinger/pkg"

	k8sClient "k8s.io/client-go/kubernetes"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	var seed = flag.Int64("seed", time.Now().UnixNano(), "seed to use for random number generators")
	var addr = flag.String("addr", ":80", "address to listen on")
	var kubeconfig = flag.String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	var namespace = flag.String("namespace", "goldpinger", "namespace to ping pods in")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("failed to load config for kubernetes client: %v", err)
	}
	client, err := k8sClient.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create kubernetes client: %v", err)
	}

	p, err := client.CoreV1().Pods(*namespace).Watch(meta_v1.ListOptions{})
	if err != nil {
		log.Fatalf("failed to watch pods in namespace %s: %v", *namespace, err)
	}
	n, err := client.CoreV1().Nodes().Watch(meta_v1.ListOptions{})
	if err != nil {
		log.Fatalf("failed to watch nodes: %v", err)
	}
	r := rand.New(rand.NewSource(seed))
	u := goldpinger.NewPinger(p.ResultChan(), n.ResultChan(), r)
	log.Printf("starting pinger")
	u.Start()

	r := http.NewServeMux()
	r.HandleFunc("/ok", goldpinger.OK)
	r.HandleFunc("/state.json", func(w http.ResponseWriter, r *http.Request) {
		goldpinger.Status(w, r, u)
	})
	r.HandleFunc("/", http.FileServer(http.Dir("./static/")).ServeHTTP)
	server := &http.Server{Addr: *addr, Handler: r}
	log.Printf("start to listen on %v", *addr)
	log.Fatalln(server.ListenAndServe())
}
