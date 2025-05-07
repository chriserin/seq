package theory

import (
	"sync"
)

// Registry is a singleton chord registry
type Registry struct {
	chords map[int]Chord
	nextID int
	mutex  sync.RWMutex
}

var (
	instance *Registry
	once     sync.Once
)

// GetRegistry returns the singleton instance of the chord registry
func GetRegistry() *Registry {
	once.Do(func() {
		instance = &Registry{
			chords: make(map[int]Chord),
			nextID: 1,
		}
	})
	return instance
}

func RegisterChord(chord Chord) int {
	return GetRegistry().RegisterChord(chord)
}

// RegisterChord adds a chord to the registry and returns its ID
func (r *Registry) RegisterChord(chord Chord) int {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	id := r.nextID
	r.chords[id] = chord
	r.nextID++
	return id
}

func GetChord(id int) (Chord, bool) {
	return GetRegistry().GetChord(id)
}

// GetChord retrieves a chord from the registry by its ID
func (r *Registry) GetChord(id int) (Chord, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	chord, exists := r.chords[id]
	return chord, exists
}

func RemoveChord(id int) bool {
	return GetRegistry().RemoveChord(id)
}

// RemoveChord removes a chord from the registry by its ID
func (r *Registry) RemoveChord(id int) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.chords[id]; !exists {
		return false
	}

	delete(r.chords, id)
	return true
}

func UpdateChord(chord Chord) bool {
	return GetRegistry().UpdateChord(chord.Id, chord)
}

func (r *Registry) UpdateChord(id int, chord Chord) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.chords[id]; exists {
		r.chords[id] = chord
	} else {
		panic("Updated chord should exist")
	}

	return true
}

// Clear removes all chords from the registry
func (r *Registry) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.chords = make(map[int]Chord)
	r.nextID = 1
}
