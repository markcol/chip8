package emulator

import (
	"fmt"
	"time"
)

const (
	// TimerFrequency holds the frequency of the sound and timer clocks (60hz).
	TimerFrequency = time.Second / 60

	// MemorySize holds the amount of RAM available in the emulator.
	MemorySize = 4096

	// DisplayHeight holds the number of lines available in the display.
	DisplayHeight = 32

	// DisplayWidth holds the number of columns available in the display.
	DisplayWidth = 64

	// Registers holds the number of registers available in the Emulator.
	Registers = 16

	// StackSize holds the size of the stack in the Emulator (max call depth).
	StackSize = 16
)

// Emulator represents an instance of the Chip8 emulator.
type Emulator struct {
	mem       [MemorySize]byte
	display   [DisplayWidth * DisplayHeight]byte
	registers [Registers]byte
	stack     [StackSize]uint16
	pc        uint16
	sp        byte
	timerChan chan bool
}

// NewEmulator creates a new Emulator.
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
	if e.timerChan == nil {
		e.timerChan = startTicker(TimerFrequency, e.timerCallback)
	}
}

// stop the background clock timer
func (e *Emulator) stopTimer() {
	close(e.timerChan)
	// Let the goroutine finish
	time.Sleep(2 * TimerFrequency)
	e.timerChan = nil
}

// Stop stops execution of the emulator.
func (e *Emulator) Stop() {
	e.stopTimer()
}

// Beep sounds the speaker.
func (e *Emulator) Beep() {
	fmt.Printf("Beep")
}

// Write sets the memory at addr..address+len(bytes) to the value of the byte slice.
func (e *Emulator) Write(addr uint16, bytes []byte) {
	beg := int(addr)
	if beg >= MemorySize {
		return
	}
	max := int(addr) + len(bytes)
	if max >= MemorySize {
		max = MemorySize
	}
	for i := beg; i < max; i++ {
		e.mem[i] = bytes[i]
	}
}

// Read returns a slice of bytes from memory.
func (e *Emulator) Read(addr uint16, len uint) []byte {
	start := int(addr)
	end := int(addr) + int(len)
	if start >= MemorySize {
		return []byte{}
	}
	if end > MemorySize {
		end = MemorySize
	}
	bytes := make([]byte, end-start)
	for i := start; i < end; i++ {
		bytes[i] = e.mem[i]
	}
	return bytes
}

// GetOpcode returns the two-byte opcode at mem[pc] << 8 | mem[pc+1] and advances the pc.
func (e *Emulator) GetOpcode() uint16 {
	opcode := uint16(e.mem[e.pc]) << 8
	opcode |= uint16(e.mem[e.pc+1])
	e.pc += 2
	return opcode
}

// ClearDisplay sets the display to all 0s.
func (e *Emulator) ClearDisplay() {
	for i := 0; i < DisplayHeight*DisplayWidth; i++ {
		e.display[i] = 0
	}
}

func (e *Emulator) call(a uint16) {
	if e.sp >= StackSize {
		panic("Emulator tack overflow")
	}
	if a >= MemorySize {
		panic("Emulator address out of bounds")
	}
	e.sp++
	e.stack[e.sp] = e.pc
	e.pc = a
}

func (e *Emulator) ret() {
	if e.sp == 0 {
		panic("Emulator stack underflow")
	}
	e.pc = e.stack[e.sp]
	e.sp--
}

func (e *Emulator) timerCallback() {
}
