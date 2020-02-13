package main

import (
	"github.com/markcol/chip8/emulator"
	"time"
)

func main() {
	e := emulator.NewEmulator()
	e.Start()
	time.Sleep(time.Second / 10)
	e.Stop()
}
