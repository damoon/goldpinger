package main

import (
	"fmt"
	"syscall/js"

	"github.com/damoon/goldpinger/pkg"
	"github.com/mohae/deepcopy"
)

type Model struct {
	Model goldpinger.Model
	Error string
}

type ModelAgent chan<- func(m *Model)

func startNewModel() ModelAgent {
	c := make(chan func(m *Model))
	m := &Model{
		Model: goldpinger.Model{
			Nodes:        []*goldpinger.Node{},
			Measurements: map[string]map[string]*goldpinger.Measurement{},
		},
		Error: "",
	}
	go func() {
		for f := range c {
			f(m)
			el := js.Global().Get("document").Call("getElementById", "thing")
			el.Set("innerHTML", m.Render())
		}
	}()
	return c
}

func (m *Model) Render() string {
	return fmt.Sprintf("%s", m)
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
