package goldpinger

import (
	"math/rand"
	"reflect"
	"testing"
)

func Test_mergeHistories(t *testing.T) {
	type args struct {
		right History
		left  History
	}
	h1 := History{
		Measurement{Delay: 1, Error: "", Timestamp: 101},
		Measurement{Delay: 1, Error: "", Timestamp: 91},
		Measurement{Delay: 1, Error: "", Timestamp: 81},
		Measurement{Delay: 1, Error: "", Timestamp: 71},
		Measurement{Delay: 1, Error: "", Timestamp: 61},
		Measurement{Delay: 1, Error: "", Timestamp: 51},
		Measurement{Delay: 1, Error: "", Timestamp: 41},
		Measurement{Delay: 1, Error: "", Timestamp: 31},
		Measurement{Delay: 1, Error: "", Timestamp: 21},
		Measurement{Delay: 1, Error: "", Timestamp: 11},
	}
	h2 := History{
		Measurement{Delay: 1, Error: "", Timestamp: 102},
		Measurement{Delay: 1, Error: "", Timestamp: 92},
		Measurement{Delay: 1, Error: "", Timestamp: 82},
		Measurement{Delay: 1, Error: "", Timestamp: 72},
		Measurement{Delay: 1, Error: "", Timestamp: 62},
		Measurement{Delay: 1, Error: "", Timestamp: 52},
		Measurement{Delay: 1, Error: "", Timestamp: 42},
		Measurement{Delay: 1, Error: "", Timestamp: 32},
		Measurement{Delay: 1, Error: "", Timestamp: 22},
		Measurement{Delay: 1, Error: "", Timestamp: 12},
	}
	tests := []struct {
		name string
		args args
		want History
	}{
		{
			name: "zero history",
			args: args{
				right: History{},
				left:  History{},
			},
			want: History{},
		},
		{
			name: "right history",
			args: args{
				right: History{
					Measurement{Delay: 1, Error: "", Timestamp: 1},
				},
				left: History{},
			},
			want: History{
				Measurement{Delay: 1, Error: "", Timestamp: 1},
			},
		},
		{
			name: "left history",
			args: args{
				right: History{},
				left: History{
					Measurement{Delay: 1, Error: "", Timestamp: 1},
				},
			},
			want: History{
				Measurement{Delay: 1, Error: "", Timestamp: 1},
			},
		},
		{
			name: "same history",
			args: args{
				right: History{
					Measurement{Delay: 1, Error: "", Timestamp: 1},
				},
				left: History{
					Measurement{Delay: 1, Error: "", Timestamp: 1},
				},
			},
			want: History{
				Measurement{Delay: 1, Error: "", Timestamp: 1},
			},
		},
		{
			name: "updated history",
			args: args{
				right: History{
					Measurement{Delay: 1, Error: "", Timestamp: 2},
					Measurement{Delay: 1, Error: "", Timestamp: 1},
				},
				left: History{
					Measurement{Delay: 1, Error: "", Timestamp: 1},
				},
			},
			want: History{
				Measurement{Delay: 1, Error: "", Timestamp: 2},
				Measurement{Delay: 1, Error: "", Timestamp: 1},
			},
		},
		{
			name: "two histories",
			args: args{
				right: History{
					Measurement{Delay: 1, Error: "", Timestamp: 2},
				},
				left: History{
					Measurement{Delay: 1, Error: "", Timestamp: 1},
				},
			},
			want: History{
				Measurement{Delay: 1, Error: "", Timestamp: 2},
				Measurement{Delay: 1, Error: "", Timestamp: 1},
			},
		},
		{
			name: "full histories",
			args: args{
				right: h1,
				left:  h2,
			},
			want: History{
				Measurement{Delay: 1, Error: "", Timestamp: 102},
				Measurement{Delay: 1, Error: "", Timestamp: 101},
				Measurement{Delay: 1, Error: "", Timestamp: 92},
				Measurement{Delay: 1, Error: "", Timestamp: 91},
				Measurement{Delay: 1, Error: "", Timestamp: 82},
				Measurement{Delay: 1, Error: "", Timestamp: 81},
				Measurement{Delay: 1, Error: "", Timestamp: 72},
				Measurement{Delay: 1, Error: "", Timestamp: 71},
				Measurement{Delay: 1, Error: "", Timestamp: 62},
				Measurement{Delay: 1, Error: "", Timestamp: 61},
			},
		},
		{
			name: "full equal histories",
			args: args{
				right: h1,
				left:  h1,
			},
			want: h1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeHistories(tt.args.right, tt.args.left); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeHistories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelAccess_randomNode(t *testing.T) {
	rand.Seed(0)
	tests := []struct {
		name    string
		ch      ModelAccess
		nodes   []Node
		want    Node
		wantErr bool
	}{
		{
			name:    "empty node list",
			ch:      StartNewModel(rand.New(rand.NewSource(0))),
			nodes:   []Node{},
			wantErr: true,
		},
		{
			name: "one node",
			ch:   StartNewModel(rand.New(rand.NewSource(0))),
			nodes: []Node{
				{HostName: "hostName1", HostIP: "1.1.1.1", PodName: "podName1", PodIP: "1.1.1.2"},
			},
			want: Node{HostName: "hostName1", HostIP: "1.1.1.1", PodName: "podName1", PodIP: "1.1.1.2"},
		},
		{
			name: "two node",
			ch:   StartNewModel(rand.New(rand.NewSource(0))),
			nodes: []Node{
				{HostName: "hostName1", HostIP: "1.1.1.1", PodName: "podName1", PodIP: "1.1.1.2"},
				{HostName: "hostName2", HostIP: "2.1.1.1", PodName: "podName2", PodIP: "2.1.1.2"},
			},
			want: Node{HostName: "hostName1", HostIP: "1.1.1.1", PodName: "podName1", PodIP: "1.1.1.2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer close(tt.ch)
			for _, n := range tt.nodes {
				tt.ch.Add(n)
			}

			got, err := tt.ch.randomNode()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModelAccess.randomNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModelAccess.randomNode() = %v, want %v", got, tt.want)
			}
		})
	}
}
