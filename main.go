package main

import (
	"fmt"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

var wailsRuntime *wails.Runtime

func main() {

	s, err := NewSequencer(16, 0, 140.0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1280,
		Height: 940,
		Title:  "Sequensual",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})

	// app bindings
	app.Bind(s)

	app.Run()
}

// WailsInit instantiates the wails runtime and contains event listeners for the js app
func (s *Sequencer) WailsInit(runtime *wails.Runtime) error {
	s.Log = runtime.Log.New("Sequencer")
	wailsRuntime = runtime

	runtime.Events.On("play", func(...interface{}) {
		s.Start()
	})

	runtime.Events.On("stop", func(...interface{}) {
		s.Stop()
	})

	runtime.Events.On("changeStep", func(data ...interface{}) {

		newStep := Step{
			Number: 0,
			Active: false,
			Trig: &Trigger{
				Note:     0,
				Velocity: 100,
				Length:   1,
				Active:   false,
			},
		}
		for k, v := range data[0].(map[string]interface{}) {
			switch k {
			case "number":
				newStep.Number = int(v.(float64))
			case "status":
				newStep.Trig.Active = v.(bool)
			case "note":
				newStep.Trig.Note = int32(v.(float64))
				// case "length":
				// case "velocity":
			}
		}
		s.SetStep(newStep)
		fmt.Println(newStep.Trig)
		runtime.Events.Emit("updateSteps", s.Steps)
	})

	return nil
}
