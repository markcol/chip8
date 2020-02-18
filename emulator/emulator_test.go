package emulator

import (
	"testing"
)

// Test that the register values default to zero at startup.
func TestRegistersZeroAtStartup(t *testing.T) {
	e := &Emulator{}
	for i := 0; i < len(e.v); i++ {
		if e.v[i] != 0 {
			t.Errorf("V%1X = %#2x, expected 0x00", i, e.v[i])
		}
	}
}

// Test that the Read and Write functions work properly given normal inputs.
func TestWriteOpcode(t *testing.T) {
	e := &Emulator{}

	addr := uint16(0x0000 & 0x0FFF)
	opcode := uint16(0x1f7F)

	if e.mem[addr] != 0 {
		t.Errorf("mem[%#04x] = %#02x, expected %#02x", addr, e.mem[addr], 0)
	}
	if e.mem[addr+1] != 0 {
		t.Errorf("mem[%#04x] = %#02x, expected %#02x", addr, e.mem[addr+1], 0)
	}

	e.WriteOpcode(opcode, addr)

	if e.mem[addr] != byte(opcode>>8) {
		t.Errorf("mem[%#04x] = %#02x, expected %#02x", addr, e.mem[addr], opcode>>8)
	}
	if e.mem[addr+1] != byte(opcode) {
		t.Errorf("mem[%#04x] = %#02x, expected %#02x", addr, e.mem[addr+1], byte(opcode))
	}
}

// Test that the Read and Write functions work properly given normal inputs.
func TestReadOpcode(t *testing.T) {
	e := &Emulator{}

	addr := uint16(0x017E & 0x0FFF)
	opcode := uint16(0x1F7F)

	e.WriteOpcode(opcode, addr)
	op := e.ReadOpcode(addr)
	if op != opcode {
		t.Errorf("opcode(%#04x) = %#04x, expected %#04x", addr, op, opcode)
	}
}

// Test that the Read and Write functions work properly given normal inputs.
func TestReadWrite(t *testing.T) {
	l := uint(16)
	addr := uint16(0x00)
	e := &Emulator{}
	// Make sure memory is initialized to zeroes
	b := e.Read(addr, l)
	for i := 0; i < len(b); i++ {
		if b[i] != 0 {
			t.Errorf("mem[%d] = %#2x, expected 0x00", i, b[i])
		}
	}

	// Write data to memory
	m := make([]byte, l)
	for i := 0; i < len(m); i++ {
		m[i] = byte(i)
	}
	e.Write(addr, m)

	// Make sure read returns same data as write
	b = e.Read(addr, uint(len(m)))
	for i := 0; i < len(b); i++ {
		if b[i] != m[i] {
			t.Errorf("mem[%d] = %#2x, expected %#2x", i, b[i], m[i])
		}
	}
}

// Test that GetOpcode returns 16-bit opcodes in big-endian order and that
// the PC advances by two with each read of the opcode.
func TestGetOpcode(t *testing.T) {
	e := &Emulator{}

	// Write data to memory
	m := []byte{0x01, 0x02, 0x03, 0x04}
	e.Write(0x0000, m)

	if e.pc != 0x0 {
		t.Errorf("PC = %#04x, expected 0x0000", e.pc)
	}

	// make sure opcode is read in big-endian form
	o := e.GetOpcode()
	if o != 0x0102 {
		t.Errorf("opcode = %#04x, expected 0x0102", o)
	}

	// make sure PC advanced
	if e.pc != 0x0002 {
		t.Errorf("PC = %#04x, expected 0x0002", e.pc)
	}

	o = e.GetOpcode()
	if o != 0x0304 {
		t.Errorf("opcode = %#04x, expected 0x0304", 0)
	}

	if e.pc != 0x0004 {
		t.Errorf("PC = %#04x, expected 0x0004", e.pc)
	}
}

func TestRetOpcode(t *testing.T) {
	e := &Emulator{}

	addr := uint16(0x100 & 0x0FFF)
	opcode := uint16(0x00EE)
	e.WriteOpcode(opcode, addr)
	oldAddr := uint16(0x135 & 0x0FFF)

	e.pc = addr
	e.sp++
	e.stack[e.sp] = oldAddr

	if e.pc != addr {
		t.Errorf("PC = %#04x, expected %#04x", e.pc, addr)
	}

	if e.sp != 1 {
		t.Errorf("SP = %#02x, expected %#02x", e.sp, 1)
	}

	if e.stack[e.sp] != oldAddr {
		t.Errorf("stack[%#02x] = %#04x, expected %#04x", e.sp, oldAddr, e.stack[e.sp])
	}

	e.runCode()

	if e.pc != oldAddr {
		t.Errorf("PC = %#04x, expected %#04x", e.pc, oldAddr)
	}

	if e.sp != 0 {
		t.Errorf("SP = %#02x, expected %#02x", e.sp, 0)
	}
}

func TestJpOpcode(t *testing.T) {
	e := &Emulator{}

	addr := uint16(0x0135 & 0x0FFF)
	opcode := 0x1000 | addr
	e.WriteOpcode(opcode, uint16(0x0000))

	if e.pc != 0 {
		t.Errorf("PC = %#04x, expected %#04x", e.pc, 0)
	}

	e.runCode()

	if e.pc != addr {
		t.Errorf("PC = %#04x, expected %#04x", e.pc, addr)
	}
}

func TestCallOpcode(t *testing.T) {
	e := &Emulator{}

	addr := uint16(0x135 & 0x0FFF)
	opcode := uint16(0x2000 | addr)
	e.WriteOpcode(opcode, 0x000)

	if e.pc != 0 {
		t.Errorf("PC = %#04x, expected %#04x", e.pc, 0)
	}
	if e.sp != 0 {
		t.Errorf("SP = %#02x, expected %#02x", e.sp, 0)
	}

	oldPc := e.pc + 2
	e.runCode()

	if e.pc != addr {
		t.Errorf("PC = %#04x, expected %#04x", e.pc, addr)
	}
	if e.sp != 1 {
		t.Errorf("SP = %#02x, expected %#02x", e.sp, 1)
	}
	if e.stack[e.sp] != oldPc {
		t.Errorf("stack[%#02x] = %#04x, expected %#04x", e.sp, e.stack[e.sp], oldPc)
	}
}

func TestSeVxBbEqualOpcode(t *testing.T) {
	e := &Emulator{}

	b := byte(0x7F)
	r := byte(6)
	opcode := uint16(0x3000 | uint16(r)<<8 | uint16(b))
	e.WriteOpcode(opcode, 0x000)

	opc := e.pc
	e.v[r] = b
	if e.v[r] != b {
		t.Errorf("V%1X = %#02x, expected %#02x", r, e.v[r], b)
	}

	e.runCode()

	if e.pc != opc+4 {
		t.Errorf("PC = %#04x, expected %#04x", e.pc, opc+4)
	}
}

func TestSeVxBbNotEqualOpcode(t *testing.T) {
	e := &Emulator{}

	b := 0x7F
	r := byte(6)
	opcode := uint16(0x3000 | uint16(r)<<8 | uint16(b))
	e.WriteOpcode(opcode, 0x000)

	opc := e.pc
	e.v[r] = 0x01
	if e.v[r] != 0x01 {
		t.Errorf("V%1X = %#02x, expected %#02x", r, e.v[r], 0x01)
	}

	e.runCode()

	if e.pc != opc+2 {
		t.Errorf("PC = %#04x, expected %#04x", e.pc, opc+2)
	}
}

func TestLdiOpcode(t *testing.T) {
	e := &Emulator{}

	addr := uint16(0x135 & 0x0FFF)
	opcode := uint16(0xA000 | addr)
	e.WriteOpcode(opcode, 0x000)

	if e.i != 0 {
		t.Errorf("I = %#04x, expected %#04x", e.i, 0)
	}

	e.runCode()

	if e.i != addr {
		t.Errorf("I = %#04x, expected %#04x", e.i, addr)
	}
}

func TestLdVxDt(t *testing.T) {
	e := &Emulator{}

	r := byte(6)
	rval := byte(0x17)
	opcode := uint16(0xF007) | uint16(r)<<8
	e.WriteOpcode(opcode, 0x000)

	e.dt = rval
	if e.v[r] != 0 {
		t.Errorf("V%1X = %#02x, expected %#02x", r, e.v[r], 0)
	}
	if e.dt != rval {
		t.Errorf("DT = %#02x, expected %#02x", e.dt, rval)
	}

	e.runCode()

	if e.v[r] != rval {
		t.Errorf("V%1X = %#02x, expected %#02x", r, e.v[r], rval)
	}
}

func TestSneVxBbEqualOpcode(t *testing.T) {
	e := &Emulator{}

	b := byte(0x7F)
	r := byte(6)
	opcode := uint16(0x4000 | uint16(r)<<8 | uint16(b))
	e.WriteOpcode(opcode, 0x000)

	opc := e.pc
	e.v[r] = b
	if e.v[r] != b {
		t.Errorf("V%1X = %#02x, expected %#02x", r, e.v[r], b)
	}

	e.runCode()

	if e.pc != opc+2 {
		t.Errorf("PC = %#04x, expected %#04x", e.pc, opc+2)
	}
}

func TestSneVxBbNotEqualOpcode(t *testing.T) {
	e := &Emulator{}

	b := 0x7F
	r := byte(6)
	opcode := uint16(0x4000 | uint16(r)<<8 | uint16(b))
	e.WriteOpcode(opcode, 0x000)

	opc := e.pc
	e.v[r] = 0x01
	if e.v[r] != 0x01 {
		t.Errorf("V%1X = %#02x, expected %#02x", r, e.v[r], 0x01)
	}

	e.runCode()

	if e.pc != opc+4 {
		t.Errorf("PC = %#04x, expected %#04x", e.pc, opc+4)
	}
}

func TestLdStVx(t *testing.T) {
	e := &Emulator{}

	r := byte(6)
	rval := byte(0x17)
	opcode := uint16(0xF018) | uint16(r)<<8
	e.WriteOpcode(opcode, 0x000)

	e.v[r] = rval
	if e.st != 0 {
		t.Errorf("ST = %#02x, expected %#02x", e.st, 0)
	}
	if e.v[r] != rval {
		t.Errorf("V%1X = %#02x, expected %#02x", r, e.v[r], rval)
	}

	e.runCode()

	if e.st != rval {
		t.Errorf("ST = %#02x, expected %#02x", e.i, rval)
	}
}

func TestAddIVx(t *testing.T) {
	e := &Emulator{}

	r := byte(6)
	rval := byte(0x17)
	opcode := uint16(0xF01E) | uint16(r)<<8
	e.WriteOpcode(opcode, 0x000)
	ival := uint16(0x1700)
	newval := ival + uint16(rval)

	e.v[r] = rval
	e.i = ival
	if e.i != ival {
		t.Errorf("I = %#04x, expected %#04x", e.i, ival)
	}
	if e.v[r] != rval {
		t.Errorf("V%1X = %#02x, expected %#02x", r, e.v[r], rval)
	}

	e.runCode()

	if e.i != newval {
		t.Errorf("I = %#04x, expected %#04x", e.i, newval)
	}
}

func TestLnOpcode(t *testing.T) {
	e := &Emulator{}

	addr := uint16(0x700)
	l := uint16(0x0F)
	regs := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
	e.v = regs
	e.i = addr
	e.WriteOpcode(0xFF55, 0x0000)

	// ensure that the target area is set
	if e.i != addr {
		t.Errorf("I = %#04x, expected %#04x", e.i, addr)
	}
	// ensure that target area is empty
	for i := uint16(0); i < uint16(l); i++ {
		if e.mem[e.i+i] != 0 {
			t.Errorf("mem[%#04x] = %#2x, expected 0x00", i, e.v[i])
		}
	}
	// ensure that v are set
	for i := uint16(0); i < uint16(len(regs)); i++ {
		if e.v[i] != regs[i] {
			t.Errorf("V%1X = %#02x, expected %#02x", i, e.v[i], regs[i])
		}
	}

	e.runCode()

	// ensure that target area is set to the register value
	for i := uint16(0); i < uint16(l); i++ {
		if e.mem[e.i+i] != e.v[i] {
			t.Errorf("mem[%#04x] = %#2x, expected %#02x", e.i+i, e.mem[e.i+i], regs[i])
		}
	}
	if e.mem[e.i+l] != 0 {
		t.Errorf("mem[%#04x] = %#2x, expected %#02x", e.i+l, e.mem[e.i+l], 0)
	}
	// ensure that I still points to the initial address
	if e.i != addr {
		t.Errorf("I = %#04x, expected %#04x", e.i, addr)
	}
}

func TestLdOpcode(t *testing.T) {
	e := &Emulator{}

	addr := uint16(0x0700)
	l := uint16(0x0F)
	regs := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
	e.i = addr
	opcode := 0xF065 | ((l << 8) & 0x0F00)
	e.WriteOpcode(opcode, 0x0000)
	e.Write(addr, regs[0:l+1])

	// ensure that the target area is set
	if e.i != addr {
		t.Errorf("I = %#04x, expected %#04x", e.i, addr)
	}
	for i := uint16(0); i < uint16(l); i++ {
		if e.v[i] != 0 {
			t.Errorf("V%1X = %#2x, expected 0x00", i, e.v[i])
		}
	}
	// ensure that memory is set
	for i := uint16(0); i < uint16(len(regs)); i++ {
		if e.mem[addr+i] != regs[i] {
			t.Errorf("mem[%#04x] = %#02x, expected %#02x", addr+i, e.mem[addr+i], regs[i])
		}
	}

	e.runCode()

	// ensure that target area is set to the register value
	for i := uint16(0); i < uint16(l); i++ {
		if e.v[i] != e.mem[e.i+i] {
			t.Errorf("V%1X = %#2x, expected %#02x", i, e.v[i], e.mem[e.i+i])
		}
	}
	// ensure that I still points to the initial address
	if e.i != addr {
		t.Errorf("I = %#04x, expected %#04x", e.i, addr)
	}
}
