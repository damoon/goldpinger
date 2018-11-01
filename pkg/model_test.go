package goldpinger

import (
	"reflect"
	"testing"
)

func Test_mergeHistories(t *testing.T) {
	type args struct {
		right History
		left  History
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
				right: History{
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
				},
				left: History{
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
				},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeHistories(tt.args.right, tt.args.left); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeHistories() = %v, want %v", got, tt.want)
			}
		})
	}
}
