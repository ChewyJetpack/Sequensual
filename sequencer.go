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

type Step struct {
	Number int
	Active bool
	Trig   *Trigger
}

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

	go func() {
		ppqnCount := 0

		for {
			select {
			case <-s.Timer.Pulses:
				if !s.Running {
					ppqnCount = 0
					s.Beat = 0
					s.Timer.Done <- true
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

func (s *Sequencer) Stop() {
	s.Running = false
}

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

// ProcessAudio is the callback function for the portaudio stream
// Attached the the current Sequencer.
// Writes all active Track Samples to the output buffer
// At the playhead for each track.
func (s *Sequencer) ProcessAudio(out []float32) {
	for i := range out {
		var data float32

		if data > 1.0 {
			data = 1.0
		}

		out[i] = data
	}
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
