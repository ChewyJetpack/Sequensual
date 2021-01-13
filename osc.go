package main

import (
	"fmt"
	"time"

	"github.com/hypebeast/go-osc/osc"
)

func oscNote(note, length int32, client *osc.Client) {
	noteOn := osc.NewMessage("/keyboard_event/main", string("note_on"), int32(0), note, float32(1))
	noteOff := osc.NewMessage("/keyboard_event/main", string("note_off"), int32(0), note, float32(1))
	client.Send(noteOn)
	fmt.Println(noteOn)
	time.Sleep(time.Second)
	client.Send(noteOff)
}
