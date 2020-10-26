// The Test package is used for testing rus.
// It provides a simple hooks which register logged messages.
package test

import (
	"io/ioutil"
	"sync"

	"github.com/shupkg/trier/rus"
)

// Hook is a hook designed for dealing with logs in test scenarios.
type Hook struct {
	// Entries is an array of all entries that have been received by this hook.
	// For safe access, use the AllEntries() method, rather than reading this
	// value directly.
	Entries []rus.Entry
	mu      sync.RWMutex
}

// NewGlobal installs a test hook for the global logger.
func NewGlobal() *Hook {

	hook := new(Hook)
	rus.AddHook(hook)

	return hook

}

// NewLocal installs a test hook for a given local logger.
func NewLocal(logger *rus.Logger) *Hook {

	hook := new(Hook)
	logger.Hooks.Add(hook)

	return hook

}

// NewNullLogger creates a discarding logger and installs the test hook.
func NewNullLogger() (*rus.Logger, *Hook) {

	logger := rus.New()
	logger.Out = ioutil.Discard

	return logger, NewLocal(logger)

}

func (t *Hook) Fire(e *rus.Entry) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Entries = append(t.Entries, *e)
	return nil
}

func (t *Hook) Levels() []rus.Level {
	return rus.AllLevels
}

// LastEntry returns the last entry that was logged or nil.
func (t *Hook) LastEntry() *rus.Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	i := len(t.Entries) - 1
	if i < 0 {
		return nil
	}
	return &t.Entries[i]
}

// AllEntries returns all entries that were logged.
func (t *Hook) AllEntries() []*rus.Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	// Make a copy so the returned value won't race with future log requests
	entries := make([]*rus.Entry, len(t.Entries))
	for i := 0; i < len(t.Entries); i++ {
		// Make a copy, for safety
		entries[i] = &t.Entries[i]
	}
	return entries
}

// Reset removes all Entries from this test hook.
func (t *Hook) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Entries = make([]rus.Entry, 0)
}
