package main

import (
	"fmt"
    "time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/yakshaveinc/go-keycodes"
)

const winWidth, winHeight int = 800, 600

func main() {
    err := sdl.Init(sdl.INIT_EVERYTHING)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer sdl.Quit()

    window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer window.Destroy()

    renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer renderer.Destroy()

    tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
    if err != nil {
        fmt.Println(err)
        return
    }
    defer tex.Destroy()

    var frameStart time.Time
    var elapsedTime float32
    var running bool = true
    for running {
        frameStart = time.Now()

        for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
            switch t := event.(type) {
                case *sdl.QuitEvent:
                    return
                case *sdl.KeyboardEvent:
                    if(t.Type == sdl.KEYDOWN) {
                        fmt.Println("key pressed: ", t.Keysym.Scancode)
                        if uint16(t.Keysym.Scancode) == keycodes.KeyEscape{
                            fmt.Println("escape pressed, exiting...")
                            running = false
                        }
                    }
            }
        }

        ret := gfx.BoxColor(renderer, 10, 10, 100, 100, sdl.Color{123, 50, 255, 100})
        if ret != true {
            fmt.Println("Error while drawing box")
        }
        ret = gfx.RectangleColor(renderer, 10, 10, 1000, 1000, sdl.Color{123, 50, 255, 100})
        if ret != true {
            fmt.Println("Error while drawing rect")
        }

        renderer.Present()


    elapsedTime = float32(time.Since(frameStart).Seconds())
    if elapsedTime < .005 {
        sdl.Delay(5 - uint32(elapsedTime*1000))
        elapsedTime = float32(time.Since(frameStart).Seconds())
    }

    fmt.Println("framerate: ", 1/elapsedTime)


    }



}
