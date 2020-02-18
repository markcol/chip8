package emulator

import (
	"testing"
)

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

func TestGetOpcode(t *testing.T) {
	e := &Emulator{}

	// Write data to memory
	m := []byte{0x01, 0x02, 0x03, 0x04}
	e.Write(0x0000, m)

	if e.pc != 0x0 {
		t.Errorf("pc = %#04x, expected 0x0000", e.pc)
	}

	o := e.GetOpcode()
	if o != 0x0102 {
		t.Errorf("opcode = %#04x, expected 0x0102", o)
	}

	o = e.GetOpcode()
	if o != 0x0304 {
		t.Errorf("opcode = %#04x, expected 0x0304", 0)
	}
}

func TestRegistersZeroAtStartup(t *testing.T) {
	e := &Emulator{}
	for i := 0; i < len(e.registers); i++ {
		if e.registers[i] != 0 {
			t.Errorf("register[%d] = %#2x, expected 0x00", i, e.registers[i])
		}
	}
}
