package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/damoon/goldpinger/pkg"

	k8sClient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	var namespace *string
	namespace = flag.String("hamespace", "goldpinger", "namespace to ping pods in")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	client, err := k8sClient.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		goldpinger.Status(w, r, client, namespace)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
