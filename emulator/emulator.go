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

	// Registers holds the number of v available in the Emulator.
	Registers = 16

	// StackSize holds the size of the stack in the Emulator (max call depth).
	StackSize = 16
)

// Emulator represents an instance of the Chip8 emulator.
type Emulator struct {
	mem       [MemorySize]byte
	display   [DisplayWidth * DisplayHeight]byte
	v         [Registers]byte
	stack     [StackSize]uint16
	pc        uint16
	i         uint16
	sp        byte
	st        byte
	dt        byte
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
	// TODO(markcol): Enable the timer
	// e.startTimer()
}

// Stop stops execution of the emulator.
func (e *Emulator) Stop() {
	e.stopTimer()
}

func (e *Emulator) runCode() {
	opcode := e.GetOpcode()
	switch {
	case opcode == 0x00E0: // CLS
		e.ClearDisplay()
	case opcode == 0x00EE: // RET
		e.ret()
	case opcode&0xF000 == 0x1000: // JP
		e.pc = opcode & 0x0FFF
	case opcode&0xF000 == 0x2000: // CALL
		e.call(opcode & 0x0FFF)
	case opcode&0xF000 == 0x3000: // SE vx, byte
		r := (opcode & 0x0F00) >> 8
		if e.v[r] == byte(opcode) {
			e.pc += 2
		}
	case opcode&0xF000 == 0xA000: // LD I, addr
		addr := opcode & 0x0FFF
		e.i = addr
	case opcode&0xF0FF == 0xF007: // LD ST, Vx
		max := (opcode & 0x0F00) >> 8
		e.v[max] = e.dt
	case opcode&0xF0FF == 0xF018: // LD ST, Vx
		r := (opcode & 0x0F00) >> 8
		e.st = e.v[r]
	case opcode&0xF0FF == 0xF01E: // ADD I, Vx
		r := (opcode & 0x0F00) >> 8
		e.i += uint16(e.v[r])
	case opcode&0xF0FF == 0xF055: // LD[I], Vx
		max := (opcode & 0x0F00) >> 8
		if e.i+max > uint16(len(e.mem)) {
			panic("Address out of range")
		}
		for i := uint16(0); i < max; i++ {
			e.mem[e.i+i] = e.v[i]
		}
	case opcode&0xF0FF == 0xF065: // LD Vx, [I]
		max := (opcode & 0x0F00) >> 8
		if e.i+max > uint16(len(e.mem)) {
			panic("Address out of range")
		}
		for i := uint16(0); i < max; i++ {
			e.v[i] = e.mem[e.i+i]
		}
	default:
	}

}

// WriteOpcode writes an opcode at the given address
func (e *Emulator) WriteOpcode(opcode uint16, addr uint16) {
	if (addr + 1) > MemorySize {
		panic("Address out of range")
	}
	e.mem[addr] = byte(opcode >> 8)
	e.mem[addr+1] = byte(opcode)
}

// ReadOpcode reads an opcode from the given address
func (e *Emulator) ReadOpcode(addr uint16) uint16 {
	if (addr + 1) > MemorySize {
		panic("Address out of range")
	}
	op := uint16(e.mem[addr])<<8 | uint16(e.mem[addr+1])
	//op |=  uint16(e.mem[addr + 1]) & 0x00FF
	return op
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
	for i := 0; i < len(bytes); i++ {
		e.mem[addr+uint16(i)] = bytes[i]
	}
}

// Read returns a slice of bytes from memory.
func (e *Emulator) Read(addr uint16, l uint) []byte {
	start := int(addr)
	end := int(addr) + int(l)
	if start >= MemorySize {
		return []byte{}
	}
	if end > MemorySize {
		end = MemorySize
	}
	bytes := make([]byte, end-start)
	for i := 0; i < int(l); i++ {
		bytes[i] = e.mem[start+i]
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
		panic("Emulator stack overflow")
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

func (e *Emulator) timerCallback() {
}
