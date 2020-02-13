package emulator

import "testing"

func TestGetRegistersZero(t *testing.T) {
	e := &Emulator{}
	r := e.GetRegisters()
	for i := 0; i < len(r); i++ {
		if r[i] != 0 {
			t.Errorf("register[%d] = %#2X, expected 0", i, r[i])
		}
	}
}
