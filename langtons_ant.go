package main

import (
	"fmt"
    "time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/yakshaveinc/go-keycodes"
)

const winWidth, winHeight int = 800, 600

func drawGrid(renderer *sdl.Renderer, nx, ny int32) {
    bw := int32(winWidth)/nx
    bh := int32(winHeight)/ny

    for y  := int32(0); y < ny; y++ {
        for x := int32(0); x < nx; x++ {
            ret := gfx.BoxColor(renderer, x*bw, y*bh, (x+1)*bw, (y+1)*bh, sdl.Color{123, 50, 255, 255})
            if !ret {
                fmt.Println("Error while drawing box")
            }
            ret = gfx.RectangleColor(renderer, x*bw, y*bh, (x+1)*bw, (y+1)*bh, sdl.Color{0, 0, 0, 255})
            if !ret {
                fmt.Println("Error while drawing rect")
            }
        }
    }
}

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

        drawGrid(renderer, 5, 5);
        renderer.Present()


    elapsedTime = float32(time.Since(frameStart).Seconds())
    if elapsedTime < .005 {
        sdl.Delay(5 - uint32(elapsedTime*1000))
        elapsedTime = float32(time.Since(frameStart).Seconds())
    }

    fmt.Println("framerate: ", 1/elapsedTime)


    }



}
