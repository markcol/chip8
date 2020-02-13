package main

import (
	"fmt"
	"github.com/markcol/chip8/emulator"
)

func main() {
	e := emulator.NewEmulator()
	b := e.Read(0, 4)
	fmt.Println(b)
	e.Write(0, []byte{0x01,0x02,0x03,0x04})
	b = e.Read(0, 4)
	fmt.Println(b)
	o := e.GetOpcode()
	fmt.Printf("opcode: %#04X\n", o)
	r := e.GetRegisters()
	fmt.Printf("registers: %v\n", r)
	e.SetRegisters([emulator.REGISTERS]byte{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16})
	r = e.GetRegisters()
	fmt.Printf("registers: %v\n", r)
	s := e.GetStack()
	fmt.Printf("stack: %v\nb", s)
	e.SetStack([emulator.STACK_SIZE]uint16{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16})
	s = e.GetStack()
	fmt.Printf("stack: %v\n", s)
	v := e.GetRegister(5)
	fmt.Printf("%#02X\n", v)
	e.SetRegister(5, 0x7F)
	v = e.GetRegister(5)
	fmt.Printf("%#02X\n", v)

}
