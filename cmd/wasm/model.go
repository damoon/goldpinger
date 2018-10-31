package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"syscall/js"
	"text/template"

	"github.com/damoon/goldpinger/pkg"
	"github.com/mohae/deepcopy"
)

type Model struct {
	Model      goldpinger.Model
	FetchError string
}

type ModelAgent chan<- func(m *Model)

func startNewModel() ModelAgent {
	c := make(chan func(m *Model))
	m := &Model{
		Model: goldpinger.Model{
			Participants: []*goldpinger.Node{},
			Worldview:    map[string]map[string]goldpinger.History{},
		},
		FetchError: "",
	}
	go func() {
		for f := range c {
			f(m)
			d := js.Global().Get("document")
			d.Call("getElementById", "measurement").Set("innerHTML", m.renderMeasurement())
			d.Call("getElementById", "fetch-error").Set("innerHTML", m.renderFetchError())
			//d.Call("getElementById", "json").Set("innerHTML", m.renderJSON())
		}
	}()
	return c
}

const measurementsTemplate = `<table>
<tr>
	<td/>
	{{range $to := $.Participants -}}
		<td ><div class="to">to {{ $to.HostName }}<div></td>
	{{ end -}}
</tr>
{{- range $from := $.Participants }}
<tr>
	<td>from {{ $from.HostName }}</td>
	{{ range $to := $.Participants -}}
	{{ $history := index $.Worldview $from.HostName $to.HostName -}}
	{{ if not $history -}}
	<td class="empty ping" />
	{{ else -}}
	{{ $measurement := index $history 0 -}}
	{{ if not $measurement -}}
	<td class="empty ping" />
	{{ else -}}
	{{ if ne $measurement.Error "" -}}
	<td><img title="{{ $measurement.Error }}" src="./alert.png" /></td>
	{{ else -}}
	{{ $delay := NanoToMilli $measurement.Delay -}}
	<td class="ping" style="color:{{ Color $delay }};">{{ printf "%.1f" $delay }}</td>
	{{ end -}}
	{{- end -}}
	{{- end -}}
	{{- end -}}
</tr>
{{- end }}
</table>
`

func (m *Model) renderMeasurement() string {

	nanoToMlli := func(n int64) float64 {
		return float64(n) / 1000000
	}

	color := func(delay float64) string {
		r := 0
		g := 0
		b := 0

		if delay > 4 {
			r = 255
		} else {
			r = int(255 / 4 * delay)
		}

		if delay < 4 {
			g = 255
		} else {
			if delay > 8 {
				g = 0
			} else {
				g = int(255 / 1 * (8 - delay))
			}
		}
		return "rgb(" + strconv.Itoa(r) + ", " + strconv.Itoa(g) + ", " + strconv.Itoa(b) + ")"
	}

	fns := template.FuncMap{"NanoToMilli": nanoToMlli, "Color": color}

	tpl, err := template.New("measurements").Funcs(fns).Parse(measurementsTemplate)
	if err != nil {
		return fmt.Sprintf("failed to parse measurements template: %v", err)
	}

	b := &bytes.Buffer{}
	err = tpl.Execute(b, m.Model)
	if err != nil {
		return fmt.Sprintf("failed to render measurements template: %v", err)
	}

	return b.String()
}

func (m *Model) renderFetchError() string {
	return m.FetchError
}

func (m *Model) renderJSON() string {
	json, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("failed to marshal model to json: %v", err)
	}
	return string(json)
}

func model(ch ModelAgent) Model {
	r := make(chan Model)
	ch <- func(m *Model) {
		c := deepcopy.Copy(*m)
		r <- c.(Model)
		close(r)
	}
	return <-r
}
