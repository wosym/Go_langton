package main

import (
	"fmt"
    "time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/yakshaveinc/go-keycodes"
)

const winWidth, winHeight int = 1000, 1000
const gridDim int = 201      //size of grid
const stepTime = 0        //time between ant steps in ms

type position struct {
    x,y int
}

const (
    NORTH = iota
    EAST
    SOUTH
    WEST
    OF
)
func printGrid(grid [][]int) {
    for y := 0; y < gridDim; y++ {
        fmt.Println(grid[y])
    }
}
func drawAnt(renderer *sdl.Renderer, nx, ny int32, antpos position) {
    bw := int32(winWidth)/nx
    bh := int32(winHeight)/ny
    cx := bw * int32(antpos.x) + bw/2
    cy := bh * int32(antpos.y) + bh/2
    //TODO: make ant ellipse-form based on direction?

    ret := gfx.FilledEllipseColor(renderer, cx, cy, bw/2, bh/2, sdl.Color{100, 50, 0, 255})
    if !ret {
        fmt.Println("Error while drawing box")
    }
}
func drawGrid(renderer *sdl.Renderer, grid [][]int, nx, ny int32) {
    bw := int32(winWidth)/nx
    bh := int32(winHeight)/ny
    var ret bool

    for y  := int32(0); y < ny; y++ {
        for x := int32(0); x < nx; x++ {
            //TODO: get color from LUT
            
            if grid[y][x] == 0 {
                ret = gfx.BoxColor(renderer, x*bw, y*bh, (x+1)*bw, (y+1)*bh, sdl.Color{123, 50, 255, 255})
            } else {
                ret = gfx.BoxColor(renderer, x*bw, y*bh, (x+1)*bw, (y+1)*bh, sdl.Color{10, 250, 40, 255})
            }
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

func moveAnt(grid [][]int, antpos *position, antdir *int) bool {
    //move
    switch *antdir {    //TODO: check for out of bounds!    --> what to do then? Stop program?
        case NORTH:
            (*antpos).y--   //pixel coordinates origin is in top left corner
        case EAST:
            (*antpos).x++
        case SOUTH:
            (*antpos).y++
        case WEST:
            (*antpos).x--
        default:
            fmt.Println("Error moving ant: illegal direction")
    }
    if (*antpos).x < 0 || (*antpos).y < 0 || (*antpos).x > gridDim-1 || (*antpos).y > gridDim-1 {
        return true     //pause
    }

    //rotate based on the square we land on
    if grid[(*antpos).y][(*antpos).x] == 0 {    //TODO: check in lookup table for more complex patterns
        *antdir++;
    } else {
        *antdir--;
    }

    //Check for overflows
    if *antdir >= OF{
        *antdir = NORTH;
    } else if *antdir <= -1 {
        *antdir = WEST;
    }

    //Update cell
    grid[(*antpos).y][(*antpos).x]++
    if grid[(*antpos).y][(*antpos).x] >= 2 {    //TODO: this number should be based on the amount of possibilities in the LUT
        grid[(*antpos).y][(*antpos).x] = 0
    }

    return false    //don't pause



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

    //Create grid for the ant
    grid := make([][]int, gridDim)
    for i := 0; i < gridDim; i++ {
        grid[i] = make([]int, gridDim)
    }

    antpos := position{gridDim/2, gridDim/2}
    antdir := NORTH
    fmt.Println("Starting position for ant: ", antpos, "in direction: ", antdir)

    var frameStart time.Time
    var elapsedTime float32
    var running bool = true
    var paused bool = false
    for running {
        frameStart = time.Now()


        if !paused {
            paused = moveAnt(grid, &antpos, &antdir)
            fmt.Println("Position: ", antpos, "direction: ", antdir)


            drawGrid(renderer, grid, int32(gridDim), int32(gridDim));
            drawAnt(renderer, int32(gridDim), int32(gridDim), antpos);
            renderer.Present()

            //printGrid(grid)
        }

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


    elapsedTime = float32(time.Since(frameStart).Seconds())
    if elapsedTime < .005 {
        sdl.Delay(5 - uint32(elapsedTime*1000))
        elapsedTime = float32(time.Since(frameStart).Seconds())
    }

    fmt.Println("framerate: ", 1/elapsedTime)

    sdl.Delay(stepTime)

    }



}
