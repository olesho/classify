package stream

import (
	"fmt"
	"time"
)

// Timer calculates elapsed time
type Timer struct {
	ts time.Time
}

// NewTimer creates new timer
func NewTimer() *Timer {
	return &Timer{}
}

// Start starts timer (set to 0 seconds)
func (t *Timer) Start() {
	t.ts = time.Now()
}

// Check prints message and time elapsed and sets to 0 seconds
func (t *Timer) Check(msg string) {
	fmt.Println(msg, "at:", time.Since(t.ts).Seconds(), "seconds")
	t.ts = time.Now()
}
