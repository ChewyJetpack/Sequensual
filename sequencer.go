package main

import (
	"fmt"

	"github.com/hypebeast/go-osc/osc"
	"github.com/wailsapp/wails"
)

// Sequencer describes the mechanism that
// Triggers and synchronizses a Pattern for audio playback.
type Sequencer struct {
	Timer   *Timer
	Running bool
	Beat    int
	Channel int
	Steps   *map[int]*Step
	Length  int
	Osc     *osc.Client
	Log     *wails.CustomLogger
}

// Step contains all the data needed for a single sequencer step, aside from the trig
type Step struct {
	Number int
	Active bool
	Trig   *Trigger
}

// Trigger - an active component of a step that contains any relevant notation or other data to be sent to a device
type Trigger struct {
	Note     int32
	Velocity float32
	Length   int
	Active   bool
}

// NewSequencer creates and returns a pointer to a New Sequencer.
// Returns an error if there is one encountered
// During initializing portaudio, or the default stream
func NewSequencer(length, channel int, tempo float32) (*Sequencer, error) {

	s := &Sequencer{
		Timer:   NewTimer(tempo),
		Running: false,
		Beat:    0,
		Channel: channel,
		Steps:   MakeSteps(length),
		Length:  length,
		//Osc:     osc.NewClient("elk-pi.local", 24024),
	}

	return s, nil
}

// MakeSteps creates and returns a map of key value pairs,
// containing integer keys, one for each step, and a pointers to steps
func MakeSteps(length int) *map[int]*Step {
	stps := make(map[int]*Step, length)

	for i := 0; i < length; i++ {
		stps[i] = &Step{
			Number: i,
			Active: false,
			Trig: &Trigger{
				Note: int32(i + 50),
			},
		}
	}

	return &stps
}

// Start starts the sequencer.
// Starts counting the Pulses Per Quarter note for the given BPM.
// Triggers samples based on each 16th note that is triggered.
func (s *Sequencer) Start() {
	s.Running = true
	wailsRuntime.Events.Emit("running", s.Running)

	go func() {
		ppqnCount := 0

		for {
			select {
			case <-s.Timer.Pulses:
				if !s.Running {
					ppqnCount = 0
					s.Beat = 0
					//s.Timer.Done <- true
					break
				}
				ppqnCount++

				// Trigger playback after 6 pulses, which is 16th note precision
				// TODO add in time signatures
				if ppqnCount%(int(Ppqn)/4) == 0 {
					activeStep((*s.Steps)[s.Beat], s)
				}

				// Reset the ppqn Count and Beat count after 4 bars of quarter notes
				// TODO add in different kinds of time signatures
				if ppqnCount == (int(Ppqn) * 4) {
					ppqnCount = 0
					s.Beat = 0
				}

			}
		}
	}()

	go s.Timer.Start()
}

// Stop stops the sequencer Start() loop
func (s *Sequencer) Stop() {
	s.Running = false
	wailsRuntime.Events.Emit("running", s.Running)
}

// activeStep sets the current step as active, emits data to the front end, and executes the step trig
func activeStep(stp *Step, s *Sequencer) {
	stp.Active = true

	wailsRuntime.Events.Emit("activeStep", stp.Number)

	if stp.Trig.Active {
		fmt.Println("Step:", stp.Number, "Trig note:", stp.Trig.Note)
		// commented out for pi-free development
		//go oscNote(stp.Trig.Note, 1, s.Osc)
	} else {
		fmt.Println("Step:", stp.Number, "No Trig.")
	}

	s.Beat++
	stp.Active = false
}

// GetLength is a utility method for the front end, to get the length of the sequence
func (s *Sequencer) GetLength() int {
	return s.Length
}

// GetSteps is a utility method for the front end, to get the Steps map
func (s *Sequencer) GetSteps() *map[int]*Step {
	return s.Steps
}

// SetStep is a setter that can be called in response to the 'toggleStep' event from the front end
func (s *Sequencer) SetStep(st Step) {
	(*s.Steps)[st.Number] = &st
}
