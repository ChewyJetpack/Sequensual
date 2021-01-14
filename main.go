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

func (s *Sequencer) wailsInit(runtime *wails.Runtime) error {

	runtime.Events.On("play", func(...interface{}) {
		s.Start()
	})

	runtime.Events.On("toggleStep", func(...interface{}) {
		// newData := data[0].([]map[string]string)
		//fmt.Println(data)
		// s.SetStep(1, true)
		runtime.Events.Emit("updateSteps", s.Steps)
	})

	return nil
}
