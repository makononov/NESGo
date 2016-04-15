package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
)

const windowHeight = 480
const windowWidth = 640

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "NESGo", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	for !window.ShouldClose() {
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
