package grid

import (
	"fmt"
	"time"
)

type GridKey struct {
	Line uint8
	Beat uint8
}

func (gk GridKey) String() string {
	return fmt.Sprintf("Grid-%0.2d-%0.2d", gk.Line, gk.Beat)
}

type Note struct {
	AccentIndex uint8
	Ratchets    Ratchet
	Action      Action
	GateIndex   uint8
	WaitIndex   uint8
}

func InitNote() Note {
	return Note{5, InitRatchet(), ACTION_NOTHING, 0, 0}
}

func InitActionNote(act Action) Note {
	return Note{0, InitRatchet(), act, 0, 0}
}

func (n Note) IncrementAccent(modifier int8) Note {
	var newAccent = int8(n.AccentIndex) - modifier
	// TODO: remove hardcoded accent length
	accentlength := 8
	if newAccent >= 1 && newAccent < int8(accentlength) {
		n.AccentIndex = uint8(newAccent)
	}
	return n
}

func (n Note) IncrementGate(modifier int8) Note {
	var newGate = int8(n.GateIndex) + modifier
	// TODO: remove hardcoded gates length
	gateslength := 8
	if newGate >= 0 && newGate < int8(gateslength) {
		n.GateIndex = uint8(newGate)
	}
	return n
}

func (n Note) IncrementWait(modifier int8) Note {
	var newWait = int8(n.WaitIndex) + modifier
	// TODO: remove hardcoded waits length
	waitslength := 8
	if newWait >= 0 && newWait < int8(waitslength) {
		n.WaitIndex = uint8(newWait)
	}
	return n
}

type Pattern map[GridKey]Note

type Ratchet struct {
	Hits   uint8
	Length uint8
	Span   uint8
}

func InitRatchet() Ratchet {
	return Ratchet{
		// We always start with one Ratchet enabled
		Hits: boolsToUint8([8]bool{true, false, false, false, false, false, false, false}),
	}
}

func boolsToUint8(bools [8]bool) uint8 {
	var result uint8
	for i := 0; i < 8; i++ {
		if bools[i] {
			result = result | (1 << i)
		}
	}
	return result
}

func (r Ratchet) Interval(beatInterval time.Duration) time.Duration {
	return (beatInterval * time.Duration(r.GetSpan())) / time.Duration(r.Length+1)
}

func (r Ratchet) GetSpan() uint8 {
	return r.Span + 1
}

func (r *Ratchet) SetRatchet(value bool, index uint8) {
	if value {
		// (1 << 1) == 010
		// 001 | 010 == 011
		r.Hits = r.Hits | (1 << index)
	} else {
		// 011 & 101 == 001
		r.Hits = r.Hits & ((1 << index) ^ 255)
	}
}

func (r *Ratchet) Toggle(index uint8) {
	// 011 ^ 010  == 001
	r.Hits = (r.Hits ^ (1 << index))
}

func (r Ratchet) HitAt(index uint8) bool {
	return (r.Hits & (1 << index)) != 0
}

type Action uint8

const (
	ACTION_NOTHING Action = iota
	ACTION_LINE_RESET
	ACTION_LINE_REVERSE
	ACTION_LINE_SKIP_BEAT
	ACTION_RESET
	ACTION_LINE_BOUNCE
	ACTION_LINE_DELAY
)
