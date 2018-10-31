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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeHistories(tt.args.right, tt.args.left); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeHistories() = %v, want %v", got, tt.want)
			}
		})
	}
}
