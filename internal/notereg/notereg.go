package notereg

import (
	"errors"
	"maps"
	"sync"
)

type NoteRegKey struct {
	Channel uint8
	Note    uint8
}

type Keyable interface {
	GetKey() NoteRegKey
}

var noteMutex = &sync.Mutex{}

type NoteReg map[NoteRegKey]Keyable

var noteReg = make(NoteReg)

func Add(note Keyable) error {
	noteMutex.Lock()
	defer noteMutex.Unlock()
	nrk := note.GetKey()
	_, existing := noteReg[nrk]
	if existing {
		return errors.New("note already exists")
	}
	noteReg[nrk] = note
	return nil
}

func Has(note Keyable) bool {
	_, existing := noteReg[note.GetKey()]
	return existing
}

func Remove(note Keyable) {
	noteMutex.Lock()
	defer noteMutex.Unlock()
	delete(noteReg, note.GetKey())
}

func Clear() []Keyable {
	noteMutex.Lock()
	defer noteMutex.Unlock()
	values := make([]Keyable, 0, len(noteReg))

	for n := range maps.Values(noteReg) {
		values = append(values, n)
	}

	noteReg = make(NoteReg)
	return values
}
