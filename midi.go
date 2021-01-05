package main

import (
	"bytes"
	"fmt"
	"io"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/realtime"
	"gitlab.com/gomidi/midi/midireader"
	"gitlab.com/gomidi/midi/midiwriter"
)

func main() {
	var bf bytes.Buffer

	wr := midiwriter.New(&bf)
	wr.Write(realtime.Channel2.Pitchbend(5000))
	wr.Write(realtime.Channel2.NoteOn(65, 90))
	wr.Write(realtime.Reset)
	wr.Write(realtime.Channel2.NoteOff(65))

	rthandler := func(m realtime.Message) {
		fmt.Printf("Realtime: %s\n", m)
	}

	rd := midireader.New(bytes.NewReader(bf.Bytes()), rthandler)

	var m midi.Message
	var err error

	for {
		m, err = rd.Read()

		// breaking at least with io.EOF
		if err != nil {
			break
		}

		// inspect
		fmt.Println(m)

		switch v := m.(type) {
		case NoteOn:
			fmt.Printf("NoteOn at channel %v: key: %v velocity: %v\n", v.Channel(), v.Key(), v.Velocity())
		case NoteOff:
			fmt.Printf("NoteOff at channel %v: key: %v\n", v.Channel(), v.Key())
		}

	}

	if err != io.EOF {
		panic("error: " + err.Error())
	}

	// Output:
	// channel.Pitchbend channel 2 value 5000 absValue 13192
	// channel.NoteOn channel 2 key 65 velocity 90
	// NoteOn at channel 2: key: 65 velocity: 90
	// Realtime: Reset
	// channel.NoteOff channel 2 key 65
	// NoteOff at channel 2: key: 65
}
