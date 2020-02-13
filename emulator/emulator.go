package emulator

import (
	"fmt"
	"time"
)

const (
	TIMERDURATION = time.Second / 60	// 60hz timer for sound and timer updates
	MEMORYSIZE = 4096 * 1024
)

// Emulator represents an instance of the Chip8 emulator.
type Emulator struct {
	mem [MEMORYSIZE]byte
	timerChan chan bool
}

func NewEmulator() *Emulator {
	return &Emulator{
		timerChan: nil,
	}
}

func startTicker(d time.Duration, f func()) chan bool {
	done := make(chan bool, 1)
	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				f()
			case <-done:
				return
			}
		}
	}()
	return done
}

// Start starts execution of the emulator.
func (e *Emulator) Start() {
	e.startTimer()
}

// start the background clock timer
func (e *Emulator) startTimer() {
	if (e.timerChan == nil) {
		e.timerChan = startTicker(TIMERDURATION, e.timerCallback)
	}
}

// stop the background clock timer
func (e *Emulator) stopTimer() {
	close(e.timerChan)
	// Let the goroutine finish
	time.Sleep(2 * TIMERDURATION)
	e.timerChan = nil
}

// Stop stops execution of the emulator.
func (e *Emulator) Stop() {
	e.stopTimer()
}

func (e *Emulator) Beep() {
	fmt.Printf("Beep")
}

func (e *Emulator) timerCallback() {
}