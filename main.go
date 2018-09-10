package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/damoon/goldpinger/pkg"

	k8sClient "k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"
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
		log.Fatalf("hostName is not set")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("failed to load config for kubernetes client: %v", err)
	}
	client, err := k8sClient.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create kubernetes client: %v", err)
	}

	pods := client.CoreV1().Pods(*namespace)
	r := rand.New(rand.NewSource(*seed))
	log.Printf("starting goldpinger")
	pinger := goldpinger.StartNew(*hostName, pods, r)

	m := http.NewServeMux()
	m.HandleFunc("/_goroutines", getGoroutinesCountHandler)
	m.HandleFunc("/ok", goldpinger.OK)
	m.HandleFunc("/status.json", pinger.Status)
	m.HandleFunc("/", http.FileServer(http.Dir("./static/")).ServeHTTP)
	server := &http.Server{Addr: *addr, Handler: m}
	log.Printf("start to listen on %v", *addr)
	log.Fatalln(server.ListenAndServe())
}

func getGoroutinesCountHandler(w http.ResponseWriter, r *http.Request) {
	count := runtime.NumGoroutine()
	w.Write([]byte(strconv.Itoa(count)))
}
