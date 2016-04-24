package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/makononov/NESGo/cartridge"
	"github.com/makononov/NESGo/cpu"
	"github.com/makononov/NESGo/ppu"
	// "github.com/go-gl/glfw/v3.1/glfw"
)

const windowHeight = 480
const windowWidth = 640

func init() {
	runtime.LockOSThread()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("Initializing console...")

	if len(os.Args) != 2 {
		panic(errors.New("You must provide a ROM file to run."))
	}

	romFile := os.Args[1]
	_, err := ioutil.ReadFile(romFile)
	check(err)

	// Initialize cartridge
	fmt.Println("Reading ROM file and initializing cartridge...")
	cart, err := cartridge.ParseROM(romFile)
	check(err)

	dataBus := make(chan uint8)
	readWriteBus := make(chan int)
	cartridgeControlBus := make(chan uint16)
	ppuControlBus := make(chan uint16)

	// Initialize cpu
	fmt.Println("Initializing CPU...")
	cpu := new(cpu.CPU)
	cpu.Init(ppuControlBus, cartridgeControlBus, readWriteBus, dataBus)

	fmt.Println("Initializing PPU...")
	ppu := new(ppu.PPU)
	ppu.Init(ppuControlBus, readWriteBus, dataBus)

	fmt.Println("Spawning threads...")
	go cart.WaitForReadWrite(cartridgeControlBus, readWriteBus, dataBus)
	go cpu.Run()
	go ppu.Run()
	for {
	}

	// if err := glfw.Init(); err != nil {
	// 	panic(err)
	// }
	// defer glfw.Terminate()
	//
	// window, err := glfw.CreateWindow(windowWidth, windowHeight, "NESGo", nil, nil)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// window.MakeContextCurrent()
	//
	// for !window.ShouldClose() {
	// 	window.SwapBuffers()
	// 	glfw.PollEvents()
	// }
}
