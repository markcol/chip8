package emulator

import (
	"fmt"
	"time"
)

const (
	TIMER_DURATION = time.Second / 60	// 60hz timer for sound and timer updates
	MEMORY_SIZE = 4096
	DISPLAY_HEIGHT = 32
	DISPLAY_WIDTH = 64
	REGISTERS = 16
	STACK_SIZE = 16
)

// Emulator represents an instance of the Chip8 emulator.
type Emulator struct {
	mem [MEMORY_SIZE]byte
	display [DISPLAY_WIDTH * DISPLAY_HEIGHT]byte
	registers [REGISTERS]byte
	stack [STACK_SIZE]uint16
	pc uint16
	sp byte
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
		e.timerChan = startTicker(TIMER_DURATION, e.timerCallback)
	}
}

// stop the background clock timer
func (e *Emulator) stopTimer() {
	close(e.timerChan)
	// Let the goroutine finish
	time.Sleep(2 * TIMER_DURATION)
	e.timerChan = nil
}

// Stop stops execution of the emulator.
func (e *Emulator) Stop() {
	e.stopTimer()
}

func (e *Emulator) Beep() {
	fmt.Printf("Beep")
}

// Write sets the memory at addr..address+len(bytes) to the value of the byte slice.
func (e *Emulator) Write(addr uint16, bytes []byte) {
	beg := int(addr)
	if beg >= MEMORY_SIZE {
		return
	}
	max := int(addr) + len(bytes)
	if max >= MEMORY_SIZE {
		max = MEMORY_SIZE
	}
	for i := beg; i < max; i++ {
		e.mem[i] = bytes[i]
	}
}

// Read returns a slice of bytes from memory.
func (e *Emulator) Read(addr uint16, len uint) []byte {
	start := int(addr)
	end := int(addr) + int(len)
	if start >= MEMORY_SIZE {
		return []byte{}
	}
	if end > MEMORY_SIZE {
		end = MEMORY_SIZE
	}
	bytes := make([]byte, end - start)
	for i := start; i < end; i++ {
		bytes[i] = e.mem[i]
	}
	return bytes
}

// GetPc returns the value of the pc.
func (e *Emulator) GetPc() uint16 {
	return e.pc
}

// SetPc sets the pc to addr. If addr is >= MEMORY_SIZE nothing happens.
func (e *Emulator) SetPc(addr uint16) {
	if addr >= MEMORY_SIZE {
		return
	}
	e.pc = addr
}

// GetPc returns the value of the sp.
func (e *Emulator) GetSp() byte {
	return e.sp
}

// SetPc sets the sp to l. If l is >= STACK_SIZE nothing happens.
func (e *Emulator) SetSp(l byte) {
	if l >= STACK_SIZE {
		return
	}
	e.sp = l
}

// GetStack returns the current stack
func (e *Emulator) GetStack() [STACK_SIZE]uint16 {
	s := [STACK_SIZE]uint16{}
	for i := 0; i < len(s); i++ {
		s[i] = e.stack[i]
	}
	return s
}

// SetStack sets the stack to the current values
func (e *Emulator) SetStack(s [STACK_SIZE]uint16) {
	for i := 0; i < len(s); i++ {
		e.stack[i] = s[i]
	}
}

// GetRegisters returns the current register values.
func (e *Emulator) GetRegisters() [REGISTERS]byte {
	r := [REGISTERS]byte{}
	for i := 0; i < len(r); i++ {
		r[i] = e.registers[i]
	}
	return r
}

// SetRegisters sets the registers to the given values.
func (e *Emulator) SetRegisters(r [REGISTERS]byte) {
	for i := 0; i < len(r); i++ {
		e.registers[i] = r[i]
	}
}

// GetRegister returns the value of register r.
func (e *Emulator) GetRegister(r byte) byte {
	return e.registers[r]
}

// SetRegister sets  register r to the given value.
func (e *Emulator) SetRegister(r byte, v byte) {
	if r >= REGISTERS {
		return
	}
	e.registers[r] = v
}

// GetOpcode returns the two-byte opcode at mem[pc] << 8 | mem[pc+1] and advances the pc.
func (e *Emulator) GetOpcode() uint16 {
	var code uint16

	code = uint16(e.mem[e.pc]) << 8
	code |= uint16(e.mem[e.pc + 1])
	e.pc += 2
	return code
}

// ClearDisplay sets the display to all 0s.
func (e *Emulator) ClearDisplay() {
	for i := 0; i < DISPLAY_HEIGHT * DISPLAY_WIDTH; i++ {
		e.display[i] = 0
	}
}

func (e *Emulator) call(a uint16) {
	if (e.sp >= STACK_SIZE) {
		panic("Emulator tack overflow")
	}
	if (a >= MEMORY_SIZE) {
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