package goldpinger

import (
	"math/rand"
	"time"
)

func randFloat(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func (p *Pinger) mock() *Model {
	return &Model{
		"host1": &Source{
			Target: Target{
				HostName: "host1",
				HostIP:   "1.2.3.4",
				PodName:  "pod1",
				PodIP:    "1.2.3.5",
			},
			Measurements: map[string]*Measurement{
				"host1": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
				"host2": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
				"host3": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
				"host4": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
			},
		},
		"host2": &Source{
			Target: Target{
				HostName: "host2",
				HostIP:   "2.2.3.4",
				PodName:  "pod2",
				PodIP:    "2.2.3.5",
			},
			Measurements: map[string]*Measurement{
				"host1": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
				"host2": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
				"host3": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
				"host4": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
			},
		},
		"host3": &Source{
			Target: Target{
				HostName: "host3",
				HostIP:   "3.2.3.4",
				PodName:  "pod3",
				PodIP:    "3.2.3.5",
			},
			Measurements: map[string]*Measurement{
				"host1": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
				"host2": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
				"host3": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
				"host4": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
			},
		},
		"host4": &Source{
			Target: Target{
				HostName: "host4",
				HostIP:   "4.2.3.4",
				PodName:  "pod4",
				PodIP:    "4.2.3.5",
			},
			Measurements: map[string]*Measurement{
				"host1": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
				"host2": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
				"host3": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
				"host4": &Measurement{Error: "", Delay: rand.Int63n(10000), Timestamp: time.Now().Unix()},
			},
		},
	}
}
