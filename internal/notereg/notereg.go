// Package notereg provides a thread-safe registry for tracking active MIDI notes.
// It manages note-on/off states by channel and note number to prevent duplicate
// note events and ensure proper note cleanup. This is essential for polyphonic
// sequencing where multiple notes may be active simultaneously.
package notereg

import (
	"errors"
	"sync"
	"time"

	"gitlab.com/gomidi/midi/v2"
)

type NoteRegKey struct {
	Channel uint8
	Note    uint8
}

type Keyable interface {
	GetKey() NoteRegKey
	GetID() int
}

var noteMutex = &sync.Mutex{}

type NoteReg map[NoteRegKey]struct{}
type TimerReg map[NoteRegKey]*time.Timer

var noteReg = make(NoteReg)
var timerReg = make(TimerReg)

func Add(note Keyable) error {
	noteMutex.Lock()
	defer noteMutex.Unlock()
	nrk := note.GetKey()
	_, existing := noteReg[nrk]
	if existing {
		return errors.New("note already exists")
	}
	noteReg[nrk] = struct{}{}
	return nil
}

func AddKey(key NoteRegKey) error {
	noteMutex.Lock()
	defer noteMutex.Unlock()
	_, existing := noteReg[key]
	if existing {
		return errors.New("note already exists")
	}
	noteReg[key] = struct{}{}
	return nil
}

func AddTimer(key NoteRegKey, timer *time.Timer) error {
	noteMutex.Lock()
	defer noteMutex.Unlock()
	_, existing := timerReg[key]
	if existing {
		return errors.New("note already exists")
	}
	timerReg[key] = timer
	return nil
}

func Has(note Keyable) bool {
	_, existing := noteReg[note.GetKey()]
	return existing
}

func HasKey(key NoteRegKey) bool {
	_, existing := noteReg[key]
	return existing
}

func GetKey(msg midi.Message) NoteRegKey {
	var channel uint8
	var note uint8
	var velocity uint8
	if msg.Type().Is(midi.NoteOnMsg) {
		msg.GetNoteOn(&channel, &note, &velocity)
	} else if msg.Type().Is(midi.NoteOffMsg) {
		msg.GetNoteOff(&channel, &note, &velocity)
	} else {
		return NoteRegKey{
			Channel: 255,
			Note:    255,
		}
	}
	return NoteRegKey{
		Channel: channel,
		Note:    note,
	}
}

func RemoveKey(key NoteRegKey) {
	noteMutex.Lock()
	_, exists := timerReg[key]
	if exists {
		delete(timerReg, key)
	}
	defer noteMutex.Unlock()
	delete(noteReg, key)
}

func Remove(note Keyable) {
	noteMutex.Lock()
	defer noteMutex.Unlock()
	delete(noteReg, note.GetKey())
}

func Clear() []midi.Message {
	noteMutex.Lock()
	defer noteMutex.Unlock()
	values := make([]midi.Message, 0, len(noteReg))

	for n := range noteReg {
		values = append(values, midi.NoteOff(n.Channel, n.Note))
	}
	for key, timer := range timerReg {
		timer.Stop()
		delete(timerReg, key)
	}

	noteReg = make(NoteReg)
	return values
}
