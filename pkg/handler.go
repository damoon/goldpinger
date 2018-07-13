package goldpinger

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type model []host

type host struct {
	Source string `json:"source"`
	Pings  []ping `json:"pings"`
}

type ping struct {
	Target    string `json:"target"`
	Delay     int    `json:"delay"`
	Timestamp int64  `json:"timestamp"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	hosts := []host{
		host{
			Source: "host1",
			Pings: []ping{
				ping{Target: "host1", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host2", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host3", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host4", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
			},
		},
		host{
			Source: "host2",
			Pings: []ping{
				ping{Target: "host1", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host2", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host3", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host4", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
			},
		},
		host{
			Source: "host3",
			Pings: []ping{
				ping{Target: "host1", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host2", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host3", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host4", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
			},
		},
		host{
			Source: "host4",
			Pings: []ping{
				ping{Target: "host1", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host2", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host3", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
				ping{Target: "host4", Delay: rand.Int() % 100, Timestamp: time.Now().Unix()},
			},
		},
	}

	json, err := json.Marshal(hosts)
	if err != nil {
		log.Fatalf("failed to marshal pings to json: %v", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, string(json))
}
