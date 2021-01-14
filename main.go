package main

import (
	"fmt"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

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

func (s *Sequencer) WailsInit(runtime *wails.Runtime) error {
	s.Log = runtime.Log.New("Sequencer")

	runtime.Events.On("play", func(...interface{}) {
		s.Start()
	})

	runtime.Events.On("changeStep", func(data ...interface{}) {

		newStep := Step{
			Active: false,
			Trig: &Trigger{
				Note: int32(50),
			},
		}
		for k, v := range data[0].(map[string]interface{}) {
			if k == "stepNumber" {
				newStep.Number = int(v.(float64))
			} else if k == "trigStatus" {
				newStep.Trig.Active = v.(bool)
			}
		}
		s.SetStep(newStep)
		fmt.Println(newStep)
		runtime.Events.Emit("updateSteps", s.Steps)
	})

	return nil
}
