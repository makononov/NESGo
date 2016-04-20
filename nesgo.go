package main

import (
	"errors"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/makononov/NESGo/cartridge"
	// "github.com/go-gl/glfw/v3.1/glfw"
)

const windowHeight = 480
const windowWidth = 640

type console struct {
	cartridge  *cartridge.Cartridge
	dataBus    chan int
	controlBus chan int
	addressBus chan int
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
	cons := new(console)
	cons.dataBus = make(chan int)
	cons.controlBus = make(chan int)
	cons.addressBus = make(chan int)

	if len(os.Args) != 2 {
		panic(errors.New("You must provide a ROM file to run."))
	}

	romFile := os.Args[1]
	_, err := ioutil.ReadFile(romFile)
	check(err)

	cons.cartridge, err = cartridge.ParseROM(romFile)
	check(err)

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
