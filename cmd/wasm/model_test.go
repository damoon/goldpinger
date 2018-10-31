package main

import (
	"testing"

	"github.com/damoon/goldpinger/pkg"
)

func TestModel_renderMeasurement(t *testing.T) {
	type fields struct {
		Model      goldpinger.Model
		FetchError string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "incomplete Worldview",
			fields: fields{
				Model: goldpinger.Model{
					Participants: []*goldpinger.Node{
						&goldpinger.Node{
							HostIP:   "1.1.1.1",
							HostName: "nodeOne",
							PodIP:    "1.1.1.2",
							PodName:  "podOne",
						},
						&goldpinger.Node{
							HostIP:   "2.1.1.1",
							HostName: "nodeTwo",
							PodIP:    "2.1.1.2",
							PodName:  "podTwo",
						},
						&goldpinger.Node{
							HostIP:   "3.1.1.1",
							HostName: "nodeThree",
							PodIP:    "3.1.1.2",
							PodName:  "podThree",
						},
					},
					Worldview: map[string]map[string]goldpinger.History{
						"nodeOne": map[string]goldpinger.History{
							"nodeOne": []*goldpinger.Measurement{
								&goldpinger.Measurement{
									Delay:     1,
									Error:     "",
									Timestamp: 1,
								},
							},
							"nodeTwo": []*goldpinger.Measurement{
								&goldpinger.Measurement{
									Delay:     1,
									Error:     "",
									Timestamp: 1,
								},
							},
							"nodeThree": []*goldpinger.Measurement{
								&goldpinger.Measurement{
									Delay:     1,
									Error:     "",
									Timestamp: 1,
								},
							},
						},
						"nodeTwo": map[string]goldpinger.History{
							"nodeOne": []*goldpinger.Measurement{},
							"nodeTwo": []*goldpinger.Measurement{},
							"nodeThree": []*goldpinger.Measurement{
								&goldpinger.Measurement{
									Delay:     1,
									Error:     "some Error",
									Timestamp: 1,
								},
							},
						},
						"nodeThree": map[string]goldpinger.History{
							"nodeOne":   []*goldpinger.Measurement{},
							"nodeTwo":   []*goldpinger.Measurement{},
							"nodeThree": []*goldpinger.Measurement{},
						},
					},
				},
				FetchError: "",
			},
			want: `<table>
<tr>
	<td/>
	<td ><div class="to">to nodeOne<div></td>
	<td ><div class="to">to nodeTwo<div></td>
	<td ><div class="to">to nodeThree<div></td>
	</tr>
<tr>
	<td>from nodeOne</td>
	<td class="ping" style="color:rgb(0, 255, 0);">0.0</td>
	<td class="ping" style="color:rgb(0, 255, 0);">0.0</td>
	<td class="ping" style="color:rgb(0, 255, 0);">0.0</td>
	</tr>
<tr>
	<td>from nodeTwo</td>
	<td class="empty ping" />
	<td class="empty ping" />
	<td><img title="some Error" src="./alert.png" /></td>
	</tr>
<tr>
	<td>from nodeThree</td>
	<td class="empty ping" />
	<td class="empty ping" />
	<td class="empty ping" />
	</tr>
</table>
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Model{
				Model:      tt.fields.Model,
				FetchError: tt.fields.FetchError,
			}
			if got := m.renderMeasurement(); got != tt.want {
				t.Errorf("Model.renderMeasurement() = %v, want %v", got, tt.want)
			}
		})
	}
}
