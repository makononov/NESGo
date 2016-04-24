package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/makononov/NESGo/cartridge"
	"github.com/makononov/NESGo/cpu"
	// "github.com/go-gl/glfw/v3.1/glfw"
)

const windowHeight = 480
const windowWidth = 640

type console struct {
	cartridge *cartridge.Cartridge
	cpu       *cpu.CPU
}

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
	cons := new(console)

	if len(os.Args) != 2 {
		panic(errors.New("You must provide a ROM file to run."))
	}

	romFile := os.Args[1]
	_, err := ioutil.ReadFile(romFile)
	check(err)

	// Initialize cartridge
	fmt.Println("Reading ROM file and initializing cartridge...")
	cons.cartridge, err = cartridge.ParseROM(romFile)
	check(err)

	dataBus := make(chan uint8)
	readWriteBus := make(chan int)
	cartridgeControlBus := make(chan uint16)

	// Initialize cpu
	fmt.Println("Initializing CPU...")
	cons.cpu = new(cpu.CPU)
	cons.cpu.Init(cartridgeControlBus, readWriteBus, dataBus)

	fmt.Println("Spawning threads...")
	go cons.cartridge.WaitForReadWrite(cartridgeControlBus, readWriteBus, dataBus)
	go cons.cpu.Run()
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
