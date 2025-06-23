package grid

import (
	"fmt"
	"time"
)

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

type SequencerType uint8

const (
	SEQTYPE_TRIGGER SequencerType = iota
	SEQTYPE_POLYPHONY
)

type GridKey struct {
	Line uint8
	Beat uint8
}

func (gk GridKey) String() string {
	return fmt.Sprintf("Grid-%0.2d-%0.2d", gk.Line, gk.Beat)
}

func Compare(a GridKey, b GridKey) int {
	if a.Beat == b.Beat {
		return int(b.Line) - int(a.Line)
	} else {
		return int(a.Beat) - int(b.Beat)
	}
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

func (n Note) IncrementAccent(modifier int8, accentsLength uint8) Note {
	var newAccent = int8(n.AccentIndex) - modifier
	if newAccent >= 1 && newAccent < int8(accentsLength) {
		n.AccentIndex = uint8(newAccent)
	}
	return n
}

func (n Note) IncrementGate(modifier int8, gatesLength int) Note {
	var newGate = int8(n.GateIndex) + modifier
	if newGate >= 0 && int(newGate) < gatesLength {
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

func (n Note) IncrementRatchet(modifier int8) Note {
	currentRatchet := n.Ratchets.Length
	// TODO: remove hardcoded ratchets length
	var ratchetsLength int8 = 8
	newRatchet := int8(currentRatchet) + modifier
	if n.AccentIndex > 0 && n.Action == ACTION_NOTHING && newRatchet < ratchetsLength && int8(currentRatchet)+modifier >= 0 {
		n.Ratchets.Length = uint8(int8(currentRatchet) + modifier)
		n.Ratchets.SetRatchet(true, n.Ratchets.Length)
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

type MessageType uint8

const (
	MESSAGE_TYPE_NOTE MessageType = iota
	MESSAGE_TYPE_CC
)

type LineDefinition struct {
	Channel uint8
	Note    uint8
	MsgType MessageType
	Name    string
}

func (l *LineDefinition) IncrementChannel() {
	if l.Channel < 16 {
		l.Channel++
	}
}

func (l *LineDefinition) DecrementChannel() {
	if l.Channel > 1 {
		l.Channel--
	}
}

func (l *LineDefinition) IncrementNote() {
	if l.Note < 128 {
		l.Note++
	}
}

func (l *LineDefinition) DecrementNote() {
	if l.Note > 1 {
		l.Note--
	}
}

func (l *LineDefinition) IncrementMessageType() {
	l.MsgType = (l.MsgType + 1) % 2
}

func (l *LineDefinition) DecrementMessageType() {
	l.MsgType = (l.MsgType - 1) % 2
}
